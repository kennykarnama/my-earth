package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kennykarnama/my-earth/api/openapi/genapi"
	"github.com/kennykarnama/my-earth/cmd/earth-server/config"
	"github.com/kennykarnama/my-earth/src/adapter"
	"github.com/kennykarnama/my-earth/src/app"
	"github.com/kennykarnama/my-earth/src/pkg/workerpool"
	"github.com/kennykarnama/my-earth/src/port"
)

func main() {
	cfg := config.Get()
	ctx, cancelFunc := context.WithCancel(context.Background())

	locRepo, err := adapter.NewLocationRepo(ctx, cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}

	meteoSourceCli, err := genapi.NewClientWithResponses(cfg.MeteoSourceBaseURL)
	if err != nil {
		log.Fatal(err)
	}

	meteoSourceRepo := adapter.NewMeteoSource(meteoSourceCli, cfg.MeteoSourceAPIKey)

	wp := workerpool.New(10)
	wp.Start()

	locSvc := app.NewLocationSvc(locRepo, meteoSourceRepo, wp)
	weatherRefresher := app.NewSimpleWeatherRefresher(ctx, locSvc, 1*time.Second)

	locHandler := port.NewHttpHandler(locSvc)

	r := gin.Default()

	genapi.RegisterHandlers(r, locHandler)

	log.Printf("serving http on port %s", cfg.HTTPPort)

	srv := &http.Server{
		Addr:              cfg.HTTPPort,
		Handler:           r.Handler(),
		ReadHeaderTimeout: 20 * time.Second,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen.http.server", slog.Any("err", err))
		}
	}()

	go func() {
		weatherRefresher.Watch()
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutdown Server ...")

	cancelFunc()

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	if err := srv.Shutdown(ctx2); err != nil {
		slog.Error("Server Shutdown", slog.Any("err", err))
	}

	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx2.Done()

	slog.Info("Server exiting")
}
