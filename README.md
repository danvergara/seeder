seeder ![unit tests](https://github.com/danvergara/seeder/actions/workflows/test.yaml/badge.svg) ![linters](https://github.com/danvergara/seeder/actions/workflows/lint.yaml/badge.svg) [![GitHub Release](https://img.shields.io/github/release/danvergara/seeder.svg)](https://github.com/danvergara/seeder/releases)
=================
<p align="center">
  <img style="float: right;" src="assets/gopher-seeder.png" alt="Seeder logo"/  width=200>
</p>

__Insert records into a database programmatically.__

## Overview

Seeder is a tool used to insert records into your relational database programmatically.

## Features

* Driver agnostic (you can choose whatever database driver you want)
* sql builder or ORM agnostic (you can run your seeds no matter what library you choose)

## Installation

### CLI:

> **_NOTE:_** The support for the CLI has been deprecated. After a year of actually using the tool, I've realized that this feature is pointless. The user might be better off running their main files by themselves or compiling custom binaries for specific use cases.

### Library:

```sh
$ go get github.com/danvergara/seeder
```

## Usage

The library provides a set of functions as the API:

* Excute
* ExecuteFunc
* ExecuteTxFunc

## Execute

Create an struct and define methods used to insert records into the database on that object.

```go
// db/seeds/seeds.go
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
```

This example uses [faker](https://github.com/bxcodec/faker) to generate random data.

```go
// db/seeds/roles.go
// db/seeds/users.go
// db/seeds/products.go
import (
	"log"
	"math/rand"

	"github.com/bxcodec/faker/v3"
)

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

// ProductsSeed seeds product data.
func (s Seed) ProductsSeed() {
	for i := 0; i < 100; i++ {
		var err error

		_, err = s.db.Exec(`INSERT INTO products(name, price) VALUES ($1, $2)`, faker.Word(), rand.Float32())
		if err != nil {
			log.Fatalf("error seeding products: %v", err)
		}
	}
}
```

Then, instantiate the `Seed` struct. The `Execute` function is gonna access to all the methods attached to `Seed`.

```go
//  db/main.go
import (
	"log"

	"github.com/danvergara/seeder/db/seeds"
	"github.com/danvergara/seeder"
	"github.com/jmoiron/sqlx"

	// postgres driver.
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "postgres-url")
	if err != nil {
		log.Fatalf("error opening a connection with the database %s\n", err)
	}

	s := seeds.NewSeed(db)

	if err := seeder.Execute(s); err != nil {
		log.Fatalf("error seeding the db %s\n", err)
	}
}

```
Unfortunately, due to `Seeder` uses reflection to guess the number and the name of the methods, the execution of methods is sorted in lexicographic order. So, if you chose this approach, make sure the order of the desired execution matches the lexicographic order of the defined methods. There's another way to deal with this limitation:

```go
import (
	"log"
	"math/rand"

	"github.com/bxcodec/faker/v3"
)

func (s Seed) rolesSeed() {
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

func (s Seed) usersSeed() {
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

func (s Seed) productsSeed() {
	for i := 0; i < 100; i++ {
		var err error

		_, err = s.db.Exec(`INSERT INTO products(name, price) VALUES ($1, $2)`, faker.Word(), rand.Float32())
		if err != nil {
			log.Fatalf("error seeding products: %v", err)
		}
	}
}

func (s Seed) PopulateDB() {
	s.rolesSeed()
	s.usersSeed()
	s.productsSeed()
}
```

By making the methods unexported and defining them in a specific order in another exported method, bypassing the limitation imposed by the `reflect` package.

This approach has a problem we recently spotted and which is that if an insertion errors out, the previous insertions can't be rollback. To tackle this problem down, we can use TXs.

```go
// db/seeds/seeds.go
package seeds

import "github.com/jmoiron/sqlx"

// Seed struct.
type Seed struct {
	tx *sqlx.Tx
}

// NewSeed return a Seed with a pool of connection to a dabase.
func NewSeed(tx *sqlx.Tx) Seed {
	return Seed{
		tx: tx,
	}
}
```
Then, handle the tx based on the error value:

```go
//  db/main.go
import (
	"log"

	"github.com/danvergara/seeder/db/seeds"
	"github.com/danvergara/seeder"
	"github.com/jmoiron/sqlx"

	// postgres driver.
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "postgres-url")
	if err != nil {
		log.Fatalf("error opening a connection with the database %s\n", err)
	}

	tx, err := db.Beginx()
	if err != nil {
		log.Fatalf("error creating a tx %s\n", err)
	}

	s := seeds.NewSeed(tx)

	if err := seeder.Execute(s); err != nil {
        tx.Rollback()
	}

    tx.Commit()
}
```


## ExecuteFunc

In case you want to direclty work with an instance of `sql.DB` from `database/sql`, you can use `ExecuteFunc` which allows you to pass one or more functions to the `ExecuteFunc` function, along with a pointer to a `sql.DB` instance.

The functions you use to seed the database need to have the following signature:

```go
func(*sql.DB) error
```

```go
// db/seeds/seeds.go 

import (
	"database/sql"
	"math/rand"

	"github.com/bxcodec/faker/v3"
	"github.com/jmoiron/sqlx"
)

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
```

Now, you can pass the function previously define to `ExecuteFunc`, along with your database connection.

```go
//  db/main.go
import (
	"log"

	"github.com/danvergara/seeder/db/seeds"
	"github.com/danvergara/seeder"
	"github.com/jmoiron/sqlx"

	// postgres driver.
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "postgres-url")
	if err != nil {
		log.Fatalf("error opening a connection with the database %s\n", err)
	}

	// Here's where you pass the functions you want to use to insert new records into the database.
	if err := seeder.ExecuteFunc(db, seeds.PopulateDB); err != nil {
		log.Fatalf("error seeding the db %s\n", err)
	}
}
```

## EexcuteTxFunc

The function recieves a `Tx` as a parameter, along with a list of functions which accept TXs, too.

Try this out:
```go
//  db/main.go
import (
	"log"

	"github.com/danvergara/seeder/db/seeds"
	"github.com/danvergara/seeder"
	"github.com/jmoiron/sqlx"

	// postgres driver.
	_ "github.com/lib/pq"
)

func main() {
	db, err := sqlx.Open("postgres", "postgres-url")
	if err != nil {
		log.Fatalf("error opening a connection with the database %s\n", err)
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("error creating a tx %s\n", err)
	}

	// Here's where you pass the functions you want to use to insert new records into the database.
	if err := seeder.ExecuteTxFunc(tx, seeds.PopulateTx); err != nil {
		log.Fatalf("error seeding the db %s\n", err)
	}
}
```

## Instructions to run the seeds

There is two options to run the seeds:

1. Run the main file:

```sh
$ go run ./example/main.go
```

```
 └── db
    ├── main.go
    └── seeds
        ├── products.go
        ├── roles.go
        ├── seeds.go
        └── users.go
```

2. Compile the project:

```sh
$ cd example && go build && ./example
```

## Contribute

- Fork this repository
- Create a new feature branch for a new functionality or bugfix
- Commit your changes
- Execute test suite
- Push your code and open a new pull request
- Use [issues](https://github.com/danvergara/seeder/issues) for any questions

## License
Apache-2.0 License. See [LICENSE](LICENSE) file for more details.
