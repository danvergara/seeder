package seeds

import (
	"math/rand"

	"github.com/bxcodec/faker/v3"
	"github.com/jmoiron/sqlx"
)

// SeedTx struct.
type SeedTx struct {
	tx *sqlx.Tx
}

// NewSeedTx returns a seed object with a database tx.
func NewSeedTx(tx *sqlx.Tx) SeedTx {
	return SeedTx{
		tx: tx,
	}
}

// RolesSeed seeds roles data.
func (s SeedTx) RolesSeed() error {
	var err error

	_, err = s.tx.Exec(`INSERT INTO roles(name) VALUES ($1)`, "admin")
	if err != nil {
		return err
	}

	_, err = s.tx.Exec(`INSERT INTO roles(name) VALUES ($1)`, "user")
	if err != nil {
		return err
	}

	return err
}

// UsersSeed seeds roles data.
func (s SeedTx) UsersSeed() error {
	var id int
	var err error

	err = s.tx.Get(&id, `SELECT id FROM roles WHERE name = 'admin'`)
	if err != nil {
		return err
	}

	for i := 0; i < 50; i++ {
		_, err = s.tx.Exec(`INSERT INTO users(username, first_name, last_name, role_id) VALUES ($1, $2, $3, $4)`, faker.Username(), faker.FirstName(), faker.LastName(), id)
		if err != nil {
			return err
		}
	}

	err = s.tx.Get(&id, `SELECT id FROM roles WHERE name = 'user'`)
	if err != nil {
		return err
	}

	for i := 0; i < 50; i++ {
		_, err = s.tx.Exec(`INSERT INTO users(username, first_name, last_name, role_id) VALUES ($1, $2, $3, $4)`, faker.Username(), faker.FirstName(), faker.LastName(), id)
		if err != nil {
			return err
		}
	}

	return err
}

// ProductsSeed seeds product data.
func (s SeedTx) ProductsSeed() error {
	for i := 0; i < 100; i++ {
		var err error

		_, err = s.tx.Exec(`INSERT INTO products(name, price) VALUES ($1, $2)`, faker.Word(), rand.Float32())
		if err != nil {
			return err
		}
	}

	return nil
}

// SeedTxWithError struct used to validate the ability to rollback on errors.
type SeedTxWithError struct {
	tx *sqlx.Tx
}

// NewSeedTxWithError returns a seed object with a database tx.
func NewSeedTxWithError(tx *sqlx.Tx) SeedTxWithError {
	return SeedTxWithError{
		tx: tx,
	}
}

// RolesSeed seeds roles data.
func (s SeedTxWithError) RolesSeed() error {
	var err error

	// intentionally wrong sql statement.
	_, err = s.tx.Exec(`INSERT INTO roles(name, last_name) VALUES ($1)`, "admin")
	if err != nil {
		return err
	}

	_, err = s.tx.Exec(`INSERT INTO roles(name) VALUES ($1)`, "user")
	if err != nil {
		return err
	}

	return err
}
