import * as command from '@brandonkal/pulumi-command'

const cmd = new command.Command('test', {
  diff: {
    command: ['echo', 'hello'],
  },
  create: {
    command: ['echo', 'hello'],
  },
  update: {
    command: ['echo', 'hello'],
  },
})
