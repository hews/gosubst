package main_test

import (
	"bytes"
	"log"
	"strings"
	"testing"

	gosubst "github.com/hews/gosubst"
)

var outbuf bytes.Buffer
var captureLogger = log.New(&outbuf, "", 0)

func TestPrintVersion(t *testing.T) {
	outbuf.Reset()

	gosubst.VersionBuild = "abcdefg"
	gosubst.VersionRelease = "v2.0.0-rc1"
	gosubst.VersionSprig = "v3.1.1"

	expected := "gosubst v2.0.0-rc1, build: abcdefg\n  └ using Sprig v3.1.1\n"
	gosubst.PrintVersion(captureLogger)

	actual := outbuf.String()
	if actual != expected {
		t.Errorf("PrintVersion() == %q; expected %q", actual, expected)
	}
}

// Get that coverage UP!
func TestPrintHelp(t *testing.T) {
	outbuf.Reset()

	gosubst.PrintHelp(captureLogger)
	actual := outbuf.String()
	if !strings.Contains(gosubst.HelpText, actual) {
		t.Errorf("PrintHelp() does not print HelpText:\nEXPECTED:%q\nACTUAL:%q", gosubst.HelpText, actual)
	}
}
