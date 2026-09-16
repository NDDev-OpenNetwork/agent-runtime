// SPDX-License-Identifier: AGPL-3.0-only

// Package staticcheckcompat pins Staticcheck and refuses a toolchain it cannot analyze.
package staticcheckcompat

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const Version = "v0.8.1"

const Module = "honnef.co/go/tools/cmd/staticcheck"

var ContractFiles = []string{
	"AGENTS.md",
	"CONTRIBUTING.md",
	".github/workflows/ci.yml",
}

type invocation struct {
	path   string
	args   []string
	dir    string
	env    []string
	stdout io.Writer
	stderr io.Writer
}

type commandRunner interface {
	run(context.Context, invocation) error
}

type osCommandRunner struct{}

func (osCommandRunner) run(ctx context.Context, call invocation) error {
	command := exec.CommandContext(ctx, call.path, call.args...)
	command.Dir, command.Env = call.dir, call.env
	command.Stdout, command.Stderr = call.stdout, call.stderr
	return command.Run()
}

func VerifyPins(root string) error {
	expected := Module + "@" + Version
	for _, relative := range ContractFiles {
		path := filepath.Join(root, relative)
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", relative, err)
		}
		text := string(data)
		if !strings.Contains(text, expected) {
			return fmt.Errorf("%s does not pin %s", relative, expected)
		}
		for _, line := range strings.Split(text, "\n") {
			index := strings.Index(line, Module+"@")
			if index < 0 {
				continue
			}
			got := strings.Fields(line[index:])[0]
			if got != expected {
				return fmt.Errorf("%s pins %s, want %s", relative, got, expected)
			}
		}
	}
	return nil
}

func Verify(ctx context.Context, root string, stdout, stderr io.Writer) error {
	return verify(ctx, root, stdout, stderr, osCommandRunner{})
}

func verify(ctx context.Context, root string, stdout, stderr io.Writer, runner commandRunner) error {
	if err := VerifyPins(root); err != nil {
		return err
	}
	env := localToolchain(os.Environ())
	versionOut, versionErr := &bytes.Buffer{}, &bytes.Buffer{}
	if err := runner.run(ctx, invocation{
		path:   "go",
		args:   []string{"run", Module + "@" + Version, "-debug.version"},
		dir:    root,
		env:    env,
		stdout: io.MultiWriter(stdout, versionOut),
		stderr: io.MultiWriter(stderr, versionErr),
	}); err != nil {
		return fmt.Errorf("staticcheck %s -debug.version: %w", Version, err)
	}
	if err := compatibilityError(versionOut.String() + versionErr.String()); err != nil {
		return err
	}
	analyzeOut, analyzeErr := &bytes.Buffer{}, &bytes.Buffer{}
	if err := runner.run(ctx, invocation{
		path:   "go",
		args:   []string{"run", Module + "@" + Version, "."},
		dir:    root,
		env:    env,
		stdout: io.MultiWriter(stdout, analyzeOut),
		stderr: io.MultiWriter(stderr, analyzeErr),
	}); err != nil {
		combined := analyzeOut.String() + analyzeErr.String()
		if compat := compatibilityError(combined); compat != nil {
			return compat
		}
		return fmt.Errorf("staticcheck %s compatibility probe: %w", Version, err)
	}
	return compatibilityError(analyzeOut.String() + analyzeErr.String())
}

func localToolchain(env []string) []string {
	result := make([]string, 0, len(env)+1)
	for _, entry := range env {
		if strings.HasPrefix(entry, "GOTOOLCHAIN=") {
			continue
		}
		result = append(result, entry)
	}
	return append(result, "GOTOOLCHAIN=local")
}

func compatibilityError(output string) error {
	if strings.Contains(output, "export data version") {
		return fmt.Errorf("staticcheck %s cannot decode this Go toolchain's export data; the pinned analyzer is older than the local compiler", Version)
	}
	return nil
}
