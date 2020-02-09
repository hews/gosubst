package main

import (
	"log"
	"strings"
)

// Version information, injected at build-time.
var (
	VersionBuild   string
	VersionRelease string
	VersionSprig   string
)

// HelpText is the poor man's man for the CLI.
var HelpText = `
Usage: gosubst [OPTION]

Substitutes the values of environment variables.

Options:
  -e, --expand-only           skip the Go templating pass
  -t, --template-only         skip env variable expansion pass
      --debug                 set the template context .Debug val true
  -h, --help                  display this help and exit
  -V, --version               output version information and exit

When gosubst is invoked standard input is copied to standard output,
with references to environment variables of the form ${VARIABLE}
being replaced with the corresponding values first (as in ` + "`envsubst`" + `), and
then passed through the Go templating engine.

For the Go template, the global context includes the environment as .Env,
information about the currently running process as .Proc, and the command
line boolean option --debug as .Debug. Also included in the template are
the suite of Sprig <http://masterminds.github.io/sprig/> functions and a
special ` + "`sh()`" + ` function that evals the given string with` + "`sh -c '...'`" + `.
Use sh at your own peril!

For more information, email <p+gosubst@hews.co>, or visit the project page
at <https://github.com/hews/gosubst>.
`

// PrintVersion prints the version information.
func PrintVersion(log *log.Logger) {
	log.Printf(
		"gosubst %s, build: %s\n  └ using Sprig %s\n",
		VersionRelease,
		VersionBuild,
		VersionSprig,
	)
}

// PrintHelp prints help information.
func PrintHelp(log *log.Logger) {
	log.Println(strings.TrimSpace(HelpText))
}
