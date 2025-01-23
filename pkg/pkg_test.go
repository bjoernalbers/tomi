package pkg

import "testing"

func TestBuild(t *testing.T) {
	if err := Build([]string{}, "1.2.3"); err == nil {
		t.Fatalf("want error, but got nil")
	}
}

func TestPreinstallScript(t *testing.T) {
	got := preinstallScript([]string{"-A", "-a", "tomedo.example.com"})
	want := `#!/bin/sh

set -eu

export PATH=/bin:/sbin:/usr/bin:/usr/sbin

tomi="$(dirname "$0")/tomi"

loggedInUser() {
local user=$(stat -f %Su /dev/console)

if [ -z "${user}" ]; then exit 1; fi
if [ "${user}" = 'root' ]; then exit 1; fi
echo "${user}"
}

if username=$(loggedInUser); then
echo "Running tomi as user '${username}'..."
sudo -u "${username}" --set-home "${tomi}" "-A" "-a" "tomedo.example.com"
else
echo "ERROR: Could not determine active user."
exit 1
fi
`
	if got != want {
		t.Fatalf("got:\n%s\n\nwant:\n%s", got, want)
	}
}
