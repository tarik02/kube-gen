# kube-gen

[![license](http://img.shields.io/badge/license-BSD-red.svg?style=flat)](https://raw.githubusercontent.com/tarik02/kube-gen/master/LICENSE)
[![Check Commit](https://github.com/tarik02/kube-gen/actions/workflows/check.yml/badge.svg)](https://github.com/tarik02/kube-gen/actions/workflows/check.yml)


`kube-gen` is a template rendering tool that generates files and executes commands based on Kubernetes events and metadata.

## Installation

#### Binary Download

See the [Releases](https://github.com/tarik02/kube-gen/releases) page

#### Docker Container

`kube-gen` can be run as a standalone container, or bundled in a container with other applications.

Images are available on [Github Package Registry](https://github.com/tarik02/kube-gen/pkgs/container/kube-gen):

```sh
$ docker run ghcr.io/tarik02/kube-gen ...
```


## Usage

When run with no arguments (or with `-h/-help`), `kube-gen` prints the following usage message.

```shell
$ kube-gen
Usage: kube-gen [options] <template> [<output>]

Render templates using Kubernetes metadata and events

Options:
  -host string
        If not set will use kubeconfig. If using proxy - set it to http://localhost:8001
  -interval int

  -kubeconfig string
        (optional) absolute path to the kubeconfig file (default "/Users/kyle/.kube/config")
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
By default, `kube-gen` will look for a kubeconfig file at `$HOME/.kube/config`. A different kubeconfig file may be specified by using the `-kubeconfig` flag. Alternatively, `kube-gen` provides a `-host` flag that, if set, will supersede the `-kubeconfig`. The `-host` flag is best paired with `kubectl proxy`, which listens on 127.0.0.1:8001 by default. The `-host` flag may also be set to the value of `kube-apiserver`'s `--insecure-bind-address` / `--insecure-port`.

#### Watching for changes

The `-watch` flag configures `kube-gen` to watch the API for changes to `Services`, `Pods`, and `Endpoints` (support for other types is forthcoming). This mode is useul when combined with the `-pre-cmd`, `-post-cmd`, and `-wait` parameters.

## Template Language

- `kube-gen` supports templates written in Go`s [text/template](https://golang.org/pkg/text/template/) language.
- it supports all of the [built in](https://golang.org/pkg/text/template/#hdr-Functions) functions.
- it supports all of the [sprig](http://masterminds.github.io/sprig/) functions.
- it supports subset of [docker-gen](https://github.com/jwilder/docker-gen) functions. Implementation, along with the accompanying License and Copyright are located in the [dockergen_template_functions.go](https://github.com/tarik02/kube-gen/blob/master/dockergen_template_functions.go) source file.
