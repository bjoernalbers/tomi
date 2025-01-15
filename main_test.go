//go:build ignore

package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to show help: %v", err)
	}
	if want := "Usage:"; !strings.Contains(string(output), want) {
		t.Fatalf("Help does not include %q\n%s", want, output)
	}
}
