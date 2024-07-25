package main

import (
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	DSN string `envconfig:"DSN"`
}

//go:embed *.sql
var fs embed.FS

func main() {
	cfg := config{}

	err := envconfig.Process("myapp", &cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	d, err := iofs.New(fs, ".")
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", d, "postgres://postgres:password@localhost:5432/mvp?sslmode=disable&x-migrations-table=schema_migrations_my_earth")
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}
}
