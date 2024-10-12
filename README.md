# Golang sockets basics

The very basics of Golang networking using pure sockets in client-server manner.

The `server` part provides basic functionality to handle arbitrary number of
clients at the one time. Upon receiving byte messages (intended to be ascii text
only), it simply forwards that message to other connected clients (except the
sender).

The `client` part provides basic functionality to connect to the `server` and
send messages (indented to be ascii text) and recevie messages (ascii text as
well) from the `server`. Really primitive stuff.

NOTE: the successor of this repo is [BunnyHop](https://github.com/VigLinat/BunnyHop) repo. Check it out for some real
chat functionality.
