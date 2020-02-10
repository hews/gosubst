package main_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gotest.tools/v3/golden"

	gosubst "github.com/hews/gosubst"
	"github.com/hews/gosubst/internal/testutils"
)

var templateTests = []struct {
	file                    string
	expand, template, debug bool
}{
	{"manifest.yaml", true, true, false},
	{"manifest.yaml", false, true, false},
	{"manifest.yaml", true, false, false},
}

func TestTemplate(t *testing.T) {
	resetEnvirnonment := testutils.ClearEnvironment(t)
	defer resetEnvirnonment()

	os.Setenv("APP_NAME", "nginx")

	for _, test := range templateTests {
		output, err :=
			gosubst.Template(contents(test.file), test.expand, test.template, test.debug)

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
