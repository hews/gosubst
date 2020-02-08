# `gosubst`

Here's the deal: I want to just do some simple command-line templating of configuration files that are used as input to a lot of cloud service tools. YAML files, confs, .envs, lots of container definitions, orchestration manifests, etc. `envsubst` does not get the job done, and I want something that just sits on the command line and is sturdy.

That is `gosubst`.

## Installation

<!-- 
Figure this all out:

Easy:

```
$ go install github.com/hews/gosubst
```

Little harder:
-->

```
$ export VERSION=v0.1.0; export OS=linux;
$ curl -LJ \
  https://github.com/hews/gosubst/releases/download/${VERSION}/gosubst.${OS} \
  -o gosubst
$ chmod +x gosubst
$ mv gosubst /usr/local/bin
```

See [releases](https://github.com/hews/gosubst/releases) for what options will work.

<!--
_Dreams do come true..._

Little easier (MacOS):

```
$ brew update && brew install gosubst
```
-->

## Use

```
$ echo 'This is a ${TEST}!' | TEST=test gosubst
```

This is meant to be minimal and dead simple, so it accepts _NO*_ options. It reads from STDIN and writes to STDOUT, and is meant to replace `envsubst`, which has only one option `-v`/`--variables` to list the variables in some text, but that doesn't really track with what we're doing, so that's that.

And there you go! [See the `/examples` directory for examples.](examples)

> __*_ – it does accept `-h`/`--help` and `-V`/`--version`, but let's ignore that for now._

### Ok, but really, how do I use it?

Yeah so, it's [`envsubst`][envsubst] but uses [Go templates][gotemplates]. Read about those!

In addition to doing vanilla Go template rendering, the things to know are:

-  The top-level context indcludes two values: `Proc` and `Env`; `Proc` contains details about the process and shell that initiated the command, and `Env` contains a copy of the [shell's environmental variables][envvars] as a string map.
-  Becuase typing `"value: {{ index .Env "SECRET_VAR" }}"` is a big headache, we actually run `envsubst` on the text _BEFORE_ we template it. This means we can use `"value: ${SECRET_VAR}"` just like always. Be careful though when mixing this with Go templating: remember we parse these first!
-  The template loads [Sprig functions][sprig] for fun and profit. See `--version` for information on the version of Sprig used.
-  An extra, very important but possibly world-destroying, function is also added called `sh()`, that in essence spawns a sub-process that runs the given string with `/bin/sh` (assuming a *nix system).

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
