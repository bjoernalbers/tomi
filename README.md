# tomi - the tomedo-installer

Setting up [tomedo](https://tomedo.de) on a Mac is simple, but tedious.
tomi automates this process to relieve lazy admins (like me) of some work.

## Installation

Download the [latest tomi binary](https://github.com/bjoernalbers/tomi/releases/latest)
and make it executable.

## Usage

Just run `tomi` with your tomedo server as argument (hostname or IP addres), i.e.:

```
tomi tomedo.example.com
```

This will download the tomedo.app from your tomedo server and install it under
`$HOME/Programme`.
