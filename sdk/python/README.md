# pulumi-resource-command

This is a simple Pulumi provider that allows one to run arbitrary commands and treat their outputs as a resource. With this, anything can be done in a Pulumi program.

It is important to ensure that the output of a command is deterministic. If it is not, use the diff command to ensure the net results are deterministic. The output of the update and and create commands should remain the same (not just the command effects). See the examples for usage details.

## Usage

See [./examples](./examples) folder for examples of plugin usage for available runtimes.

> Note: `python` and `nodejs` runtimes will pull required plugin binaries automatically, for `dotnet` and `go` runtimes check [Installation](#Installation) instruction below

## Installation

Find available versions on [releases](https://github.com/brandonkal/pulumi-command/releases) page and install prebuild plugin with this command:
```sh
pulumi plugin install resource command v<version> --server https://github.com/brandonkal/pulumi-command/releases/download/v<version>/
```

To build and install plugin from source you can do this:

1. Checkout this repo
2. Run this commands:

```sh
make provider
make install
```

## Developing

### Pre-requisites

Install the `pulumictl` cli from the [releases](https://github.com/pulumi/pulumictl/releases) page or follow the [install instructions](https://github.com/pulumi/pulumictl#installation)

> NB: Usage of `pulumictl` is optional. If not using it, hard code the version in the [Makefile](Makefile) of when building explicitly pass version as `VERSION=0.0.1 make build`

### Build and Test

```bash
# build and install the resource provider plugin
$ make build install
# test
$ cd examples/simple
$ npm install
$ pulumi stack init test
$ pulumi up
```

## Attribution

Thank you to [Luke Hoban](https://github.com/lukehoban) for his help answering my Pulumi questions on Slack.

© Brandon Kalinowski. Apache-2.0.
