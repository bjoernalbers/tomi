# tomi - the tomedo-installer

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

Just run `tomi` with your tomedo server as argument (hostname or IP addres), i.e.:

```
tomi tomedo.example.com
```

This will download the tomedo.app from your tomedo server and install it under
`$HOME/Programme`.
