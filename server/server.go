package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/google/uuid"
)

var hostPort = "50160"
var hostAddr = "localhost"
var peerChan chan string = make(chan string)

var room chan *Message = make(chan *Message)// One common room for now

type Client struct {
    Conn net.Conn
    // Name string // TODO: implement later
    uuid string // SHA-1
}

// Read reads data from Client connection to room
func (client *Client) Read(room chan <- *Message) {
    var reader *bufio.Reader = bufio.NewReader(client.Conn)
    defer client.Conn.Close()
    for {
        if content, err := reader.ReadBytes('\n'); err == nil {
            fmt.Fprintf(os.Stdout, "[%s]: %s\n", client.Conn.LocalAddr().String(), content)
            message := &Message{client, content} 
            room <- message
        } else {
            if err == io.EOF {
                fmt.Fprintf(os.Stdout, "Client [%s] disconnected\n", client.Conn.LocalAddr().String())
            } else {
                fmt.Fprintf(os.Stdout, "Client [%s] Read error: %s\n", client.Conn.LocalAddr().String(), err.Error())
            }
            break;
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
