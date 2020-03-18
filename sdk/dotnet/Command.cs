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

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;

namespace Pulumi.Command
{
  public partial class Command : Pulumi.CustomResource
  {
        /// <summary>
        /// stdout of the command
        /// </summary>
        [Output("stdout")]
        public Output<string?> StdOut { get; private set; } = null!;

        /// <summary>
        /// stderr of the command
        /// </summary>
        [Output("stderr")]
        public Output<string?> StdErr { get; private set; } = null!;

        /// <summary>
        /// Create a Command resource with the given unique name, arguments, and options.
        /// </summary>
        ///
        /// <param name="name">The unique name of the resource</param>
        /// <param name="args">The arguments used to populate this resource's properties</param>
        /// <param name="options">A bag of options that control this resource's behavior</param>
        public Command(string name, CommandSet args, CustomResourceOptions? options = null)
            : base("command:v1:exec", name, args ?? ResourceArgs.Empty, MakeResourceOptions(options, ""))
        {

          if (args == null){
            throw new ArgumentNullException(nameof(args));
          }

          if (args.Create == null){
            throw new ArgumentNullException(nameof(args.Create));
          }
        }

        private static CustomResourceOptions MakeResourceOptions(CustomResourceOptions? options, Input<string>? id)
        {
            var defaultOptions = new CustomResourceOptions
            {
                Version = Utilities.Version,
            };
            var merged = CustomResourceOptions.Merge(defaultOptions, options);
            // Override the ID if one was specified for consistency with other language SDKs.
            merged.Id = id ?? merged.Id;
            return merged;
        }
  }

  public sealed class CommandSet : Pulumi.ResourceArgs
  {
        /// <summary>
        /// diff
        /// </summary>
        [Input("diff")]
        public Input<CommandArgs>? Diff { get; set; }

        /// <summary>
        /// compare
        /// </summary>
        [Input("compare")]
        public Input<CommandArgs>? Compare { get; set; }

        /// <summary>
        /// create
        /// </summary>
        [Input("create")]
        public Input<CommandArgs>? Create { get; set; }

        /// <summary>
        /// read
        /// </summary>
        [Input("read")]
        public Input<CommandArgs>? Read { get; set; }

        /// <summary>
        /// update
        /// </summary>
        [Input("update")]
        public Input<CommandArgs>? Update { get; set; }

        /// <summary>
        /// delete
        /// </summary>
        [Input("delete")]
        public Input<CommandArgs>? Delete { get; set; }

        public sealed class CommandArgs : Pulumi.ResourceArgs
        {
          [Input("command")]
          private InputList<string>? _command;

          /// <summary>
          /// Specify the command to run as an array of arguments (list)
          /// </summary>
          public InputList<string> Command
          {
              get => _command ?? (_command = new InputList<string>());
              set => _command = value;
          }

          /// <summary>
          /// Pass the stdin to a command (string)
          /// </summary>
          [Input("stdin")]
          public Input<string>? Stdin { get; set; }

          [Input("environment")]
          private InputMap<object>? _environment;

          /// <summary>
          /// Set environment variables for the running command (list)
          /// </summary>
          public InputMap<object> Environment
          {
              get => _environment ?? (_environment = new InputMap<object>());
              set => _environment = value;
          }
        }
  }
}
