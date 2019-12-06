import { Command } from '@brandonkal/pulumi-command'

const cmd = new Command('demo', {
  diff: {
    command: ['false'],
  },
  create: {
    command: ['ls'],
  },
  update: {
    command: ['bash', '-c', 'echo $VAR'],
    environment: {
      VAR: 'Hello Pulumi!',
    },
  },
})
