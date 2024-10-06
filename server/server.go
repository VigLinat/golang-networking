package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func handleConn(c net.Conn) {

	logConnReceived(c)
	io.Copy(c, c)
	fmt.Fprintf(os.Stdout, "MESSAGE [%s] Received byte sequence:", time.Now().Format(time.TimeOnly))
	io.Copy(os.Stdout, c) // TEMPORARY
	fmt.Fprintln(os.Stdout)
	c.Close()
	logConnClosing(c)
}

var hostPort = "50160"
var hostAddr = "localhost"

func main() {
	parseArgs()
	host := makeSocketAddress(hostAddr, hostPort)
	listener, err := net.Listen("tcp4", host)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
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
