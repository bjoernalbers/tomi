//go:build ignore

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

func TestVersion(t *testing.T) {
	got, err := getTomiVersion(tomi)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if got == "" {
		t.Fatalf("version not set")
	}
	// Regular expression taken from here:
	// https://semver.org/#is-there-a-suggested-regular-expression-regex-to-check-a-semver-string
	re := regexp.MustCompile(`^(?P<major>0|[1-9]\d*)\.(?P<minor>0|[1-9]\d*)\.(?P<patch>0|[1-9]\d*)(?:-(?P<prerelease>(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+(?P<buildmetadata>[0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	m := re.FindStringSubmatch(got)
	if m == nil {
		t.Fatalf("no version captured in output: %q", got)
	}
}

func getTomiVersion(command string) (string, error) {
	cmd := exec.Command(command, "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(string(output))
	}
	re := regexp.MustCompile(`\(version (.+)\)`)
	m := re.FindStringSubmatch(string(output))
	if m == nil || len(m) < 2 {
		return "", fmt.Errorf("no version captured: %v", m)
	}
	return m[1], nil
}
