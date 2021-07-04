package seeder_test

import (
	"testing"

	"github.com/danvergara/seeder/pkg/seeder"
)

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
