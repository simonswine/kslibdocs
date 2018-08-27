# kslibdocs

Generate a reference docs website for ksonnet-lib

## Installation

This is an early version of `kslibdocs`. To use it, you build it from source:

```sh
$ mkdir -p ${GOPATH}/src/github.com/ksonnet
$ cd ${GOPATH}/src/github.com/ksonnet
$ git clone https://github.com/ksonnet/kslibdocs
$ go install github.com/ksonnet/kslibdocs/cmd/kslibdocgen
```

## Usage

Build reference website HTML

```sh
$ kslibdocgen --k8sLib=/path/to/k8s.libsonnet -outDir=/tmp/build
```

Hack on HTML templates

```sh
$ hack/init-template.sh /tmp/custom-template
$ <...> hack on the template
$ kslibdocgen --k8sLib=/path/to/k8s.libsonnet -outDir=/tmp/build --templateDir=/tmp/custom-template
```

Update embedded template

```sh
$ hack/update-template.sh /tmp/custom-template
$ go generate ./pkg/...
```
