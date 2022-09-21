// Package seeds seeds the database.
package seeds

import (
	"database/sql"
	"math/rand"

	faker "github.com/bxcodec/faker/v3"
	"github.com/jmoiron/sqlx"
)

// Seed struct.
type Seed struct {
	db *sqlx.DB
}

// NewSeed return a Seed with a pool of connection to a dabase.
func NewSeed(db *sqlx.DB) Seed {
	return Seed{
		db: db,
	}
}

// PopulateDB inserts data into the database.
func PopulateDB(db *sql.DB) error {
	var id int

	// inserts roles.
	if _, err := db.Exec(`INSERT INTO roles(name) VALUES ($1)`, "admin"); err != nil {
		return err
	}

	if _, err := db.Exec(`INSERT INTO roles(name) VALUES ($1)`, "user"); err != nil {
		return err
	}

	// inserts users with admin permissions.
	if err := db.QueryRow(`SELECT id FROM roles WHERE name = 'admin'`).Scan(&id); err != nil {
		return err
	}

	for i := 0; i < 50; i++ {
		_, err := db.Exec(
			`INSERT INTO users(username, first_name, last_name, role_id) VALUES ($1, $2, $3, $4)`,
			faker.Username(),
			faker.FirstName(),
			faker.LastName(),
			id,
		)
		if err != nil {
			return err
		}
	}

	// inserts users with regular permissions.
	if err := db.QueryRow(`SELECT id FROM roles WHERE name = 'user'`).Scan(&id); err != nil {
		return err
	}

	for i := 0; i < 50; i++ {
		_, err := db.Exec(
			`INSERT INTO users(username, first_name, last_name, role_id) VALUES ($1, $2, $3, $4)`,
			faker.Username(),
			faker.FirstName(),
			faker.LastName(),
			id,
		)
		if err != nil {
			return err
		}
	}

	// inserts products.
	for i := 0; i < 100; i++ {
		var err error

		if _, err = db.Exec(
			`INSERT INTO products(name, price) VALUES ($1, $2)`, faker.Word(), rand.Float32(),
		); err != nil {
			return err
		}
	}
	return nil
}

// PopulateTx inserts data into the database thoguh the use of database tx.
func PopulateTx(db *sql.Tx) error {
	var id int

	// inserts roles.
	if _, err := db.Exec(`INSERT INTO roles(name) VALUES ($1)`, "admin"); err != nil {
		return err
	}

	if _, err := db.Exec(`INSERT INTO roles(name) VALUES ($1)`, "user"); err != nil {
		return err
	}

	// inserts users with admin permissions.
	if err := db.QueryRow(`SELECT id FROM roles WHERE name = 'admin'`).Scan(&id); err != nil {
		return err
	}

	for i := 0; i < 50; i++ {
		_, err := db.Exec(
			`INSERT INTO users(username, first_name, last_name, role_id) VALUES ($1, $2, $3, $4)`,
			faker.Username(),
			faker.FirstName(),
			faker.LastName(),
			id,
		)
		if err != nil {
			return err
		}
	}

	// inserts users with regular permissions.
	if err := db.QueryRow(`SELECT id FROM roles WHERE name = 'user'`).Scan(&id); err != nil {
		return err
	}

	for i := 0; i < 50; i++ {
		_, err := db.Exec(
			`INSERT INTO users(username, first_name, last_name, role_id) VALUES ($1, $2, $3, $4)`,
			faker.Username(),
			faker.FirstName(),
			faker.LastName(),
			id,
		)
		if err != nil {
			return err
		}
	}

	// inserts products.
	for i := 0; i < 100; i++ {
		var err error

		if _, err = db.Exec(
			`INSERT INTO products(name, price) VALUES ($1, $2)`, faker.Word(), rand.Float32(),
		); err != nil {
			return err
		}
	}
	return nil
}
