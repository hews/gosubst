package main_test

import (
	"testing"

	. "github.com/hews/gosubst"
)

func testGetenv(s string) string {
	switch s {
	case "*":
		return "all the args"
	case "#":
		return "NARGS"
	case "$":
		return "PID"
	case "1":
		return "ARGUMENT1"
	case "HOME":
		return "/usr/gopher"
	case "H":
		return "(Value of H)"
	case "home_1":
		return "/usr/foo"
	case "_":
		return "underscore"
	}
	return ""
}

var expandTests = []struct {
	in, out string
}{
	{"", ""},
	{"$*", "$*"},
	{"${*}", "all the args"},
	{"$$", "$$"},
	{"${$}", "PID"},
	{"$1", "$1"},
	{"${1}", "ARGUMENT1"},
	{"now is the time", "now is the time"},
	{"$home_1", "$home_1"},
	{"${home_1}", "/usr/foo"},
	{"$${home_1}", "${home_1}"},
	{"$HOME", "$HOME"},
	{"${HOME}", "/usr/gopher"},
	{"$${HOME}", "${HOME}"},
	{"${H}OME", "(Value of H)OME"},
	{"$${H}OME", "${H}OME"},
	{"$", "$"},
	{"$}", "$}"},
	{"start$+middle$^end$", "start$+middle$^end$"},
	{"mixed$|bag$$$", "mixed$|bag$$$"},
	{"mixed$|bag${$}$", "mixed$|bagPID$"},
	{"A$$$#$1$H$home_1*B", "A$$$#$1$H$home_1*B"},
	{"A$${$}#${1}${H}${home_1}*B", "A${$}#ARGUMENT1(Value of H)/usr/foo*B"},
	{"Hello {{ printf \"${HOME}\" }}", "Hello {{ printf \"/usr/gopher\" }}"},
	{"Hello {{ printf \"$HOME\" }}", "Hello {{ printf \"$HOME\" }}"},
	{"Hello {{ printf \"$$HOME\" }}", "Hello {{ printf \"$$HOME\" }}"},
	{"Hello {{ printf \"$${HOME}\" }}", "Hello {{ printf \"${HOME}\" }}"},
	// invalid syntax; eat up the characters
	{"${", ""},
	{"${}", ""},
	{"start${+middle}${^end}$", "start$"},
}

func TestExpand(t *testing.T) {
	for _, test := range expandTests {
		result := Expand(test.in, testGetenv)
		if result != test.out {
			t.Errorf("Expand(%q) == %q; expected %q", test.in, result, test.out)
		}
	}
}
