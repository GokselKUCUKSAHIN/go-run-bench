package container

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
)

const (
	FixtureDir = "/fixture"
	WorkDir    = "/work"
)

type CLIContainer struct {
	ctx       context.Context
	container testcontainers.Container
}

func NewCLIContainer(ctx context.Context) *CLIContainer {
	return &CLIContainer{ctx: ctx}
}

func (c *CLIContainer) Run() error {
	repoRoot, err := findRepoRoot()
	if err != nil {
		return err
	}

	req := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    repoRoot,
			Dockerfile: "internal/integration-test/Dockerfile",
		},
		WorkingDir: WorkDir,
		Cmd:        []string{"sleep", "infinity"},
	}

	c.container, err = testcontainers.GenericContainer(c.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	return err
}

func (c *CLIContainer) Terminate() error {
	if c.container == nil {
		return nil
	}
	return c.container.Terminate(c.ctx)
}

func (c *CLIContainer) Exec(args ...string) (stdout string, stderr string, exitCode int, err error) {
	exitCode, reader, err := c.container.Exec(c.ctx, args, tcexec.Multiplexed())
	if err != nil {
		return "", "", exitCode, err
	}
	out, readErr := io.ReadAll(reader)
	if readErr != nil {
		return "", "", exitCode, readErr
	}
	return string(out), "", exitCode, nil
}

func (c *CLIContainer) RunCLI(args ...string) (stdout string, stderr string, exitCode int, err error) {
	cmd := append([]string{"go-run-bench"}, args...)
	return c.Exec(cmd...)
}

func (c *CLIContainer) ListWorkFiles() ([]string, error) {
	stdout, _, exitCode, err := c.Exec("sh", "-c", "ls -1 "+WorkDir)
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, fmt.Errorf("list work files failed: exit=%d out=%s", exitCode, stdout)
	}
	lines := strings.Split(strings.TrimSpace(stdout), "\n")
	var files []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

func (c *CLIContainer) ReadWorkFile(name string) ([]byte, error) {
	path := filepath.Join(WorkDir, name)
	stdout, _, exitCode, err := c.Exec("cat", path)
	if err != nil {
		return nil, err
	}
	if exitCode != 0 {
		return nil, fmt.Errorf("read %s failed: exit=%d out=%s", path, exitCode, stdout)
	}
	return []byte(stdout), nil
}

func (c *CLIContainer) FindResultFile(ext string) (string, error) {
	files, err := c.ListWorkFiles()
	if err != nil {
		return "", err
	}
	prefix := "benchmark_results_"
	suffix := "." + strings.TrimPrefix(ext, ".")
	for _, f := range files {
		if strings.HasPrefix(f, prefix) && strings.HasSuffix(f, suffix) {
			return f, nil
		}
	}
	return "", fmt.Errorf("no result file with extension %s in %v", suffix, files)
}

func (c *CLIContainer) CleanWorkDir() error {
	_, _, exitCode, err := c.Exec("sh", "-c", "rm -f "+WorkDir+"/benchmark_results_*")
	if err != nil {
		return err
	}
	if exitCode != 0 {
		return fmt.Errorf("clean work dir failed: exit=%d", exitCode)
	}
	return nil
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		modPath := filepath.Join(dir, "go.mod")
		data, readErr := os.ReadFile(modPath)
		if readErr == nil && bytes.Contains(data, []byte("module github.com/GokselKUCUKSAHIN/go-run-bench\n")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("repo root not found from %s", dir)
		}
		dir = parent
	}
}
