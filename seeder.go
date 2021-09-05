// Package seeder provides a set of tools to seed databases.
package seeder

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
)

// Execute executes all the methods given a struct.
func Execute(s interface{}, seedMethodNames ...string) error {
	sType := reflect.TypeOf(s)
	sKind := sType.Kind()

	if sKind != reflect.Struct {
		return fmt.Errorf("receive a %s instead of a struct", sType.String())
	}

	// Execute all seeders if no method name is given.
	if len(seedMethodNames) == 0 {
		// We are looping over the method on a Seed struct.
		for i := 0; i < sType.NumMethod(); i++ {
			// Get the method in the current iteration.
			method := sType.Method(i)
			// Execute seeder.
			if err := seed(s, method.Name); err != nil {
				return err
			}
		}
	}

	// Execute only the given method names
	for _, item := range seedMethodNames {
		if err := seed(s, item); err != nil {
			return err
		}
	}

	return nil
}

// ExecuteFunc execute a function passing a functionat receives a database
// pool of connections.
func ExecuteFunc(db *sql.DB, f func(*sql.DB) error) error {
	return f(db)
}

func seed(s interface{}, methodName string) error {
	m := reflect.ValueOf(s).MethodByName(methodName)
	if !m.IsValid() {
		return fmt.Errorf("invalid method name: %s", methodName)
	}

	log.Println("seeding ", methodName, "...")
	m.Call(nil)
	log.Println("seed ", methodName, "succeed")

	return nil
}
