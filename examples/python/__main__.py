import pulumi_command as command
from pulumi import Config, export

config = Config()
dir = config.get("dir") or "/tmp"
content = config.get("content") or "mytest111"

exec = command.Command("demo",
    create=command.CmdArgs(command=["bash", "-c", f"echo '{content}' > {dir}/mytest.txt"]),
    diff=command.CmdArgs(command=["bash", "-c", f"test '{content}' != \"$(cat {dir}/mytest.txt)\""]),
    delete=command.CmdArgs(command=["rm", f"{dir}/mytest.txt"]),
    update=command.CmdArgs(command=["bash", "-c", f"echo '{content}' > {dir}/mytest.txt"]),
    read=command.CmdArgs(command=["cat", f"{dir}/mytest.txt"]),
)

export("stdout", exec.stdout)
export("stderr", exec.stderr)
