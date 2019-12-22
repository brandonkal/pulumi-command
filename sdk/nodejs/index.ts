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
import { createHash } from 'crypto'

function hash(input: any) {
  const str = JSON.stringify(input) || 'undefined'
  return createHash('sha256')
    .update(str)
    .digest('hex')
}

export interface Cmd {
  /** Specifiy the command to run as an array of arguments */
  command: pulumi.Input<string[]>
  /** Pass the stdin to a command */
  stdin?: pulumi.Input<string>
  /** Set environment variables for the running command */
  environment?: pulumi.Input<Record<string, string>>
}

export interface CommandSet {
  /** Specify a command to run to diff the resource.
   *
   * Exit 0 to run update.
   * Exit with a non-zero value or omit to disable update.
   * Hint: an easy method to always run update is to set diff to `['true']` */
  diff?: pulumi.Input<Cmd> | string[]
  compare?: pulumi.Input<any>
  /** Define a command to create a resource. */
  create: pulumi.Input<Cmd> | string[]
  /** Define a command to create read the resource. */
  read?: pulumi.Input<Cmd> | string[]
  /** If unspecified, create definition will be used. Define to provide an alternate update command. */
  update?: pulumi.Input<Cmd> | string[]
  /** Define a command to delete the resource. If unspecified, a delete operation is a no-op. */
  delete?: pulumi.Input<Cmd> | string[]
}

// fix unifies schema passed to the provider allowing for convenience array support
function fix(item: any) {
  return item ? (Array.isArray(item) ? { command: item } : item) : undefined
}

/** Execute a Command and save it as a resource.
 *
 * Each command can be specified as an object or a convenience array.
 * If only `create` is specified, `update` will use the create definition.
 * The `compare` property will be JSON serialized with a hash saved in the state. This is useful to ensure update is run if dependendent resources change.
 *
 * An update will occur in these cases:
 * 1. The `compare` hash or the `update` arguments change.
 * 2. The specified `diff` command exits with an error.
 */
export class Command extends pulumi.CustomResource {
  public readonly stdout: pulumi.Output<string>
  public readonly stderr: pulumi.Output<string>

  constructor(
    name: string,
    args: CommandSet,
    opts?: pulumi.CustomResourceOptions
  ) {
    const inputs: CommandSet = {
      create: fix(args.create),
      read: fix(args.read),
      update: fix(args.update),
      delete: fix(args.delete),
      diff: fix(args.diff),
    }
    ;(inputs as any).stdout = undefined /* out */
    ;(inputs as any).stderr = undefined /* out */
    if (typeof args.update === 'undefined') {
      inputs.update = args.create
    }
    if (
      typeof args.compare === 'undefined' &&
      typeof args.compare !== 'string'
    ) {
      try {
        inputs.compare = hash(args.compare)
      } catch (e) {
        throw new pulumi.RunError(`Could not serialize compare prop ${e}`)
      }
    }
    if (inputs.create === undefined) {
      throw new Error("Missing required property 'create'")
    }
    super('command:v1:exec', name, inputs, opts)
  }
}
