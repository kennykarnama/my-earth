CREATE TABLE my_earth.location (
    id SERIAL PRIMARY KEY,
    city VARCHAR(1024) UNIQUE NOT NULL,
    weather_summary VARCHAR(1024),
    temperature REAL,
    wind_speed REAL,
    wind_angle REAL,
    wind_direction VARCHAR(1024),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);