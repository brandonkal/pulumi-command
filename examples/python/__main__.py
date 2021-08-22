import pulumi_command as command
from pulumi import export

exec = command.Command("demo",
    create=command.CmdArgs(command=["bash", "-c", "echo 'mytest111' > /tmp/mytest.txt"]),
    diff=command.CmdArgs(command=["bash", "-c", "test 'mytest111' != \"$(cat /tmp/mytest.txt)\""]),
    delete=command.CmdArgs(command=["rm", "/tmp/mytest.txt"]),
    update=command.CmdArgs(command=["bash", "-c", "echo 'mytest111' > /tmp/mytest.txt"]),
    read=command.CmdArgs(command=["cat", "/tmp/mytest.txt"]),
)

export("stdout", exec.stdout)
export("stderr", exec.stderr)