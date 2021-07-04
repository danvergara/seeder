package main

import (
	"log"

	"github.com/danvergara/seeder/db/seeds"
	"github.com/danvergara/seeder/pkg/seeder"
)

func main() {
	s := seeds.Seed{}

	if err := seeder.Execute(s); err != nil {
		log.Fatalf("error seeding the db %s\n", err)
	}
}
