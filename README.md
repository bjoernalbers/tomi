# tomi - the missing tomedo-installer

Setting up [tomedo](https://tomedo.de) on a Mac is simple, but tedious.
tomi automates this process to relieve lazy admins (like me) of some work.

## Features

tomi performs the following steps for the current user:

1. Download tomedo.app from tomedo server and install into `$HOME/Applications`
2. Configure tomedo.app to connect to tomedo server
3. Add tomedo.app to Dock

Updates are handled by tomedo itself, so tomi does nothing at all if
`$HOME/Applications/tomedo.app` already exists.

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
$tomi -h
tomi - the missing tomedo-installer (version unset)

Usage: tomi [options]

Options:
  -A    install Arzeko as well
  -P string
        path of tomedo server (default "/tomedo_live/")
  -a string
        address of tomedo server (default "allgemeinmedizin.demo.tomedo.org")
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
