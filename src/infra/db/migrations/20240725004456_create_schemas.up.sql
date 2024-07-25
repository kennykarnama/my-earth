CREATE SCHEMA IF NOT EXISTS my_earth;

GRANT USAGE, CREATE ON SCHEMA my_earth TO me_read_write;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA my_earth TO me_read_write;
ALTER DEFAULT PRIVILEGES IN SCHEMA my_earth GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO me_read_write;

GRANT USAGE ON ALL SEQUENCES IN SCHEMA my_earth TO me_read_write;
ALTER DEFAULT PRIVILEGES IN SCHEMA my_earth GRANT USAGE ON SEQUENCES TO me_read_write;

GRANT me_read_write TO my_earth_service;
