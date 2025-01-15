//go:build ignore

package main

import (
	"os"
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

func TestHomeUnset(t *testing.T) {
	cmd := exec.Command(tomi)
	os.Setenv("HOME", "")
	if err := cmd.Run(); err == nil {
		t.Fatalf("tomi should fail if $HOME is unset")
	}
}
