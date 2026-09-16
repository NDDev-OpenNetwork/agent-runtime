// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/NDDev-OpenNetwork/agent-runtime/internal/staticcheckcompat"
)

func main() {
	if err := staticcheckcompat.Verify(context.Background(), ".", os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "staticcheck compatibility failed:", err)
		os.Exit(1)
	}
	fmt.Printf("staticcheck %s pin and toolchain probe succeeded\n", staticcheckcompat.Version)
}
