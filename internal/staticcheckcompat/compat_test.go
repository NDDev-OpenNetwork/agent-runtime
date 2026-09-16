// SPDX-License-Identifier: AGPL-3.0-only

package staticcheckcompat

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type runnerFunc func(context.Context, invocation) error

func (fn runnerFunc) run(ctx context.Context, call invocation) error { return fn(ctx, call) }

func TestVerifyPinsMatchesRepositoryContracts(t *testing.T) {
	t.Parallel()
	root := filepath.Join("..", "..")
	if err := VerifyPins(root); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyPinsFailsClosedOnMissingOrDriftedPin(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	writeContract(t, root, "AGENTS.md", "go run honnef.co/go/tools/cmd/staticcheck@v0.8.1 ./...\n")
	writeContract(t, root, "CONTRIBUTING.md", "go run honnef.co/go/tools/cmd/staticcheck@v0.8.1 ./...\n")
	if err := os.MkdirAll(filepath.Join(root, ".github", "workflows"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := VerifyPins(root); err == nil {
		t.Fatal("missing workflow pin accepted")
	}
	writeContract(t, root, ".github/workflows/ci.yml", "go run honnef.co/go/tools/cmd/staticcheck@v0.7.0 ./...\n")
	if err := VerifyPins(root); err == nil {
		t.Fatal("drifted workflow pin accepted")
	}
	writeContract(t, root, ".github/workflows/ci.yml", "go run honnef.co/go/tools/cmd/staticcheck@v0.8.1 ./...\n")
	if err := VerifyPins(root); err != nil {
		t.Fatal(err)
	}
}

func TestVerifyProbesPinnedStaticcheckWithLocalToolchain(t *testing.T) {
	t.Parallel()
	root := contractRoot(t)
	var calls []invocation
	runner := runnerFunc(func(_ context.Context, call invocation) error {
		calls = append(calls, call)
		return nil
	})
	var stdout, stderr bytes.Buffer
	if err := verify(context.Background(), root, &stdout, &stderr, runner); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 2 {
		t.Fatalf("calls=%d", len(calls))
	}
	if calls[0].path != "go" || strings.Join(calls[0].args, " ") != "run honnef.co/go/tools/cmd/staticcheck@v0.8.1 -debug.version" {
		t.Fatalf("version call=%#v", calls[0])
	}
	if calls[1].path != "go" || strings.Join(calls[1].args, " ") != "run honnef.co/go/tools/cmd/staticcheck@v0.8.1 ." {
		t.Fatalf("probe call=%#v", calls[1])
	}
	for index, call := range calls {
		if call.dir != root || call.stdout == nil || call.stderr == nil || count(call.env, "GOTOOLCHAIN=local") != 1 {
			t.Fatalf("call %d lost directory, output, or toolchain boundary", index)
		}
	}
}

func TestVerifyFailsClosedOnExportDataIncompatibility(t *testing.T) {
	t.Parallel()
	root := contractRoot(t)
	err := verify(context.Background(), root, &bytes.Buffer{}, &bytes.Buffer{}, runnerFunc(func(_ context.Context, call invocation) error {
		if strings.Contains(strings.Join(call.args, " "), "-debug.version") {
			return nil
		}
		_, _ = io.WriteString(call.stderr, "cannot decode \"internal/byteorder\", export data version 4 is greater than maximum supported version 2\n")
		return errors.New("exit status 1")
	}))
	if err == nil || !strings.Contains(err.Error(), "cannot decode this Go toolchain's export data") {
		t.Fatalf("incompatibility not reported: %v", err)
	}
}

func writeContract(t *testing.T, root, relative, body string) {
	t.Helper()
	path := filepath.Join(root, relative)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

func contractRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeContract(t, root, "AGENTS.md", "go run honnef.co/go/tools/cmd/staticcheck@v0.8.1 ./...\n")
	writeContract(t, root, "CONTRIBUTING.md", "go run honnef.co/go/tools/cmd/staticcheck@v0.8.1 ./...\n")
	writeContract(t, root, ".github/workflows/ci.yml", "go run honnef.co/go/tools/cmd/staticcheck@v0.8.1 ./...\n")
	return root
}

func count(values []string, want string) int {
	found := 0
	for _, value := range values {
		if value == want {
			found++
		}
	}
	return found
}


