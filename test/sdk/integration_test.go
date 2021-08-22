package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
)

func TestIntegration(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		stackPath   string
		preCommands [][]string
	}{
		{
			name:      "dotnet",
			stackPath: "../../examples/dotnet",
		},
		{
			name:      "go",
			stackPath: "../../examples/go",
		},
		{
			name:        "nodejs",
			stackPath:   "../../examples/nodejs",
			preCommands: [][]string{{"yarn", "install"}},
		},
		{
			name:      "python",
			stackPath: "../../examples/python",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			stackName := "test"
			workDir := filepath.Join(".", tc.stackPath)

			if len(tc.preCommands) > 0 {
				for _, args := range tc.preCommands {
					cmd := exec.Command(args[0], args[1:]...)
					cmd.Dir = workDir
					err := cmd.Run()
					if err != nil {
						t.Fatalf("failed pre command: %v", err)
					}
				}
			}

			stack, err := auto.UpsertStackLocalSource(ctx, stackName, workDir)
			if err != nil {
				t.Fatalf("failed stack init: %v", err)
			}

			// save logs and output in case of failure
			logs := &bytes.Buffer{}
			logsDestroy := optdestroy.ProgressStreams(logs)
			logsUp := optup.ProgressStreams(logs)

			// prepare test values
			dir := t.TempDir()
			stack.SetConfig(ctx, "dir", auto.ConfigValue{Value: dir})

			content := random(10)
			stack.SetConfig(ctx, "content", auto.ConfigValue{Value: content})

			// 1. Run stack
			_, err = stack.Up(ctx, logsUp)
			if err != nil {
				t.Errorf("failed stack creation: %v", err)
			}

			err = checkFile(true, dir, content)
			if err != nil {
				t.Errorf("stack creation wrong result: %v", err)
			}

			// 2. Update stack
			content = random(10)
			stack.SetConfig(ctx, "content", auto.ConfigValue{Value: content})

			_, err = stack.Up(ctx, logsUp)
			if err != nil {
				t.Errorf("failed stack update: %v", err)
			}

			err = checkFile(true, dir, content)
			if err != nil {
				t.Errorf("stack update wrong result: %v", err)
			}

			// 3. Destroy stack
			_, err = stack.Destroy(ctx, logsDestroy)
			if err != nil {
				t.Errorf("failed stack destroy: %v", err)
			}

			err = checkFile(false, dir, content)
			if err != nil {
				t.Errorf("stack destroy wrong result: %v", err)
			}

			t.Log(logs.String())
		})
	}
}

func checkFile(exist bool, dir, content string) error {
	path := filepath.Join(dir, "mytest.txt")
	_, err := os.Stat(path)

	if os.IsNotExist(err) && !exist {
		return nil
	}

	if err != nil && exist {
		return fmt.Errorf("file does not exist")
	}

	if err == nil && !exist {
		return fmt.Errorf("file not deleted")
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if v := strings.TrimSuffix(string(data), "\n"); v != content {
		return fmt.Errorf("wrong file content, expected: %s, got: %s", content, v)
	}

	return nil
}

func random(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
