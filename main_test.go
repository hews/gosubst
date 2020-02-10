package main_test

import (
	"reflect"
	"testing"

	gosubst "github.com/hews/gosubst"
)

func TestProcess(t *testing.T) {
	proc := gosubst.Process()

	proc_rv := reflect.ValueOf(proc)
	proc_rt := reflect.TypeOf(proc)
	for i := 0; i < proc_rv.NumField(); i++ {
		if proc_rv.Field(i).Interface() == nil {
			t.Errorf("Process().%s == nil, should have a value", proc_rt.Field(i).Name)
		}
	}
}

func TestOptions(t *testing.T) {
	t.Skip("TODO implementation")
}

func TestStdinModePipe(t *testing.T) {
	t.Skip("TODO implementation")
}

func TestStdinModeInteractive(t *testing.T) {
	t.Skip("TODO implementation")
}
