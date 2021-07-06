package seeds

import (
	"log"

	"github.com/bxcodec/faker/v3"
)

// UsersSeed seeds roles data.
func (s Seed) UsersSeed() {
	var id int
	var err error

	err = s.db.Get(&id, `SELECT id FROM roles WHERE name = 'admin'`)
	if err != nil {
		log.Fatalf("error querying the roles table: %v", err)
	}

	for i := 0; i < 50; i++ {
		_, err = s.db.Exec(`INSERT INTO users(username, first_name, last_name, role_id) VALUES ($1, $2, $3, $4)`, faker.Username(), faker.FirstName(), faker.LastName(), id)
		if err != nil {
			log.Fatalf("error seeding roles: %v", err)
		}
	}

	err = s.db.Get(&id, `SELECT id FROM roles WHERE name = 'user'`)
	if err != nil {
		log.Fatalf("error querying the roles table: %v", err)
	}

	for i := 0; i < 50; i++ {
		_, err = s.db.Exec(`INSERT INTO users(username, first_name, last_name, role_id) VALUES ($1, $2, $3, $4)`, faker.Username(), faker.FirstName(), faker.LastName(), id)
		if err != nil {
			log.Fatalf("error seeding roles: %v", err)
		}
	}
}
