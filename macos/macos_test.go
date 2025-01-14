package macos

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAppExists(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/System/Library/CoreServices/Finder.app", true},
		{"/Applications/Missing.app", false},
	}
	for _, tt := range tests {
		a := App{Path: tt.path}
		if got := a.Exists(); got != tt.want {
			t.Errorf("%v.Exists():\ngot:\t%v\nwant:\t%v", a, got, tt.want)
		}
	}
}

func TestAppDir(t *testing.T) {
	a := App{Path: "/Applications/Foo.app"}
	want := "/Applications"
	if got := a.Dir(); got != want {
		t.Fatalf("%#v.Dir():\ngot:\t%q\nwant:\t%q", a, got, want)
	}
}

func TestCreateUserAppsDir(t *testing.T) {
	testHome, err := os.MkdirTemp("", "home")
	if err != nil {
		t.Fatalf("create test home: %v", err)
	}
	defer os.RemoveAll(testHome)
	dir, err := CreateUserAppsDir(testHome)
	if err != nil {
		t.Fatalf("create user apps dir: %v", err)
	}
	if _, err := os.Stat(filepath.Join(testHome, "Applications")); err != nil {
		t.Fatalf("user apps dir missing: %v", err)
	}
	if filepath.Base(dir) != "Applications" {
		t.Fatalf("wrong dirname: %q", dir)
	}
}
