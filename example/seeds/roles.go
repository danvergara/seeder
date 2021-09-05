package seeds

import "log"

// RolesSeed seeds roles data.
func (s Seed) RolesSeed() {
	var err error

	_, err = s.db.Exec(`INSERT INTO roles(name) VALUES ($1)`, "admin")
	if err != nil {
		log.Fatalf("error seeding roles: %v", err)
	}

	_, err = s.db.Exec(`INSERT INTO roles(name) VALUES ($1)`, "user")
	if err != nil {
		log.Fatalf("error seeding roles: %v", err)
	}
}
