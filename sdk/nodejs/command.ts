import * as pulumi from '@pulumi/pulumi'

interface Cmd {
  command: string[]
  environment?: Record<string, string>
}

export interface CommandArgs {
  diff: Cmd
  create: Cmd
  read?: Cmd
  update: Cmd
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
