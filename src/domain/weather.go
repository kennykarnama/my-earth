package domain

import (
	"context"
	"fmt"
	"time"
)

type Temperature struct {
	Value float32
	Unit  string
}

type Weather struct {
	Summary       string
	Temperature   Temperature
	WindSpeed     float32
	WindAngle     float32
	WindDirection string
	expiredAt     time.Time
	tz            string
}

func (w *Weather) SetExpiredAt(s string, timeLoc *time.Location, layout string) error {
	expiredAt, err := time.ParseInLocation(layout, s, timeLoc)
	if err != nil {
		return fmt.Errorf("weather.setExpiredAt err: %w", err)
	}

	w.expiredAt = expiredAt

	w.tz = timeLoc.String()

	return nil
}

func (w *Weather) ExpiredAt() time.Time {
	return w.expiredAt
}

func (w *Weather) ExpiredAtUTC() time.Time {
	return w.expiredAt.UTC()
}

type WeatherByPoint struct {
	Weather
	Lat float64
	Lon float64
}

type WeatherRepository interface {
	FindByPoint(ctx context.Context, lat, lon float64) (*WeatherByPoint, error)
}

type WeatherRefresher interface {
	Refresh(ctx context.Context, req any) error
}
