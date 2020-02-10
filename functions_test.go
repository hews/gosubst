package main_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"text/template"

	"github.com/spf13/afero"

	gosubst "github.com/hews/gosubst"
	"github.com/hews/gosubst/internal/testutils"
)

func TestRequiredEnvs(t *testing.T) {
	recoverEnvironment := testutils.ClearEnvironment(t)
	defer recoverEnvironment()

	os.Setenv("VAR1", "yes")
	os.Setenv("VAR2", "")

	tests := []struct {
		tmp, out string
		err      error
	}{
		{`{{ requiredEnvs "VAR1" }}`, "", nil},
		{`{{ requiredEnvs "VAR2" }}`, "", nil},
		{`{{ requiredEnvs "VAR3" }}`, "", errors.New("required environmental variable missing: ${VAR3}")},
		{`{{ requiredEnvs "VAR1" "VAR2" }}`, "", nil},
		{`{{ requiredEnvs "VAR1" "VAR2" "VAR3" }}`, "", errors.New("required environmental variable missing: ${VAR3}")},
		{`{{ requiredEnvs "VAR1" "VAR3" "VAR2" }}`, "", errors.New("required environmental variable missing: ${VAR3}")},
	}
	for _, test := range tests {
		err := runt(test.tmp, test.out)
		if (err != nil && test.err == nil) || (err == nil && test.err != nil) {
			t.Errorf("%s did not have expected error %q: instead was %q", test.tmp, test.err, err)
		} else if err != nil && test.err != nil &&
			!strings.Contains(err.Error(), test.err.Error()) {
			t.Errorf("%s expected error containing %q: instead was %q", test.tmp, test.err, err)
		}
	}
}

func TestRequiredVals(t *testing.T) {
	dict := map[string]interface{}{"top": map[string]interface{}{"Thing": 1}}
	serr := errors.New("required value is empty")

	tests := []struct {
		tmp, out string
		err      error
	}{
		{`{{ requiredVals 1 }}`, "", nil},
		{`{{ requiredVals 0 }}`, "", serr},
		{`{{ requiredVals 1 1 }}`, "", nil},
		{`{{ requiredVals 1 0 }}`, "", serr},
		{`{{ requiredVals "" }}`, "", serr},
		{`{{ requiredVals 0.0 }}`, "", serr},
		{`{{ requiredVals false }}`, "", serr},
		{`{{ requiredVals 1 "hi" true }}`, "", nil},
		{`{{ requiredVals 1 "" true }}`, "", serr},
		{`{{ requiredVals .top }}`, "", nil},
		{`{{ requiredVals .top.Thing }}`, "", nil},
		{`{{ requiredVals .top.Thing 1 "hello" }}`, "", nil},
		{`{{ requiredVals .top.Thing .top.NoSuchThing }}`, "", serr},
		{`{{ requiredVals .top.NoSuchThing }}`, "", serr},
		{`{{ requiredVals .bottom.NoSuchThing }}`, "", serr},
	}
	for _, test := range tests {
		err := runtv(test.tmp, test.out, dict)
		if (err != nil && test.err == nil) || (err == nil && test.err != nil) {
			t.Errorf("%s did not have expected error %q: instead was %q", test.tmp, test.err, err)
		} else if err != nil && test.err != nil &&
			!strings.Contains(err.Error(), test.err.Error()) {
			t.Errorf("%s expected error containing %q: instead was %q", test.tmp, test.err, err)
		}
	}
}

func TestRequiredFiles(t *testing.T) {
	fs := gosubst.FsBackend
	defer func() {
		gosubst.FsBackend = fs
	}()

	gosubst.FsBackend = afero.NewMemMapFs()
	gosubst.FsBackend.MkdirAll("tmp/im-a-dir", 0755)
	afero.WriteFile(gosubst.FsBackend, "tmp/im-soo-here.conf", []byte("example"), 0644)
	afero.WriteFile(gosubst.FsBackend, "tmp/im-also-here", []byte("example"), 0644)

	tests := []struct {
		tmp, out string
		err      error
	}{
		{`{{ requiredFiles "tmp/im-soo-here.conf" }}`, "", nil},
		{`{{ requiredFiles "tmp/im-not-here" }}`, "", errors.New("required file missing: tmp/im-not-here")},
		{`{{ requiredFiles "tmp/im-a-dir" }}`, "", errors.New("required file missing: \"tmp/im-a-dir\" is a directory")},
		{`{{ requiredFiles "tmp/im-soo-here.conf" "tmp/im-also-here" }}`, "", nil},
		{`{{ requiredFiles "tmp/im-soo-here.conf" "tmp/im-not-here" "tmp/im-also-here" }}`, "", errors.New("required file missing")},
	}
	for _, test := range tests {
		err := runt(test.tmp, test.out)
		if (err != nil && test.err == nil) || (err == nil && test.err != nil) {
			t.Errorf("%s did not have expected error %q: instead was %q", test.tmp, test.err, err)
		} else if err != nil && test.err != nil &&
			!strings.Contains(err.Error(), test.err.Error()) {
			t.Errorf("%s expected error containing %q: instead was %q", test.tmp, test.err, err)
		}
	}
}

func TestSh(t *testing.T) {
	tests := []struct {
		tmp, out string
	}{
		{`{{ sh "" }}`, ""},
		{`{{ sh "echo 'SOME TEXT'" }}`, "SOME TEXT\n"},
		{`{{ sh "printf 'MORE %s' 'TEXT\n'" }}`, "MORE TEXT\n"},
		{`{{ sh "pwd" }}`, must(os.Getwd()) + "\n"},
	}
	for _, test := range tests {
		if err := runt(test.tmp, test.out); err != nil {
			t.Error(err)
		}
	}
}

func TestShErrors(t *testing.T) {
	tests := []struct {
		cmd, out string
		err      error
	}{
		// Are errors.
		{"$(return 1)", "", fmt.Errorf("exit status 1")},
		{"xBBB/argh", "", fmt.Errorf("exit status 127")},
		// Are NOT errors.
		{"", "", nil},
		{"echo 'IM NOT AN ERR' >&2", "", nil},
		{"echo 'IM NOT AN ERR' | tee /dev/stderr", "IM NOT AN ERR\n", nil},
	}
	for _, test := range tests {
		out, err := gosubst.Sh(test.cmd)
		if test.out != out {
			t.Errorf("Sh(%q) == %q; expected %q", test.cmd, out, test.out)
		}
		if (test.err == nil && err != nil) ||
			(test.err != nil && err == nil) ||
			(test.err != nil && err != nil && (test.err.Error() != err.Error())) {
			t.Errorf("Sh(%q) has error %q; expected %q", test.cmd, err, test.err)
		}
	}
}

func must(str string, err error) string {
	if err != nil {
		panic(err)
	}
	return str
}

func runt(tpl, expect string) error {
	return runtv(tpl, expect, map[string]string{})
}

func runtv(tpl, expect string, values interface{}) error {
	t := template.Must(template.New("test").Funcs(gosubst.FuncMap()).Parse(tpl))
	var b bytes.Buffer
	err := t.Execute(&b, values)
	if err != nil {
		return err
	}
	if expect != b.String() {
		return fmt.Errorf("Expected '%s', got '%s'", expect, b.String())
	}
	return nil
}
