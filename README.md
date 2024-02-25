# docker edit <volume> <file-path>

`sudoedit` but for files in Docker volumes.

## Installation

```shell
$ echo "FROM vlnd/docker-edit" | docker buildx build - -o type=local,dest=$HOME/.docker/cli-plugins --pull
```

## Usage

```shell
$ docker edit <volume> <file-path-in-volume>
```

## Building

```shell
$ make         # Build bin/docker-edit binary
$ make install # Install binary to ~/.docker/cli-plugins
```

## Demo

![demo](https://raw.githubusercontent.com/vvoland/docker-plugin-edit/master/.demo.gif)
