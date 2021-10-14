seeder ![unit tests](https://github.com/danvergara/seeder/actions/workflows/test.yaml/badge.svg) ![linters](https://github.com/danvergara/seeder/actions/workflows/lint.yaml/badge.svg) [![GitHub Release](https://img.shields.io/github/release/danvergara/seeder.svg)](https://github.com/danvergara/seeder/releases)
=================
<p align="center">
  <img style="float: right;" src="assets/gopher-seeder.png" alt="Seeder logo"/  width=200>
</p>

__Database seeds. CLI and Golang library.__

## Overview

Seeder is an agnostic cli and library intended to seeds databases using Go code.

## Features

* Driver agnostic (you can choose whatever database driver you want)
* sql builder or ORM agnostic (you can run your seeds no matter what library you choose)
* Flexibility (you can use seeder either as cli tool or as library and included it in your codebase)

## Installation

### CLI:
### Homebrew

It works with Linux, too.

```
$ brew install danvergara/tools/seeder
```

Or

```
$ brew tap danvergara/tools
$ brew install seeder
```

### Binary Release (Linux/OSX/Windows)
You can manually download a binary release from [the release page](https://github.com/danvergara/seeder/releases).

Automated install/update, don't forget to always verify what you're piping into bash:

```sh
curl https://raw.githubusercontent.com/danvergara/seeder/master/scripts/install_update_linux.sh | bash
```

The script installs downloaded binary to `/usr/local/bin` directory by default, but it can be changed by setting `DIR` environment variable.

### Library:

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
  -p, --path string    (default "path/to/db")

Use "seeder [command] --help" for more information about a command.

```

## Usage

`seeder` is simple because it's flexible, actually we could define seeder as a tool that runs all the methods attached to a single entity at once.

Here is the proposed pattern. Define an object and name it whatever you want. Then, add it a field called `db` with the database connection (pool of connections) from your favorite database library. In the example we chose `slqx`.

## Execute function

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

## ExecuteFunc function 

In case you want to direclty work with an instance of `sql.DB` from `database/sql`, you can use `ExecuteFunc` which allows you to pass one or more functions to the `ExecuteFunc` function, along with a pointer to an instance of `sql.DB`.

The functions you want to use to seed the database are required to have following signature:

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

This time, you can pass the function defined previously to `ExecuteFunc`, along with your database connection.

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

## Instructions to run the seeds

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

### Docker Usage

```sh
$ docker run -v "$(pwd):/seeder" seeder --path path/to/main.go
```

If you neeed set to up environment variables to connect with the database:

```sh
$ docker run -v "$(pwd):/seeder" --network host -e DB_HOST='localhost' -e DB_USER='postgres' -e DB_PASSWORD='password' -e DB_NAME='users' -e DB_PORT='5432' -e DB_DRIVER='postgres' seeder --path path/to/main.go
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
