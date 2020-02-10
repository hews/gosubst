package main

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"text/template"

	"github.com/spf13/afero"
)

// FuncMap bindles up the functions defined in this file for injection
// into a Go template definition.
func FuncMap() template.FuncMap {
	return template.FuncMap(map[string]interface{}{
		"requiredEnvs":  RequiredEnvs,
		"requiredVals":  RequiredVals,
		"requiredFiles": RequiredFiles,
		"sh":            Sh,
	})
}

// RequiredEnvs raises an error if the given strings do not reference
// defined environmental variables. Inspired by Helm's `required()`.
func RequiredEnvs(envvars ...string) (string, error) {
	for _, envvar := range envvars {
		if _, defined := os.LookupEnv(envvar); !defined {
			return "", fmt.Errorf("required environmental variable missing: ${%s}", envvar)
		}
	}
	return "", nil
}

// RequiredVals raises an error if any of the given values return true
// from Sprig's `empty()` function. Inspired by Helm's `required()`.
func RequiredVals(vals ...interface{}) (string, error) {
	for _, val := range vals {
		if empty(val) {
			return "", fmt.Errorf("required value is empty: %v", val)
		}
	}
	return "", nil
}

// FsBackend is a mockable reference to the filesystem.
var FsBackend = afero.NewOsFs()

// RequiredFiles raises an error if the given strings refer to a path
// for non-existent files. Inspired by Helm's `required()`. Uses afero
// to wrap FS calls for testing.
func RequiredFiles(paths ...string) (string, error) {
	for _, path := range paths {
		stat, err := FsBackend.Stat(path)
		if os.IsNotExist(err) {
			return "", fmt.Errorf("required file missing: %s", path)
		}
		if stat.IsDir() {
			return "", fmt.Errorf("required file missing: %q is a directory", path)
		}
	}
	return "", nil
}

// Sh implements the `sh()` function used in the template to run
// basic shell commands and inject their STDOUT back into the document.
// STDERR output is attached to the err, but then is promptly ignored.
func Sh(cmdstr string) (string, error) {
	out, err := exec.Command("sh", "-c", cmdstr).Output()
	return string(out), err
}

// Copied wholesale from Sprig v3.0.2:
// https://github.com/Masterminds/sprig/blob/3c4c60440a40b962bd9949f90a547bd24566158c/defaults.go#L28-L54
//
// empty returns true if the given value has the zero value for its type.
func empty(given interface{}) bool {
	g := reflect.ValueOf(given)
	if !g.IsValid() {
		return true
	}

	// Basically adapted from text/template.isTrue
	switch g.Kind() {
	default:
		return g.IsNil()
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String:
		return g.Len() == 0
	case reflect.Bool:
		return !g.Bool()
	case reflect.Complex64, reflect.Complex128:
		return g.Complex() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return g.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return g.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return g.Float() == 0
	case reflect.Struct:
		return false
	}
}
