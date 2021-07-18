# readme

Database seeds. CLI and Golang library.


![unit tests](https://github.com/danvergara/seeder/actions/workflows/test.yaml/badge.svg)
![linters](https://github.com/danvergara/seeder/actions/workflows/lint.yaml/badge.svg)

## Overview

Seeder is an agnostic cli and library intended to seeds databases using Go code.

## Features

* Driver agnostic (you can choose whatever database driver you want)
* sql builder or ORM agnostic (you can run your seeds no matter what library you choose)
* Flexibility (you can use seeder either as cli tool or as library and included it in your codebase)

## Installation

CLI:

- [Precompiled binaries](https://github.com/danvergara/seeder/releases) for supported
operating systems are available.


Library:

```sh
$ go get github.com/danvergara/seeder
```

## Help



```
Seeder is a ClI tool and Golang library that helps to
seeds databases using golang code. ORM or SQL driver agnostic.

Usage:
  seeder [flags]
  seeder [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  version     A brief description of your command

Flags:
  -h, --help          help for seeder
  -p, --path string    (default "/home/danvergara/Documents/goprojects/seeder/db")

Use "seeder [command] --help" for more information about a command.

```

## Usage

`seeder` is simple because it's flexible, actually we could define seeder as a tool that runs all the methods attached to a single entity at once.

Here is the proposed pattern. Define an object and name it whatever you want. Then, add it a field called `db` with the database connection (pool of connections) from your favorite database library. In the example we chose `slqx`.

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

Now you can insert the rows you need using go code and the database library you prefer. In this example we use [faker](https://github.com/bxcodec/faker) to generate random data.

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

There is two options to run the seeds:


1. Run or compile the main files as usual.

2. Run the cli utility:

```sh
$ seeder --path path/to/main.go
```

You can skip the --path flag is yout main is located at `db/main.go`. The recommended project structure for seeds is the following (this isn't set in stone, suggestions are welcome):

```
 └── db
    ├── main.go
    └── seeds
        ├── products.go
        ├── roles.go
        ├── seeds.go
        └── users.go
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
