package pkg

import "testing"

func TestPreinstallScript(t *testing.T) {
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
sudo -u "${username}" --set-home "${tomi}"
else
echo "ERROR: Could not determine active user."
exit 1
fi
`
	got := preinstallScript()
	if got != want {
		t.Fatalf("got:\n%s\n\nwant:\n%s", got, want)
	}
}
