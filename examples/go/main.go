package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"

	"github.com/brandonkal/pulumi-command/sdk/go/command"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		conf := config.New(ctx, "")
		dir := conf.Get("dir")
		if dir == "" {
			dir = "/tmp"
		}
		content := conf.Get("content")
		if content == "" {
			content = "mytest111"
		}

		exec, err := command.NewCommand(ctx, "demo", &command.CommandArgs{
			Create: command.CmdArgs{
				Command: pulumi.ToStringArray([]string{"bash", "-c", fmt.Sprintf("echo '%s' > %s/mytest.txt", content, dir)}),
			},
			Diff: command.CmdArgs{
				Command: pulumi.ToStringArray([]string{"bash", "-c", fmt.Sprintf("test '%s' != \"$(cat %s/mytest.txt)\"", content, dir)}),
			},
			Delete: command.CmdArgs{
				Command: pulumi.ToStringArray([]string{"rm", fmt.Sprintf("%s/mytest.txt", dir)}),
			},
			Update: command.CmdArgs{
				Command: pulumi.ToStringArray([]string{"bash", "-c", fmt.Sprintf("echo '%s' > %s/mytest.txt", content, dir)}),
			},
			Read: command.CmdArgs{
				Command: pulumi.ToStringArray([]string{"cat", fmt.Sprintf("%s/mytest.txt", dir)}),
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
