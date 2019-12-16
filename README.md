# pulumi-resource-command

This is a simple Pulumi provider that allows one to run arbitrary commands and treat their outputs as a resource. With this, anything can be done in a Pulumi program.

It is important to ensure that the output of a command is deterministic. If it is not, use the diff command to ensure the net results are deterministic. The output of the update and and create commands should remain the same (not just the command effects). See the examples for usage details.

## Install

After installing the node module into your pulumi project, you will need to install the go binary globally.
You may place it in your PATH or place it in the plugins directory:

```sh
cd cmd/pulumi-resource-command
mkdir -p ~/.pulumi/plugins/resource-command-v1.0.4/
go build && mv pulumi-resource-command ~/.pulumi/plugins/resource-command-v1.0.4/
```

## Attribution

Thank you to [Luke Hoban @lukehoban](https://github.com/lukehoban) for his help answering my Pulumi questions on Slack.

© Brandon Kalinowski. Apache-2.0.
