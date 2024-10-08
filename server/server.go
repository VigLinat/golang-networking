package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
    "github.com/google/uuid"
)

var hostPort = "50160"
var hostAddr = "localhost"
var peerChan chan string = make(chan string)

var room chan *Message // One common room for now

type Client struct {
    Conn net.Conn
    // Name string // TODO: implement later
    uuid string // SHA-1
}

// Read reads data from Client connection to room
func (client *Client) Read(room chan <- *Message) {
    buffer := make([]byte, 2048)
    for {
        if _, err := client.Conn.Read(buffer); err != nil {
            message := &Message{client, buffer} 
            room <- message
        }
    }
}

// Write writes data from room to Client connection
func (client *Client) Write(room <-chan *Message) {
    for message := range room {
        client.Conn.Write(message.Content)
    }
}

type Message struct {
    *Client
    Content []byte
}

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
        logConnReceived(conn)
        handleNewClient(conn)
	}
}

func handleNewClient(c net.Conn) {
    logConnReceived(c)
    clientUUID := uuid.NewString()
    newClient := &Client{c, clientUUID}
    go newClient.Read(room)
    go newClient.Write(room)
}

func logConnReceived(c net.Conn) {
	fmt.Fprintf(os.Stderr, "LOG [%s] Received connection: %s\n", time.Now().Format(time.TimeOnly), c.RemoteAddr().String())
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
