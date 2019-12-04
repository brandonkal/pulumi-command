package main

import (
	"github.com/brandonkal/pulumi-command/pkg/provider"
	"github.com/brandonkal/pulumi-command/pkg/version"
)

var providerName = "command"

func main() {
	provider.Serve(providerName, version.Version)
}
