package seeder_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/danvergara/seeder"
	"github.com/danvergara/seeder/example/seeds"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"

	// postgres driver for testing.
	txdb "github.com/DATA-DOG/go-txdb"
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

	txdb.Register("pgsqltx", "postgres", databaseURL())

	os.Exit(m.Run())
}

func databaseURL() string {
	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname,
	)

	return url
}

func prepareDB(t *testing.T) *sqlx.DB {
	db, err := sqlx.Open("postgres", databaseURL())
	if err != nil {
		log.Fatalf("error opening a connection with the database: %s\n", err)
	}

	t.Cleanup(func() {
		db.MustExec("TRUNCATE users, roles, products CASCADE;")
		db.Close()
	})

	db.MustExec("TRUNCATE users, roles, products CASCADE;")

	return db
}

func assertDBContent(t *testing.T, db *sqlx.DB) {
	t.Helper()

	var count int

	if err := db.Get(&count, `SELECT COUNT(*) FROM roles`); err != nil {
		t.Errorf("error getting the number of roles in db: %s\n", err)
	}

	assert.Equal(t, 2, count)

	if err := db.Get(&count, `SELECT COUNT(*) FROM users`); err != nil {
		t.Errorf("error getting the number of users in db: %s\n", err)
	}

	assert.Equal(t, 100, count)

	if err := db.Get(&count, `SELECT COUNT(*) FROM products`); err != nil {
		t.Errorf("error getting the number of products in db: %s\n", err)
	}

	assert.Equal(t, 100, count)
}

func assertDEmptyBContent(t *testing.T, db *sqlx.DB) {
	t.Helper()

	var count int

	if err := db.Get(&count, `SELECT COUNT(*) FROM roles`); err != nil {
		t.Errorf("error getting the number of roles in db: %s\n", err)
	}

	assert.Equal(t, 0, count)

	if err := db.Get(&count, `SELECT COUNT(*) FROM users`); err != nil {
		t.Errorf("error getting the number of users in db: %s\n", err)
	}

	assert.Equal(t, 0, count)

	if err := db.Get(&count, `SELECT COUNT(*) FROM products`); err != nil {
		t.Errorf("error getting the number of products in db: %s\n", err)
	}

	assert.Equal(t, 0, count)
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

	db := prepareDB(t)

	s := seeds.NewSeed(db)

	if err := seeder.Execute(s); err != nil {
		t.Errorf("error seeding the db at running tests: %s\n", err)
	}

	assertDBContent(t, db)
}

func TestExecuteRealDBTxRollback(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	db := prepareDB(t)
	tx, _ := db.Beginx()

	s := seeds.NewSeedTxWithError(tx)

	if err := seeder.Execute(s); err != nil {
		_ = tx.Rollback()
		assertDEmptyBContent(t, db)
		return
	}

	t.Error("the seeding process did not throw an error as expected")
}

func TestExecuteRealDBTx(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	db := prepareDB(t)

	tx, _ := db.Beginx()

	s := seeds.NewSeedTx(tx)

	if err := seeder.Execute(s); err != nil {
		_ = tx.Rollback()
		t.Errorf("error seeding the db at running tests: %s\n", err)
	}

	_ = tx.Commit()

	assertDBContent(t, db)
}

func TestExecuteParallelDB(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	for i := 0; i <= 10; i++ {
		t.Run(fmt.Sprintf("test #%d", i), func(t *testing.T) {
			t.Parallel()

			cName := fmt.Sprintf("connection_%d", time.Now().UnixNano())
			db, _ := sqlx.Open("pgsqltx", cName)

			db.MustExec("TRUNCATE users, roles, products CASCADE;")

			s := seeds.NewSeed(db)

			if err := seeder.Execute(s); err != nil {
				t.Errorf("error seeding the db at running tests: %s\n", err)
			}

			fmt.Printf("test #%d succeed", i)

			assertDBContent(t, db)

			db.Close()
		})
	}
}

func TestExecuteGivenMethodName(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	var count int

	db := prepareDB(t)

	s := seeds.NewSeed(db)

	if err := seeder.Execute(s, "RolesSeed"); err != nil {
		t.Errorf("error seeding the db at running tests: %s\n", err)
	}

	if err := db.Get(&count, `SELECT COUNT(*) FROM roles`); err != nil {
		t.Errorf("error getting the number of roles in db: %s\n", err)
	}

	assert.Equal(t, 2, count)

	if err := db.Get(&count, `SELECT COUNT(*) FROM users`); err != nil {
		t.Errorf("error getting the number of users in db: %s\n", err)
	}

	assert.Equal(t, 0, count)
}

func TestExecuteFunc(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping short mode")
	}

	db := prepareDB(t)

	if err := seeder.ExecuteFunc(db.DB, seeds.PopulateDB); err != nil {
		t.Errorf("error seeding the db at running tests: %s\n", err)
	}

	assertDBContent(t, db)
}
