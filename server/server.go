package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var hostPort = "50160"
var hostAddr = "localhost"

var connectedClients map[string]net.Conn = make(map[string]net.Conn)

func main() {
	parseArgs()
	host := makeSocketAddress(hostAddr, hostPort)
	listener, err := net.Listen("tcp4", host)
	if err != nil {
		log.Fatal(err)
	}
    go listenForClients(listener)
}

func listenForClients(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
        handleNewClient(conn)
		go handleConn(conn)
	}
}

func listenToClient(client net.Conn) <-chan string {
    res := make(chan string)
    go func(conn net.Conn) {
        buff := make([]byte, 2048)
        for {
            if n, err := conn.Read(buff); err != nil {
                res <- string(buff[:n])
            }
        }
    }(client)
    return res
}

func handleNewClient(c net.Conn) {
    logConnReceived(c)
    connectedClients[c.LocalAddr().String()] = c
}

func handleConn(c net.Conn) {
	fmt.Fprintf(os.Stdout, "MESSAGE [%s] Received byte sequence:", time.Now().Format(time.TimeOnly))
	io.Copy(os.Stdout, c) // TEMPORARY
	fmt.Fprintln(os.Stdout)

	logConnClosing(c)
}

func logConnReceived(c net.Conn) {
	fmt.Fprintf(os.Stderr, "LOG [%s] Received connection: %s\n", time.Now().Format(time.TimeOnly), c.RemoteAddr().String())
}

func logConnClosing(c net.Conn) {
	fmt.Fprintf(os.Stderr, "LOG [%s] Closing connection %s\n", time.Now().Format(time.TimeOnly), c.RemoteAddr().String())
}

// TODO: use built-in flag parsing
func parseArgs() {
	if len(os.Args) > 1 {
		hostPort = os.Args[1]
	}
	if len(os.Args) > 2 {
		hostAddr = os.Args[2]
	}
}

func makeSocketAddress(address, port string) string {
	return fmt.Sprintf("%s:%s", address, port)
}
