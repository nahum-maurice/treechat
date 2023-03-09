# Treechat

A TCP-based, console-based chat system.

## Requirements

Go 1.18 or higher
> This project is under active development. Things may change quickly.

## Installation

```
$ make build
```

## Container

### Build the container

```
$ cd treechat
$ docker build -t treechat .
```

### Run it

```
$ docker run -it -p 3000:3000 --name gotreechat treechat
```

## Use
You can interact with the server using tools such as telnet.

```
$ telnet localhost 3000
```

