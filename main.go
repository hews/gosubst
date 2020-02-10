package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// GlobalContext represents the values that will be available at the
// top level (ie "$.") in the template. This is where .Proc and .Debug
// come from.
type GlobalContext struct {
	Proc  ProcessDetails
	Debug bool
}

// ProcessDetails are just a grab bag of things we may want to know and
// it'd be nice to have a simple interface for. Most are just what the
// Go package "os" offer up simply, but User, Shell, Term, Path and PWD
// are also pulled from the environment as local shell vars.
type ProcessDetails struct {
	PID           int
	PPID          int
	UID           int
	GID           int
	CWD           string
	Hostname      string
	Executable    string
	TempDir       string
	UserCacheDir  string
	UserConfigDir string
	UserHomeDir   string
	User          string
	Shell         string
	Term          string
	Path          string
	PWD           string
}

// Allow us to use log.Fatalf w/o timestamps, and to test output.
var elog = log.New(os.Stderr, "gosubst: ", 0)
var olog = log.New(os.Stdout, "", 0)

func main() {
	var (
		doExpand   = true
		doTemplate = true
		debug      = false
	)

	// NOTE: the "flags" package is ugly, and this is simple.
	for _, arg := range os.Args[1:] {
		if arg == "-V" || arg == "--version" {
			PrintVersion(olog)
			os.Exit(0)
		}
		if arg == "-h" || arg == "--help" {
			PrintHelp(olog)
			os.Exit(0)
		}
		if arg == "-t" || arg == "--template-only" {
			if !doTemplate {
				elog.Fatalf("invalid options: must expand or template")
			}
			doExpand = false
		}
		if arg == "-e" || arg == "--expand-only" {
			if !doExpand {
				elog.Fatalf("invalid options: must expand or template")
			}
			doTemplate = false
		}
		if arg == "--debug" {
			debug = true
		}
	}

	// Check the current mode of STDIN.
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(os.Stdin)

	// Slurp up whatever has been piped if it's hanging out in STDIN...
	if (info.Mode() & os.ModeCharDevice) == 0 {
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		output, err := Template(string(bytes), doExpand, doTemplate, debug)
		if err != nil {
			elog.Fatalf("input is invalid: %s\n", err)
		}
		fmt.Print(output)
		os.Exit(0)
	}

	// ... otherwise run in interactive mode, ie terminal input.
	for {
		str, err := reader.ReadString('\n')
		if err == io.EOF {
			os.Exit(0)
		}
		if err != nil {
			panic(err)
		}
		output, err := Template(str, doExpand, doTemplate, debug)
		if err != nil {
			elog.Fatalf("input is invalid: %s\n", err)
		}
		fmt.Print(output)
	}
}

// Process gathers the basic .Proc values.
func Process() ProcessDetails {
	return ProcessDetails{
		PID:           os.Getpid(),
		PPID:          os.Getppid(),
		UID:           os.Getuid(),
		GID:           os.Getgid(),
		CWD:           must(os.Getwd()),
		Hostname:      must(os.Hostname()),
		Executable:    must(os.Executable()),
		TempDir:       os.TempDir(),
		UserCacheDir:  must(os.UserCacheDir()),
		UserConfigDir: must(os.UserConfigDir()),
		UserHomeDir:   must(os.UserHomeDir()),
		User:          os.Getenv("USER"),
		Shell:         os.Getenv("SHELL"),
		Term:          os.Getenv("TERM"),
		Path:          os.Getenv("PATH"),
		PWD:           os.Getenv("PWD"),
	}
}

// Template actually runs the templating mechanisms over input, returning
// the result if no errors are encountered.
func Template(input string, doExpand, doTemplate, debug bool) (string, error) {
	var buf bytes.Buffer
	var str string

	// Expand env vars in the input.
	if doExpand {
		str = Expand(input, os.Getenv)
	} else {
		str = input
	}

	// Compile and then execute the input as a Go template, including the
	// functions from Sprig (and sh()).
	if doTemplate {
		tmpl, err := template.New("<stdin>").
			Funcs(sprig.TxtFuncMap()).
			Funcs(FuncMap()).
			Parse(str)
		if err != nil {
			return "", err
		}
		err = tmpl.Execute(&buf, &GlobalContext{
			Proc:  Process(),
			Debug: debug,
		})
		if err != nil {
			return "", err
		}
		str = buf.String()
	}

	return str, nil
}

func must(str string, err error) string {
	if err != nil {
		panic(err)
	}
	return str
}
