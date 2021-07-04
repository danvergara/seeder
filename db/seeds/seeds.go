// Package seeds seeds the database.
package seeds

import "github.com/jmoiron/sqlx"

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
