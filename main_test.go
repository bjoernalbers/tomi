//go:build ignore

package main

import (
	"os"
	"os/exec"
	"path/filepath"
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

func TestBuildPackage(t *testing.T) {
	cmd := exec.Command(tomi, "-b")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s", string(output))
	}
	path := "tomedo-installer.pkg"
	defer os.Remove(path)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("%v", err)
	}
	tmp, err := os.MkdirTemp("", "tomi-test-*")
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer os.RemoveAll(tmp)
	expandedPackage := filepath.Join(tmp, "tomedo-installer.pkg")
	if output, err := exec.Command("/usr/sbin/pkgutil", "--expand", path, expandedPackage).CombinedOutput(); err != nil {
		t.Fatalf("%s", string(output))
	}
	if _, err := os.Stat(filepath.Join(expandedPackage, "Scripts", "preinstall")); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := os.Stat(filepath.Join(expandedPackage, "Scripts", "tomi")); err != nil {
		t.Fatalf("%v", err)
	}
}
