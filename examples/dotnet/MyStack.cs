// Pulumi Command Provider .NET SDK
// Copyright 2020, Mitchell Maler.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

using Pulumi;
using Pulumi.Command;
using Pulumi.Serialization;

class MyStack : Stack
{
    public MyStack()
    {
        var config = new Pulumi.Config();
        var dir = config.Get("dir") ?? "/tmp";
        var content = config.Get("content") ?? "mytest111";

        var test = new Command("testCommand", new CommandSet
        {
            Diff = new CommandSet.CommandArgs
            {
                Command =
                {
                    "bash",
                    "-c",
                    $"test '{content}' != \"$(cat {dir}/mytest.txt)\""
                }
            },
            Create = new CommandSet.CommandArgs
            {
                Command =
                {
                    "bash",
                    "-c",
                    $"echo '{content}' > {dir}/mytest.txt"
                }
            },
            Update = new CommandSet.CommandArgs
            {
                Command =
                {
                    "bash",
                    "-c",
                    $"echo $VAR > {dir}/mytest.txt"
                },
                Environment =
                {
                    {"VAR", $"{content}"}
                }
            },
        });

        this.Stdout = test.StdOut;
        this.Stderr = test.StdErr;
    }

    [Output]
    public Output<string> Stdout { get; set; }

    [Output]
    public Output<string> Stderr { get; set; }
}
