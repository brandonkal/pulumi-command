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
        var test = new Command("testCommand", new CommandSet
        {
            Diff = new CommandSet.CommandArgs
            {
                Command =
                {
                    "bash",
                    "-c",
                    "test 'mytest111' != \"$(cat /tmp/mytest.txt)\""
                }
            },
            Create = new CommandSet.CommandArgs
            {
                Command =
                {
                    "bash",
                    "-c",
                    "echo 'mytest111' > /tmp/mytest.txt"
                }
            },
            Update = new CommandSet.CommandArgs
            {
                Command =
                {
                    "bash",
                    "-c",
                    "echo $VAR > /tmp/mytest.txt"
                },
                Environment =
                {
                    {"VAR", "mytest111"}
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
