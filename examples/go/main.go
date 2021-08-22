package main

import (
	"github.com/brandonkal/pulumi-command/sdk/go/command"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		exec, err := command.NewCommand(ctx, "demo", &command.CommandArgs{
			Create: command.CmdArgs{
				Command: pulumi.ToStringArray([]string{"bash", "-c", "echo 'mytest111' > /tmp/mytest.txt"}),
			},
			Diff: command.CmdArgs{
				Command: pulumi.ToStringArray([]string{"bash", "-c", "test 'mytest111' != \"$(cat /tmp/mytest.txt)\""}),
			},
			Delete: command.CmdArgs{
				Command: pulumi.ToStringArray([]string{"rm", "/tmp/mytest.txt"}),
			},
			Update: command.CmdArgs{
				Command: pulumi.ToStringArray([]string{"bash", "-c", "echo 'mytest111' > /tmp/mytest.txt"}),
			},
			Read: command.CmdArgs{
				Command: pulumi.ToStringArray([]string{"cat", "/tmp/mytest.txt"}),
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("stdout", exec.Stdout)
		ctx.Export("stderr", exec.Stderr)

		return nil
	})
}
