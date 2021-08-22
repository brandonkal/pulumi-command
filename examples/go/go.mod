module github.com/brandonkal/pulumi-command/examples/go

go 1.16

replace github.com/brandonkal/pulumi-command/sdk => ../../sdk

require (
	github.com/brandonkal/pulumi-command/sdk v0.0.0-00010101000000-000000000000
	github.com/pulumi/pulumi/sdk/v3 v3.10.1
)
