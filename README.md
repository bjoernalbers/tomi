# tomi - the missing tomedo-installer

Setting up [tomedo](https://tomedo.de) on a Mac is simple, but tedious:

1. Create new directory (since tomedo cannot be installed to `/Applications`)
2. Download (a probably outdated) tomedo.app from somewhere into this folder
3. Add tomedo to Dock
4. Launch tomedo.app for the first time and enter address of your tomedo-server
6. Wait a few minutes for tomedo's auto-update to complete

Just let tomi take care of that!
tomi is a command line program that builds custom installation packages to set up
tomedo (and Arzeko) automatically.

(Updates are still handled by tomedo itself.)

## Installation

Open Terminal.app and download the
[latest release of tomi](https://github.com/bjoernalbers/tomi/releases/latest)
using this command:

```
tomi=$(mktemp) && \
    curl -o $tomi -sL https://github.com/bjoernalbers/tomi/releases/latest/download/tomi && \
    chmod +x $tomi
```

tomi will be removed automatically on the next reboot.

## Usage

Use `$tomi -h` to show the help:

```
tomi - the missing tomedo-installer (version unset)

Usage: tomi [options]

Options:
  -A    install Arzeko as well
  -P string
        path of tomedo server (default "/tomedo_live/")
  -a string
        address of tomedo server (default "allgemeinmedizin.demo.tomedo.org")
  -b    build installation package
  -p string
        port of tomedo server (default "8080")
```

Running `$tomi` without any options will install tomedo from the official demo
server.
If you want to use *your* tomedo server, you must at least change the server
address!
For example this command will install both tomedo and
[Arzeko](https://tomedo.de/praxissoftware/praxisorganisation/)
from `192.168.0.42`:

```
$tomi -A -a 192.128.0.42
```

Add the build flag (-b) to create a custom `tomedo-installer.pkg`:

```
$tomi -A -a 192.128.0.42 -b
```
