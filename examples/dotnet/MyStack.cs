using Pulumi;
using Pulumi.Command;
using Pulumi.Serialization;

class MyStack : Stack
{
    public MyStack()
    {
        var test = new Command("testCommand", new CommandSet
        {
            Diff = new CommandSet.CommandArgs
            {
                Command =
                {
                    "true"
                }
            },
            Create = new CommandSet.CommandArgs
            {
                Command =
                {
                    "ls",
                    "-lh"
                }
            },
            Update = new CommandSet.CommandArgs
            {
                Command =
                {
                    "bash",
                    "-c",
                    "echo $VAR"
                },
                Environment =
                {
                    {"VAR", "Hello Pulumi!"}
                }
            },
        });

        Stdout = test.StdOut;
    }

    [Output]
    public Output<string> Stdout { get; set; }
}
