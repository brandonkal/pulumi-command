// Pulumi Command Provider Node SDK
// Copyright 2019, Brandon Kalinowski.
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

import * as pulumi from '@pulumi/pulumi'

export interface Cmd {
  /** Specifiy the command to run as an array of arguments */
  readonly command: pulumi.Input<string[]>
  /** Pass the stdin to a command */
  readonly stdin?: pulumi.Input<string>
  /** Set environment variables for the running command */
  readonly environment?: pulumi.Input<Record<string, string>>
}

export interface CommandSet {
  /** Specify a command to run to diff the resource.
   *
   * Exit 0 to run update.
   * Exit with a non-zero value or omit to disable update.
   * Hint: an easy method to always run update is to set diff to `['true']` */
  readonly diff?: pulumi.Input<Cmd>
  /** Define a command to create a resource. */
  readonly create: pulumi.Input<Cmd>
  /** Define a command to create read the resource. */
  readonly read?: pulumi.Input<Cmd>
  /** If unspecified, create definition will be used. Define to provide an alternate update command. */
  readonly update?: pulumi.Input<Cmd>
  /** Define a command to delete the resource. If unspecified, a delete operation is a no-op. */
  readonly delete?: pulumi.Input<Cmd>
}

/** Execute a command and save it as a resource */
export class Command extends pulumi.CustomResource {
  public readonly stdout: pulumi.Output<string>
  public readonly stderr: pulumi.Output<string>

  constructor(
    name: string,
    args: CommandSet,
    opts?: pulumi.CustomResourceOptions
  ) {
    ;(args as any).stdout = undefined /* out */
    ;(args as any).stderr = undefined /* out */
    if (typeof args.update === 'undefined') {
      ;(args as any).update = args.create
    }
    super('command:exec:command', name, args, opts)
  }
}
