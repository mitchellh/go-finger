# go-finger

go-finger is a [finger](https://en.wikipedia.org/wiki/Finger_protocol)
library written in Go. This contains both a client and server implementation.

The [finger protocol](https://tools.ietf.org/html/rfc1288) is an extremely
simple TCP protocol. It can be implemented without a library cleanly in only
a few dozen lines of code but a library helps ensure correctness and handles
RFC-compliant request parsing automatically.
