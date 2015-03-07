# parody

A simple http proxy that is designed to create a fire and forget copy of any
POST or PUT requests to a second upstream client.

This is raw and my use case for it is in a firewalled environment, so no
special security measures were taken.

Additionaly, there seems to be about ~41ms or latency added when load testing
on a localhost port. Timers indicate that this is mainly due to the
`res.Write(w)` call. Any suggestions on how to futher optimize this would be
fantasic.

### Install

```
go get -u github.com/jmervine/parody/v1
```

#### Download

Currently I only have the `linux/x86_64` binaries ready.

```
$ curl -sS -O http://static.mervine.net/go/linux/x86_64/parody && chmod 755 parody
```


### Usage

```
$ parody help
NAME:
   parody - simple http proxy for copying posts and puts to two locations

USAGE:
   parody [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHOR:
  Joshua Mervine - <joshua@mervine.net>

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --listen, -l "localhost:8888"        listener address
   --main, -m                           main upstream location [host:port]
   --copy, -c                           copy upstream location [host:port]
   --main-name "main"                   main upstream name in logger
   --copy-name "copy"                   copy upstream name in logger
   --help, -h                           show help
   --version, -v                        print the version

```

### From Source

```
git clone https://github.com/jmervine/parody.git && cd parody
make clean build install
```
