# kube-gen

![latest 0.1.3](https://img.shields.io/badge/latest-0.1.3-green.svg?style=flat)
[![license](http://img.shields.io/badge/license-BSD-red.svg?style=flat)](https://raw.githubusercontent.com/kylemcc/kube-gen/master/LICENSE)


`kube-gen` is a template rendering tool that generates files based on Kubernetes events and metadata. It is based on and inspired by Jason Wilder's excellent [docker-gen](https://github.com/jwilder/docker-gen).

## Installation

#### Binary Download

See the [Releases](https://github.com/kylemcc/kube-gen/releases) page

#### Docker Container

`kube-gen` can be run as a standalone container, or bundled in a container with other applications.

TODO: create official images


## Usage

When run with no arguments (or with `-h/-help`), `kube-gen` prints the following usage message.

```shell
$ kube-gen
Usage: kube-gen [options] <template> [<output>]

Render templates using Kubernetes metadata and events

Options:
  -host string
    	 (default "http://localhost:8001")
  -interval int

  -log-cmd
    	log the output of the pre/post commands (default true)
  -overwrite
    	overwrite the output file if it exists (default true)
  -post-cmd string
    	command to run after template generation in complete
  -pre-cmd string
    	command to run before template generation
  -quiet
    	when set to true, nothing is logged
  -type value
    	types of resources to pull [pods, services, endpoints] - May be specified multiple times. If not specified, all types will be returned
  -version
    	display version information
  -wait string
    	<minimum>[:<maximum>] - the minimum and optional maximum time to wait after an event fires.E.g.: 500ms:5s
  -watch
    	watch for new events

Arguments:
  template: path or URL of the template file to render, or - to read from STDIN
  output: (Optional) path to write the rendered content. If not specified,
          rendered content is printed to STDOUT. By default, this file will
          be overwritten if it exists. Use -overwrite=false to return an
          error instead
```

#### Authentication / Connecting to the Kubernetes API
Currently, `kube-gen` does not support authentication, or reading the `~/.kube/config` file. These features are forthcoming. In the mean time, the `-host` parameter should either be set to the value of `kube-apiserver`'s `--insecure-bind-address` / `--insecure-port` values. Alternatively, as reflected by the default value, `kube-gen` can be run alongside `kubectl proxy` - which, by default, listens on `127.0.0.1:8001`.

#### Watching for changes

The `-watch` flag configures `kube-gen` to watch the API for changes to `Services`, `Pods`, and `Endpoints` (support for other types is forthcoming). This mode is useul when combined with the `-pre-cmd`, `-post-cmd`, and `-wait` parameters.

## Template Language

`kube-gen` supports templates written in Go`s [text/template](https://golang.org/pkg/text/template/) language. It supports all of the [built in](https://golang.org/pkg/text/template/#hdr-Functions) functions, as well as numerous custom functions described below. Many of the custom functions (and the documentation for those functions) have been borrowed from [docker-gen](https://github.com/jwilder/docker-gen). Those functions, along with the accompanying License and Copyright are located in the [dockergen_template_functions.go](https://github.com/kylemcc/kube-gen/blob/master/dockergen_template_functions.go) source file.
