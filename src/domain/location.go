package domain

import (
	"context"
	"fmt"
	"time"
)

type Location struct {
	ID        int32
	RefID     string
	Name      string
	Lat       float64
	Lon       float64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type LocationWeather struct {
	Location
	Weather
}

type ListLocationsQuery struct {
	City *string
	ID   *int32
}

type LocationRepository interface {
	SaveLoc(ctx context.Context, loc *Location) error
	GetLocDetailByID(ctx context.Context, ID int) (*LocationWeather, error)
	UpdateWeather(ctx context.Context, locID int, w *Weather) error
	GetExpiringWeather(ctx context.Context) ([]LocationWeather, error)
	ListLocations(ctx context.Context, q ListLocationsQuery) ([]LocationWeather, error)
}

type ListLocations struct {
	Items     []Location
	expiredAt time.Time
	tz        string
}

func (l *ListLocations) SetExpiredAt(s string, layout string) error {
	expiredAt, err := time.Parse(layout, s)
	if err != nil {
		return fmt.Errorf("weather.setExpiredAt err: %w", err)
	}

	l.expiredAt = expiredAt

	l.tz, _ = expiredAt.Zone()

	return nil
}

func (l *ListLocations) ExpiredAt() time.Time {
	return l.expiredAt
}

func (l *ListLocations) ExpiredAtUTC() time.Time {
	return l.expiredAt.UTC()
}

type LocationQuery interface {
	FindLocationPoint(ctx context.Context, name string) (*ListLocations, error)
}
