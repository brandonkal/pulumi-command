import { Command } from '@brandonkal/pulumi-command'

const cmd = new Command('test', {
  diff: {
    command: ['false'],
  },
  create: {
    command: ['ls'],
    environment: {
      VAR: 'hello world!',
    },
  },
  update: {
    command: ['echo', 'hello'],
  },
})
