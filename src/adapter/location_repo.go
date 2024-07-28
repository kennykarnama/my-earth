package adapter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/kennykarnama/my-earth/src/adapter/db"
	"github.com/kennykarnama/my-earth/src/domain"
	"github.com/kennykarnama/my-earth/src/pkg/generr"
	"github.com/kennykarnama/my-earth/src/pkg/psql"
	pgxlog "github.com/mcosta74/pgx-slog"
)

type LocationRepo struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

func NewLocationRepo(ctx context.Context, dsn string) (*LocationRepo, error) {
	// Create new pool with default config
	// * max number of connections: max(4,number of CPUs)
	// * idle timeout of 30 minutes.
	connConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse dsn config: %w", err)
	}

	slogger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// FYI: NewLoggerAdapter: https://github.com/mcosta74/pgx-slog
	adapterLogger := pgxlog.NewLogger(slogger)

	ms := psql.MultiQueryTracer{
		Tracers: []pgx.QueryTracer{
			// tracer: https://github.com/exaring/otelpgx
			otelpgx.NewTracer(),

			// logger
			&tracelog.TraceLog{
				Logger:   adapterLogger,
				LogLevel: tracelog.LogLevelTrace,
			},
		},
	}

	connConfig.ConnConfig.Tracer = &ms

	pool, err := pgxpool.NewWithConfig(ctx, connConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	queries := db.New(pool)

	return &LocationRepo{
		pool: pool,
		q:    queries,
	}, nil
}

func (l *LocationRepo) SaveLoc(ctx context.Context, loc *domain.Location) error {
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("location.save failed to begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("saveLoc", slog.String("err", err.Error()))
		}
	}()

	ptx := l.q.WithTx(tx)
	genID, err := ptx.CreateNewLocation(ctx, db.CreateNewLocationParams{
		City:      loc.Name,
		Latitude:  psql.F64ToPGFloat8(loc.Lat),
		Longitude: psql.F64ToPGFloat8(loc.Lon),
		CreatedAt: psql.TimeToPGTimestampz(loc.CreatedAt),
		UpdatedAt: psql.TimeToPGTimestampz(loc.CreatedAt),
	})
	if err != nil {
		if psql.IsUniqueViolationErr(err) {
			return generr.NewDuplicateErr(err.Error())
		}

		return fmt.Errorf("location.save err: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("location.save err: %w", err)
	}

	loc.ID = genID

	return nil
}

func (l *LocationRepo) GetLocDetailByID(ctx context.Context, id int) (*domain.LocationWeather, error) {
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("location.get failed to begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("getLoc", slog.String("err", err.Error()))
		}
	}()
	ptx := l.q.WithTx(tx)
	dbLoc, err := ptx.GetLocationByID(ctx, int32(id))
	if err != nil {
		return nil, fmt.Errorf("location.get err: %w", err)
	}

	lw := &domain.LocationWeather{
		Location: domain.Location{
			Name:      dbLoc.City,
			Lat:       dbLoc.Latitude.Float64,
			Lon:       dbLoc.Longitude.Float64,
			CreatedAt: psql.PGTimestampzToTimestamp(dbLoc.CreatedAt),
			UpdatedAt: psql.PGTimestampzToTimestamp(dbLoc.UpdatedAt),
		},
		Weather: domain.Weather{
			Summary: psql.PGTextToString(dbLoc.WeatherSummary),
			Temperature: domain.Temperature{
				Value: dbLoc.Temperature.Float32,
				Unit:  "metric",
			},
			WindSpeed:     dbLoc.Temperature.Float32,
			WindAngle:     dbLoc.Temperature.Float32,
			WindDirection: psql.PGTextToString(dbLoc.WindDirection),
		},
	}

	return lw, nil
}

func (l *LocationRepo) UpdateWeather(ctx context.Context, locID int, w *domain.Weather) error {
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("location.update.weather failed to begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("updateWeather", slog.String("err", err.Error()))
		}
	}()
	ptx := l.q.WithTx(tx)

	err = ptx.EnrichWeatherInfo(ctx, db.EnrichWeatherInfoParams{
		WeatherSummary:  psql.StrToPGText(w.Summary),
		WindAngle:       psql.F32ToPGFloat4(w.WindAngle),
		WindDirection:   psql.StrToPGText(w.WindDirection),
		WindSpeed:       psql.F32ToPGFloat4(w.WindSpeed),
		UpdatedAt:       psql.TimeToPGTimestampz(time.Now().UTC()),
		Temperature:     psql.F32ToPGFloat4(w.Temperature.Value),
		TemperatureUnit: psql.StrToPGText(w.Temperature.Unit),
		ID:              int32(locID),
		ExpiredAt:       psql.TimeToPGTimestampz(w.ExpiredAtUTC()),
	})
	if err != nil {
		return fmt.Errorf("location.update.weather err: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("location.update.weather err: %w", err)
	}

	return nil
}

func (l *LocationRepo) GetExpiringWeather(ctx context.Context) ([]domain.LocationWeather, error) {
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("location.get failed to begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("getExpiringWeather", slog.String("err", err.Error()))
		}
	}()

	ptx := l.q.WithTx(tx)

	tc := time.Now().UTC()
	log.Println("getExpiringWeather by expiration < ", tc.String())

	dbLocations, err := ptx.SelectAllLocationsExpiring(ctx, psql.TimeToPGTimestampz(tc))
	if err != nil {
		return nil, fmt.Errorf("location.getExpiring err: %w", err)
	}

	var locs []domain.LocationWeather

	for _, dbLoc := range dbLocations {
		locs = append(locs, domain.LocationWeather{
			Location: domain.Location{
				ID:        dbLoc.ID,
				Name:      dbLoc.City,
				Lat:       dbLoc.Latitude.Float64,
				Lon:       dbLoc.Longitude.Float64,
				CreatedAt: psql.PGTimestampzToTimestamp(dbLoc.CreatedAt),
				UpdatedAt: psql.PGTimestampzToTimestamp(dbLoc.UpdatedAt),
			},
			Weather: domain.Weather{
				Summary: psql.PGTextToString(dbLoc.WeatherSummary),
				Temperature: domain.Temperature{
					Unit:  psql.PGTextToString(dbLoc.TemperatureUnit),
					Value: dbLoc.Temperature.Float32,
				},
				WindSpeed:     dbLoc.WindSpeed.Float32,
				WindAngle:     dbLoc.WindAngle.Float32,
				WindDirection: psql.PGTextToString(dbLoc.WindDirection),
			},
		})
	}

	return locs, nil
}

func (l *LocationRepo) ListLocations(ctx context.Context, q domain.ListLocationsQuery) ([]domain.LocationWeather, error) {
	tx, err := l.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("location.list failed to begin transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			slog.Error("listLocations", slog.String("err", err.Error()))
		}
	}()

	ptx := l.q.WithTx(tx)

	var dbLocations []db.MyEarthLocation

	if q.ID != nil {
		dbLocation, err := ptx.GetLocationByID(ctx, *q.ID)
		if err != nil {
			if psql.NoRowsErr(err) {
				return []domain.LocationWeather{}, nil
			}

			return nil, fmt.Errorf("location.list failed to begin transaction: %w", err)
		}

		dbLocations = append(dbLocations, dbLocation)
	} else if q.City != nil {
		dbLocations, err = ptx.GetLocationByName(ctx, *q.City)
		if err != nil {
			return nil, fmt.Errorf("location.list err: %w", err)
		}
	}

	var locs []domain.LocationWeather

	for _, dbLoc := range dbLocations {
		locs = append(locs, domain.LocationWeather{
			Location: domain.Location{
				ID:        dbLoc.ID,
				Name:      dbLoc.City,
				Lat:       dbLoc.Latitude.Float64,
				Lon:       dbLoc.Longitude.Float64,
				CreatedAt: psql.PGTimestampzToTimestamp(dbLoc.CreatedAt),
				UpdatedAt: psql.PGTimestampzToTimestamp(dbLoc.UpdatedAt),
			},
			Weather: domain.Weather{
				Summary: psql.PGTextToString(dbLoc.WeatherSummary),
				Temperature: domain.Temperature{
					Unit:  psql.PGTextToString(dbLoc.TemperatureUnit),
					Value: dbLoc.Temperature.Float32,
				},
				WindSpeed:     dbLoc.WindSpeed.Float32,
				WindAngle:     dbLoc.WindAngle.Float32,
				WindDirection: psql.PGTextToString(dbLoc.WindDirection),
			},
		})
	}

	return locs, nil
}
