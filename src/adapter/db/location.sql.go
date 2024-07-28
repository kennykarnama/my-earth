// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: location.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createNewLocation = `-- name: CreateNewLocation :one
INSERT INTO my_earth.location (city, latitude, longitude, created_at, updated_at)
 VALUES($1, $2, $3, $4, $5) RETURNING id
`

type CreateNewLocationParams struct {
	City      string
	Latitude  pgtype.Float8
	Longitude pgtype.Float8
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func (q *Queries) CreateNewLocation(ctx context.Context, arg CreateNewLocationParams) (int32, error) {
	row := q.db.QueryRow(ctx, createNewLocation,
		arg.City,
		arg.Latitude,
		arg.Longitude,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

const enrichWeatherInfo = `-- name: EnrichWeatherInfo :exec
UPDATE my_earth.location SET weather_summary = $1,
  wind_angle = $2, wind_direction = $3,
  wind_speed = $4, updated_at = $5, temperature = $6, 
  expired_at = $7, temperature_unit = $8
WHERE id = $9
`

type EnrichWeatherInfoParams struct {
	WeatherSummary  pgtype.Text
	WindAngle       pgtype.Float4
	WindDirection   pgtype.Text
	WindSpeed       pgtype.Float4
	UpdatedAt       pgtype.Timestamptz
	Temperature     pgtype.Float4
	ExpiredAt       pgtype.Timestamptz
	TemperatureUnit pgtype.Text
	ID              int32
}

func (q *Queries) EnrichWeatherInfo(ctx context.Context, arg EnrichWeatherInfoParams) error {
	_, err := q.db.Exec(ctx, enrichWeatherInfo,
		arg.WeatherSummary,
		arg.WindAngle,
		arg.WindDirection,
		arg.WindSpeed,
		arg.UpdatedAt,
		arg.Temperature,
		arg.ExpiredAt,
		arg.TemperatureUnit,
		arg.ID,
	)
	return err
}

const getLocationByID = `-- name: GetLocationByID :one
SELECT id, city, weather_summary, temperature, wind_speed, wind_angle, wind_direction, latitude, longitude, created_at, updated_at, deleted_at, expired_at, temperature_unit FROM my_earth.location 
 WHERE location.id = $1 AND location.deleted_at IS NULL
`

func (q *Queries) GetLocationByID(ctx context.Context, id int32) (MyEarthLocation, error) {
	row := q.db.QueryRow(ctx, getLocationByID, id)
	var i MyEarthLocation
	err := row.Scan(
		&i.ID,
		&i.City,
		&i.WeatherSummary,
		&i.Temperature,
		&i.WindSpeed,
		&i.WindAngle,
		&i.WindDirection,
		&i.Latitude,
		&i.Longitude,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.ExpiredAt,
		&i.TemperatureUnit,
	)
	return i, err
}

const getLocationByName = `-- name: GetLocationByName :many
SELECT id, city, weather_summary, temperature, wind_speed, wind_angle, wind_direction, latitude, longitude, created_at, updated_at, deleted_at, expired_at, temperature_unit FROM my_earth.location 
 WHERE city LIKE $1 AND deleted_at IS NULL
`

func (q *Queries) GetLocationByName(ctx context.Context, city string) ([]MyEarthLocation, error) {
	rows, err := q.db.Query(ctx, getLocationByName, city)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MyEarthLocation
	for rows.Next() {
		var i MyEarthLocation
		if err := rows.Scan(
			&i.ID,
			&i.City,
			&i.WeatherSummary,
			&i.Temperature,
			&i.WindSpeed,
			&i.WindAngle,
			&i.WindDirection,
			&i.Latitude,
			&i.Longitude,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.ExpiredAt,
			&i.TemperatureUnit,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectAllLocations = `-- name: SelectAllLocations :many
SELECT id, city, weather_summary, temperature, wind_speed, wind_angle, wind_direction, latitude, longitude, created_at, updated_at, deleted_at, expired_at, temperature_unit FROM my_earth.location WHERE location.deleted_at IS NULL
`

func (q *Queries) SelectAllLocations(ctx context.Context) ([]MyEarthLocation, error) {
	rows, err := q.db.Query(ctx, selectAllLocations)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MyEarthLocation
	for rows.Next() {
		var i MyEarthLocation
		if err := rows.Scan(
			&i.ID,
			&i.City,
			&i.WeatherSummary,
			&i.Temperature,
			&i.WindSpeed,
			&i.WindAngle,
			&i.WindDirection,
			&i.Latitude,
			&i.Longitude,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.ExpiredAt,
			&i.TemperatureUnit,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectAllLocationsExpiring = `-- name: SelectAllLocationsExpiring :many
SELECT id, city, weather_summary, temperature, wind_speed, wind_angle, wind_direction, latitude, longitude, created_at, updated_at, deleted_at, expired_at, temperature_unit FROM my_earth.location
WHERE location.deleted_at IS NULL 
AND (location.expired_at IS NULL OR location.expired_at < $1)
ORDER BY location.expired_at
`

func (q *Queries) SelectAllLocationsExpiring(ctx context.Context, expiredAt pgtype.Timestamptz) ([]MyEarthLocation, error) {
	rows, err := q.db.Query(ctx, selectAllLocationsExpiring, expiredAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []MyEarthLocation
	for rows.Next() {
		var i MyEarthLocation
		if err := rows.Scan(
			&i.ID,
			&i.City,
			&i.WeatherSummary,
			&i.Temperature,
			&i.WindSpeed,
			&i.WindAngle,
			&i.WindDirection,
			&i.Latitude,
			&i.Longitude,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
			&i.ExpiredAt,
			&i.TemperatureUnit,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
