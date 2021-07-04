// Package seeder provides a set of tools to seed databases.
package seeder

import (
	"fmt"
	"log"
	"reflect"
)

// Execute executes all the methods given a struct.
func Execute(s interface{}) error {
	sType := reflect.TypeOf(s)
	sKind := sType.Kind()

	if sKind != reflect.Struct {
		return fmt.Errorf("receive a %s instead of a struct", sType.String())
	}

	for i := 0; i < sType.NumMethod(); i++ {
		method := sType.Method(i)
		if err := seed(s, method.Name); err != nil {
			return err
		}
	}
	return nil
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
