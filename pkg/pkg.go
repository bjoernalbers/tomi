package pkg

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Build(version string) error {
	payloadDir, err := os.MkdirTemp("", "payload-*")
	if err != nil {
		return fmt.Errorf("create payload dir: %v", err)
	}
	defer os.RemoveAll(payloadDir)
	scriptsDir, err := os.MkdirTemp("", "scripts-*")
	if err != nil {
		return fmt.Errorf("create scripts dir: %v", err)
	}
	defer os.RemoveAll(scriptsDir)
	preinstall, err := os.OpenFile(filepath.Join(scriptsDir, "preinstall"), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0777)
	if err != nil {
		return fmt.Errorf("create preinstall script: %v", err)
	}
	if _, err := io.Copy(preinstall, strings.NewReader(preinstallScript())); err != nil {
		return fmt.Errorf("copy preinstall script: %v", err)
	}
	preinstall.Close()
	tomi, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get tomi executable: %v", err)
	}
	if output, err := exec.Command("/bin/cp", "-p", tomi, filepath.Join(scriptsDir, "tomi")).CombinedOutput(); err != nil {
		return fmt.Errorf("copy tomi: %v", string(output))
	}
	cmd := exec.Command("/usr/bin/pkgbuild",
		"--identifier", "de.bjoernalbers.tomedo-installer",
		"--version", version,
		"--scripts", scriptsDir,
		"--root", payloadDir,
		"tomedo-installer.pkg")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", string(output))
	}
	return nil
}

func preinstallScript() string {
	return `#!/bin/sh

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
}
