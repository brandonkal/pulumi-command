import * as pulumi from "@pulumi/pulumi"
import { Command } from "@brandonkal/pulumi-command"

const cmd = new Command('demo', {
  create: {
    command: ["bash", "-c", "echo 'mytest111' > /tmp/mytest.txt"],
  },
  diff: {
    command: ["bash", "-c", "test 'mytest111' != \"$(cat /tmp/mytest.txt)\""],
  },
  delete: {
    command: ["rm", "/tmp/mytest.txt"]
  },
  update: {
    command: ["bash", "-c", "echo 'mytest111' > /tmp/mytest.txt"]
  },
  read: {
    command: ["cat", "/tmp/mytest.txt"]
  }
})

export const out = cmd.stdout
