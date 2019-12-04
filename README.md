# pulumi-resource-command

This is a simple Pulumi provider that allows one to run arbitrary commands and treat their outputs as a resource. With this, anything can be done in a Pulumi program.

It is important to ensure that the output of a command is deterministic. If it is not, use the diff command to ensure the net results are deterministic.

(c) Brandon Kalinowski
