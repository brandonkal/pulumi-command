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
   * Hint: an easy method to always run update is to set diff to `['false']` */
  readonly diff?: pulumi.Input<Cmd>
  /** Specify a command to create a resource. */
  readonly create: pulumi.Input<Cmd>
  /** Specify a command to create read the resource. */
  readonly read?: pulumi.Input<Cmd>
  /** Specify a command to update the resource. */
  readonly update: pulumi.Input<Cmd>
  /** Specify a command to delete the resource. If unspecified, a delete operation is a no-op. */
  readonly delete?: pulumi.Input<Cmd>
}

/**
 * Provides a Command Resource
 */
export class Command extends pulumi.CustomResource {
  public readonly stdout: pulumi.Output<string>
  public readonly stderr: pulumi.Output<string>

  constructor(name: string, args: CommandSet, opts?: pulumi.ResourceOptions) {
    super('command:exec:command', name, args, opts)
  }
}
