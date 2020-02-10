# `gosubst`

Here's the deal: I just wanted to do some simple command-line templating of configuration files that are used as input to a lot of cloud service tools. YAML files, confs, .envs, lots of container definitions, orchestration manifests, etc. `envsubst` does not get the job done, and I want something that just sits on the command line, is sturdy, and allows me to share configuration publicly (or even internally) without exposing secrets.

That is `gosubst`.

## Installation

Easy:

```
$ go install github.com/hews/gosubst
```

If you don't have a working Go environment, it's only a little harder:

```
$ export VERSION=v0.4.0; export OS=linux; # ... or "darwin" (MacOS), or "windows" ...
$ curl -LJ \
  https://github.com/hews/gosubst/releases/download/${VERSION}/gosubst.${OS} \
  -o gosubst
$ chmod +x gosubst
$ mv gosubst /usr/local/bin
```

> See [releases](https://github.com/hews/gosubst/releases) for what options will work.

However, if you use **MacOS**, it's as easy as can be:

```
$ brew tap hews/tap && brew update
# > ==> Tapping hews/tap
# > ...
$ brew install gosubst
# > ==> Installing gosubst from hews/tap
# > ...
```

## Use

```
$ echo '{{ printf "This is a %s!" "${TEST}" }}' | TEST=test gosubst
# > This is a test!
```

This is meant to be a minimal and dead simple replacement for `envsubst`. Though it does not implement a superset of `envsubst` exactly, it is used for the same purpose, and any file being used with `envsubst` can be very easily changed to work with **`gosubst`**.

Like `envsubst` it reads from STDIN and writes to STDOUT (both with pipes, and in interactive mode), and expands environmental variables in the form `${ENVVAR}`. Unlike envsubst, there is no way to filter or print the variables that are used, and variables in the form `$ENVVAR` are ignored.

After variable expansion, however, comes the fun part! The input is treated like a [Go template][gotemplates], and the context for the calling process is injected into it. This context includes some shell variables, details about the process, and debugging flags.

A full suite of functions is available to use in templating via [Sprig][sprig]! There is also an available function `sh("...")` that hands off to `sh -c '...'`, so that we can nest shell commands into the template (and a few more utility functions on top of that).

**[See the `/examples` directory for examples.](examples)**

<!-- TODO: move to documentation.

### Ok, but really, how do I use it?

Yeah so, it's [`envsubst`][envsubst] but uses [Go templates][gotemplates]. Read about those!

In addition to doing vanilla Go template rendering, the things to know are:

-  The top-level context indcludes two values: `Proc` and `Debug`; `Proc` contains details about the process and shell that initiated the command, and `Debug` identifies if the --debug command line option was passed.
-  We actually expand env vars in the text _BEFORE_ we template it. This means we can use `"value: ${SECRET_VAR}"` just like always. Be careful though when mixing this with Go templating: remember we expand these first!
-  The template loads [Sprig functions][sprig] for fun and profit. See `--version` for information on the version of Sprig used.
-  An extra, very important but possibly world-destroying, function is also added called `sh()`, that in essence spawns a sub-process that runs the given string with `/bin/sh` (assuming a *nix system).

-->

## Some Q&A, in which I lob myself softballs and knock 'em outta the park...

![Donald Rumsfeld, self-satisfied, speaks to a group of people as if they were fools, his soul drenched in blood.](rummy.jpg)

> _Are we concerned about the code running fast and efficiently?_
>
> Not at all. Speed is not of the essence.
>
> _Are we worried about getting or handling errors?_
>
> Nope. If we get an error we just spit it out. 
>
> _Well it seems like this was slapped together. It's not particularly robust, and it's got a bunch of loopholes._
>
> You munge data with the tools you have, not the tools you wish you had. When I have a specific problem I'll fix it. You have a problem? You fix it. Easy-peasy.

---

Copyright (c) 2020 Philip Hughes   
[License](LICENSE.md)

<!-- LINKS -->

[envsubst]:    https://www.gnu.org/software/gettext/manual/html_node/envsubst-Invocation.html
[gotemplates]: https://golang.org/pkg/text/template/
[envvars]:     https://www.gnu.org/software/bash/manual/html_node/Environment.html
[sprig]:       https://github.com/Masterminds/sprig
