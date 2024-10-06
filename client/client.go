package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var remoteAddr = "localhost"
var remotePort = "50160"

func main() {
	parseArgs()
	remote := makeSocketAddress(remoteAddr, remotePort)
	logConnecting(remote)
	conn, err := net.Dial("tcp4", remote)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		logClosingConnection(conn)
		conn.Close()
	}()

	go mustCopy(os.Stdout, conn)
	mustCopy(conn, os.Stdin)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func logConnecting(remote string) {
	fmt.Fprintf(os.Stderr, "LOG [%s] Connecting to remote: %s\n", time.Now().Format(time.TimeOnly), remote)
}

func logClosingConnection(c net.Conn) {
	fmt.Fprintf(os.Stderr, "LOG [%s] Closing connection %s\n", time.Now().Format(time.TimeOnly), c.RemoteAddr().String())
}

// TODO: use built-in flag parsing
func parseArgs() {
	if len(os.Args) > 1 {
		remotePort = os.Args[1]
	}
	if len(os.Args) > 2 {
		remoteAddr = os.Args[2]
	}
}

func makeSocketAddress(address, port string) string {
	return fmt.Sprintf("%s:%s", address, port)
}
