# tomi - the missing tomedo-installer

Setting up [tomedo](https://tomedo.de) on a Mac is simple, but tedious.
tomi automates this process to relieve lazy admins (like me) of some work.

## Installation

To install the [latest release](https://github.com/bjoernalbers/tomi/releases/latest)
on a Mac, run this command:

```
curl -sLO https://github.com/bjoernalbers/tomi/releases/latest/download/tomi && \
    chmod +x tomi && \
    mv tomi /usr/local/bin
```

Running `rm /usr/local/bin/tomi` will uninstall it again.

## Usage

Just run `tomi` with the URL of *your* tomedo server, i.e.:

```
tomi http://192.128.0.42:8080/tomedo_live/
```

This will download the tomedo.app from your tomedo server, install it under
`$HOME/Programme` and add it to the Dock.

If tomedo clients can reach your tomedo server by its internal hostname,
using that instead of an IP address works too:

```
tomi http://tomedo:8080/tomedo_live/
tomi http://tomedo.example.com:8080/tomedo_live/
```
