package domain

import (
	"context"
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

type LocationRepository interface {
	SaveLoc(ctx context.Context, loc *Location) error
	GetLocDetailByID(ctx context.Context, ID int) (*LocationWeather, error)
	UpdateWeather(ctx context.Context, locID int, w *Weather) error
	GetExpiringWeather(ctx context.Context) ([]LocationWeather, error)
}

type ListLocations struct {
	Items []Location
}
type LocationQuery interface {
	FindLocationPoint(ctx context.Context, name string) (*ListLocations, error)
}
