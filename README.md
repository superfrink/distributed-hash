distributed-hash
================

A distributed hash table written in Go


Architecutre
============

client/hc.go is the client that stores/retrieves the data across multiple
parts of the hash table.

server/hd.go is the daemon that runs as part of the hash table.
