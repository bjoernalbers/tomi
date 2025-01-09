# tomi - the missing tomedo-installer

Setting up [tomedo](https://tomedo.de) on a Mac is simple, but tedious.
tomi automates this process to relieve lazy admins (like me) of some work.

## Features

tomi performs the following steps for the current user executing the command:

1. Download the latest tomedo.app from your tomedo server
2. Install it into `$HOME/Applications`
3. Configure the app to connect to your tomedo server
4. Add tomedo.app to the Dock

However, tomi does nothing at all if `$HOME/Applications/tomedo.app` already exists.
Updates are handled by tomedo itself.

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

To install Arzeko as well, run:

```
tomi -A http://192.128.0.42:8080/tomedo_live/
```

If tomedo clients can reach your tomedo server by its internal hostname,
using that instead of an IP address works too:

```
tomi http://tomedo:8080/tomedo_live/
tomi http://tomedo.example.com:8080/tomedo_live/
```
