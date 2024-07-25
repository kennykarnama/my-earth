package main

import (
	"context"
	"log"

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
	ctx := context.Background()

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

	locHandler := port.NewHttpHandler(locSvc)

	r := gin.Default()

	genapi.RegisterHandlers(r, locHandler)

	log.Printf("serving http on port %s", cfg.HTTPPort)

	err = r.Run(cfg.HTTPPort)
	if err != nil {
		log.Fatal(err)
	}
}
