-- name: CreateNewLocation :one
INSERT INTO my_earth.location (city, latitude, longitude, created_at, updated_at)
 VALUES($1, $2, $3, $4, $5) RETURNING id;

-- name: GetLocationByID :one
SELECT * FROM my_earth.location 
 WHERE location.id = $1 AND location.deleted_at IS NULL;

-- name: EnrichWeatherInfo :exec
UPDATE my_earth.location SET weather_summary = $1,
  wind_angle = $2, wind_direction = $3,
  wind_speed = $4, updated_at = $5, temperature = $6, 
  expired_at = $7, temperature_unit = $8
WHERE id = $9;

-- name: SelectAllLocations :many
SELECT * FROM my_earth.location WHERE location.deleted_at IS NULL;

-- name: SelectAllLocationsExpiring :many
SELECT * FROM my_earth.location
WHERE location.deleted_at IS NULL 
AND (location.expired_at IS NULL OR location.expired_at < $1)
ORDER BY location.expired_at;