//go:build ignore

package main

import (
	"os/exec"
	"strings"
	"testing"
)

const tomi = "./tomi"

func TestHelp(t *testing.T) {
	cmd := exec.Command(tomi, "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to show help: %v", err)
	}
	if want := "Usage:"; !strings.Contains(string(output), want) {
		t.Fatalf("Help does not include %q\n%s", want, output)
	}
}
