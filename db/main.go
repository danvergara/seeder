package main

import (
	"log"

	"github.com/danvergara/seeder/db/seeds"
	"github.com/danvergara/seeder/pkg/seeder"
	"github.com/jmoiron/sqlx"

	// postgres driver
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "postgres://postgres:password@db:5432/users?sslmode=disable")
	if err != nil {
		log.Fatalf("error opening a connection with the database %s\n", err)
	}
	s := seeds.NewSeed(db)

	if err := seeder.Execute(s); err != nil {
		log.Fatalf("error seeding the db %s\n", err)
	}
}
