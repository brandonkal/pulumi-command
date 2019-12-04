import * as pulumi from '@pulumi/pulumi'

export interface Cmd {
  /** Specifiy the command to run as an array of arguments */
  command: string[]
  /** Set environment variables for the running command */
  environment?: Record<string, string>
  /** Kill the Process after a specified number of milliseconds if specified */
  timeout?: number
}

export interface CommandArgs {
  /** Specify a command to run to diff the resource. If not specified, update will always run. */
  diff?: Cmd
  /** Specify a command to create a resource. */
  create: Cmd
  /** Specify a command to create read the resource. */
  read?: Cmd
  /** Specify a command to update the resource. */
  update: Cmd
  /** Specify a command to delete the resource. If unspecified, a delete operation is a no-op. */
  delete?: Cmd
}

/**
 * Provides a Command Resource
 */
export class Command extends pulumi.CustomResource {
  constructor(name: string, args: CommandArgs, opts?: pulumi.ResourceOptions) {
    super('command:exec:command', name, args, opts)
  }
}
