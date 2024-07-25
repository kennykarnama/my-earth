SELECT 'CREATE DATABASE mvp' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'mvp')\gexec
CREATE USER my_earth_service WITH ENCRYPTED PASSWORD 'my_earth_service_password';

CREATE ROLE me_read_write;

CREATE TABLE IF NOT EXISTS public.schema_migrations(
  version bigint NOT NULL,
  dirty boolean NOT NULL
);

REVOKE SELECT, INSERT, UPDATE, DELETE ON public.schema_migrations FROM PUBLIC;