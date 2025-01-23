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
	pkg := "tomedo-installer.pkg"
	defer os.Remove(pkg)
	dir, err := extractPkg(pkg)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer os.RemoveAll(dir)
	if _, err := os.Stat(filepath.Join(dir, "Scripts", "preinstall")); err != nil {
		t.Fatalf("%v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "Scripts", "tomi")); err != nil {
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

func TestPackageVersion(t *testing.T) {
	if output, err := exec.Command(tomi, "-b").CombinedOutput(); err != nil {
		t.Fatalf(string(output))
	}
	pkg := "tomedo-installer.pkg"
	defer os.Remove(pkg)
	want, err := getTomiVersion(tomi)
	if err != nil {
		t.Fatalf("%v", err)
	}
	got, err := getPackageVersion(pkg)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if got != want {
		t.Fatalf("got: %q, want: %q", got, want)
	}
}

func TestPackageIdentifier(t *testing.T) {
	if output, err := exec.Command(tomi, "-b").CombinedOutput(); err != nil {
		t.Fatalf(string(output))
	}
	pkg := "tomedo-installer.pkg"
	defer os.Remove(pkg)
	want := "de.bjoernalbers.tomedo-installer"
	got, err := getPackageIdentifier(pkg)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if got != want {
		t.Fatalf("got: %q, want: %q", got, want)
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

func getPackageVersion(path string) (string, error) {
	dir, err := extractPkg(path)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	output, err := os.ReadFile(filepath.Join(dir, "PackageInfo"))
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`<pkg-info.* version="([^"]+)" `)
	m := re.FindStringSubmatch(string(output))
	if m == nil || len(m) < 2 {
		return "", fmt.Errorf("no version captured: %v", m)
	}
	return m[1], nil
}

func getPackageIdentifier(path string) (string, error) {
	dir, err := extractPkg(path)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	output, err := os.ReadFile(filepath.Join(dir, "PackageInfo"))
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`<pkg-info.* identifier="([^"]+)" `)
	m := re.FindStringSubmatch(string(output))
	if m == nil || len(m) < 2 {
		return "", fmt.Errorf("no identifier captured: %v", m)
	}
	return m[1], nil
}

// extractPkg extracts content of .pkg file into temp. dir and returns its path.
func extractPkg(path string) (string, error) {
	pkg := filepath.Base(path)
	tmp, err := os.MkdirTemp("", pkg)
	if err != nil {
		return "", err
	}
	dir := filepath.Join(tmp, pkg)
	output, err := exec.Command("/usr/sbin/pkgutil", "--expand", path, dir).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf(string(output))
	}
	return dir, nil
}

func TestStripBuildFlag(t *testing.T) {
	tests := []struct {
		in   []string
		want []string
	}{
		{
			[]string{},
			[]string{},
		},
		{
			[]string{"-A", "-a", "tomedo.example.com"},
			[]string{"-A", "-a", "tomedo.example.com"},
		},
		{
			[]string{"-A", "-a", "tomedo.example.com", "-b"},
			[]string{"-A", "-a", "tomedo.example.com"},
		},
	}
	for _, tt := range tests {
		if got, want := stripBuildFlag(tt.in), tt.want; !equal(got, want) {
			t.Errorf("got: %v want: %v", got, want)
		}
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
