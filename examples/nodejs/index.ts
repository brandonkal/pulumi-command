import * as pulumi from "@pulumi/pulumi"
import { Command } from "@brandonkal/pulumi-command"

let config = new pulumi.Config();
let dir = config.get("dir") || "/tmp";
let content = config.get("content") || "mytest111";

const cmd = new Command('demo', {
  create: {
    command: ["bash", "-c", `echo '${content}' > ${dir}/mytest.txt`],
  },
  diff: {
    command: ["bash", "-c", `test '${content}' != \"$(cat ${dir}/mytest.txt)\"`],
  },
  delete: {
    command: ["rm", `${dir}/mytest.txt`]
  },
  update: {
    command: ["bash", "-c", `echo '${content}' > ${dir}/mytest.txt`]
  },
  read: {
    command: ["cat", `${dir}/mytest.txt`]
  }
})

export const out = cmd.stdout
