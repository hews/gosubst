package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gotest.tools/v3/golden"

	gosubst "github.com/hews/gosubst"
)

var templateTests = []struct {
	file             string
	expand, template bool
}{
	{"manifest.yaml", true, true},
	{"manifest.yaml", false, true},
	{"manifest.yaml", true, false},
}

func setEnvironment() {
	home := os.Getenv("HOME") // Panics if this is lost.
	os.Clearenv()
	os.Setenv("HOME", home)
	os.Setenv("APP_NAME", "nginx")
}

func TestTemplate(t *testing.T) {
	setEnvironment()

	for _, test := range templateTests {
		output, err := gosubst.Template(contents(test.file), test.expand, test.template)

		if err != nil {
			t.Errorf(
				"Template(<%s>,%t,%t) returned error %q; expected nil",
				test.file,
				test.expand,
				test.template,
				err,
			)
		}

		golden.Assert(t, output, goldenName(test.file, test.expand, test.template))
	}
}

func contents(filename string) string {
	byt, err := ioutil.ReadFile("./examples/" + filename)
	if err != nil {
		panic(err)
	}
	return string(byt)
}

func goldenName(filename string, expand, template bool) string {
	var expandName, templateName string

	if expand {
		expandName = "+expand"
	}
	if template {
		templateName = "+template"
	}

	return fmt.Sprintf("%s%s%s.golden", filename, expandName, templateName)
}
