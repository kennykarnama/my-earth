package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/kennykarnama/my-earth/src/domain"
	"github.com/kennykarnama/my-earth/src/pkg/workerpool"
)

type LocationSvc struct {
	repo        domain.LocationRepository
	qrepo       domain.LocationQuery
	weatherRepo domain.WeatherRepository
	workerPool  *workerpool.WorkerPool
}

func NewLocationSvc(repo domain.LocationRepository,
	locQueryRepo domain.LocationQuery,
	weatherRepo domain.WeatherRepository,
	wp *workerpool.WorkerPool,
) *LocationSvc {
	return &LocationSvc{
		repo:        repo,
		weatherRepo: weatherRepo,
		workerPool:  wp,
		qrepo:       locQueryRepo,
	}
}

type SaveLocReq struct {
	Name string
	Lat  float64
	Lon  float64
}

type SaveLocResp struct {
	Loc *domain.Location
}

func (l *LocationSvc) SaveLoc(ctx context.Context, req *SaveLocReq) (*SaveLocResp, error) {
	savedLoc := &domain.Location{
		Name:      req.Name,
		Lat:       req.Lat,
		Lon:       req.Lon,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	err := l.repo.SaveLoc(ctx, savedLoc)
	if err != nil {
		return nil, err
	}

	return &SaveLocResp{
		Loc: savedLoc,
	}, nil
}

type RefreshWeatherResp struct {
	NumUpdated int
}

func (l *LocationSvc) RefreshWeather(ctx context.Context) (*RefreshWeatherResp, error) {
	locs, err := l.repo.GetExpiringWeather(ctx)
	if err != nil {
		return nil, err
	}

	slog.Info("expired weathers", slog.Any("len", len(locs)))

	for _, loc := range locs {

		fn := func() (any, error) {
			weather, err := l.weatherRepo.FindByPoint(context.Background(), loc.Lat, loc.Lon)
			if err != nil {
				return nil, fmt.Errorf("weatherEnricher.executer err: %w", err)
			}

			err = l.repo.UpdateWeather(context.Background(), int(loc.ID), &weather.Weather)
			if err != nil {
				return nil, fmt.Errorf("weatherEnricher.executer err: %w", err)
			}

			slog.Info(fmt.Sprintf("Data: %s successfully saved asynchronously", loc.Name))

			return weather, nil
		}

		t := workerpool.Task{
			ID:         loc.Name,
			Payload:    loc,
			Executor:   fn,
			OmitResult: true,
		}

		slog.Info("submit job", slog.Any("taskID", t.ID))

		l.workerPool.Submit(t)
	}

	return &RefreshWeatherResp{
		NumUpdated: len(locs),
	}, nil
}

type FindLocationsCoordinatesResp struct {
	Matches *domain.ListLocations
}

func (l *LocationSvc) FindLocationsCoordinates(ctx context.Context, q string) (*FindLocationsCoordinatesResp, error) {
	matches, err := l.qrepo.FindLocationPoint(ctx, q)
	if err != nil {
		return nil, err
	}

	return &FindLocationsCoordinatesResp{
		Matches: matches,
	}, nil
}
