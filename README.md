# distributed-hash

A distributed hash table written in Go


# Architecture

## Components

client/hc.go is the client that stores/retrieves the data across multiple
parts of the hash table.

server/hd.go is the daemon that runs as part of the hash table.

## Overview

The hc command hashes the key to pick one of the hd servers to connect to.

## Server pool and notes

The client does not monitor the servers.

The client does not write the same key to nor read from multiple servers.
If a server process goes down the client will not try another server.

The server only stores data in volatile memory.  Data is no longer
retrievable after a server process terminates.  (GitHub Issue 5.)

# Client

## Usage

To store the value 123 with key A run:

```
hc put A 123
```

To retrieve the value for key A run:

```
hc get A
```

Keys are case-sensitive.  The GET and PUT commands are case-insensitive.

## Configuration

hc reads the list of servers out of hc.conf.  The file contains JSON. eg:

```
  {
    "Servers": [
      "127.0.0.1:1750",
      "127.0.0.1:1759"
    ]
  }
```


The hc.conf must be in the current working directory. (GitHub Issue 4.)

# Server process

## Usage

Start the process on port 1234 by running:

```
hd -p 1234
```

The default port is 1742.

## Running multiple servers on localhost

There are scripts to create and terminate 10 hd processes on localhost.
They store PID files in a directory named run-dir .
- server/start-hd
- server/stop-hd
