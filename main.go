package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// GlobalContext represents the values that will be available at the
// top level (ie "$.") in the template. This is where .Env and .Proc
// come from.
type GlobalContext struct {
	Env  map[string]string
	Proc ProcessDetails
}

// ProcessDetails are just a grab bag of things we may want to know and
// it'd be nice to have a simple interface for. Most are just what the
// Go package "os" offer up simply, but User, Shell, Term, Path are
// also pulled from the environment as local shell vars, and would be
// repeated in .Env.
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
}

// Allow us to use log.Fatalf w/o timestamps, and to test output.
var elog = log.New(os.Stderr, "gosubst: ", 0)
var olog = log.New(os.Stdout, "", 0)

func main() {
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
		Run(string(bytes))
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
		Run(str)
	}
}

// Environment gathers the environment variables for .Env.
func Environment() map[string]string {
	rawenv := os.Environ()
	envmap := make(map[string]string, len(rawenv))

	for _, envvar := range rawenv {
		kv := strings.SplitN(envvar, "=", 2)
		envmap[kv[0]] = kv[1]
	}
	return envmap
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
	}
}

// Sh implements the `sh()` function used in the template to run
// basic shell commands and inject their output back into the document.
func Sh(cmdstr string) string {
	out, err := exec.Command("/bin/sh", "-c", cmdstr).Output()
	if err != nil {
		panic(err)
	}
	return string(out)
}

// Run actually runs the templating mechanisms over input, writing
// output to STDOUT.
func Run(input string) {
	// Expand env vars in the input.
	input = os.ExpandEnv(input)

	// Compile and then execute the input as a Go template, including the
	// functions from Sprig (and sh()).
	tmpl := template.Must(template.New("gosubst.main").
		Funcs(sprig.TxtFuncMap()).
		Funcs(template.FuncMap(map[string]interface{}{
			"sh": Sh,
		})).
		Parse(input))
	err := tmpl.Execute(os.Stdout, &GlobalContext{
		Env:  Environment(),
		Proc: Process(),
	})

	if err != nil {
		panic(err)
	}
}

func must(str string, err error) string {
	if err != nil {
		panic(err)
	}
	return str
}
