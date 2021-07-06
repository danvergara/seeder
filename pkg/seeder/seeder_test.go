package seeder_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/danvergara/seeder/db/seeds"
	"github.com/danvergara/seeder/pkg/seeder"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	// postgres driver for testing.
	_ "github.com/lib/pq"
)

var (
	user     string
	password string
	host     string
	port     string
	dbname   string
)

func TestMain(m *testing.M) {
	user = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	host = os.Getenv("DB_HOST")
	port = os.Getenv("DB_PORT")
	dbname = os.Getenv("DB_NAME")

	os.Exit(m.Run())
}

type Foo struct{}

func (f Foo) Bar()       {}
func (f Foo) Greetings() {}

func TestExecute(t *testing.T) {
	f := Foo{}

	t.Log("Given the need to test the Execute function.")
	{
		err := seeder.Execute(f)
		if err != nil {
			t.Errorf("error calling Execute %s", err)
		}
	}
}

func TestExecuteNoStruct(t *testing.T) {
	s := make(map[string]string)

	t.Log("Given the need to test the Execute function.")
	{
		err := seeder.Execute(s)
		if err != nil {
			t.Logf("should receive an error (%v) with type %T", err, s)
		} else {
			t.Errorf("should not recieve nil as error %v", err)
		}
	}
}

func TestExecuteRealDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname,
	)

	var count int

	db, err := sqlx.Open("postgres", url)
	if err != nil {
		log.Fatalf("error opening a connection with the database: %s\n", err)
	}

	defer db.Close()

	db.MustExec("TRUNCATE users, roles, products CASCADE;")

	s := seeds.NewSeed(db)

	if err := seeder.Execute(s); err != nil {
		t.Errorf("error seeding the db at running tests: %s\n", err)
	}

	err = db.Get(&count, `SELECT COUNT(*) FROM roles`)
	if err != nil {
		t.Errorf("error getting the number of roles in db: %s\n", err)
	}

	assert.Equal(t, 2, count)

	err = db.Get(&count, `SELECT COUNT(*) FROM users`)
	if err != nil {
		t.Errorf("error getting the number of users in db: %s\n", err)
	}

	assert.Equal(t, 100, count)

	err = db.Get(&count, `SELECT COUNT(*) FROM products`)
	if err != nil {
		t.Errorf("error getting the number of products in db: %s\n", err)
	}

	assert.Equal(t, 100, count)
}
