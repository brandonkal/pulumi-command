import * as pulumi from '@pulumi/pulumi'
import { Command } from '@brandonkal/pulumi-command'

const cmd = new Command('demo', {
  diff: {
    command: ['true'],
  },
  create: {
    command: ['ls', '-lh'],
  },
  update: {
    command: ['bash', '-c', 'echo $VAR'],
    environment: {
      VAR: 'Hello Pulumi!',
    },
  },
})

export const out = cmd.stdout
