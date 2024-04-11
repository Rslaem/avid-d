package network

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

type TCPServer struct {
	Address string
}

func NewTCPServer(address string) *TCPServer {
	return &TCPServer{Address: address}
}

// Listening the messages
func (server *TCPServer) Listen(handleMessage func(JsonMessage)) {
	listener, err := net.Listen("tcp", server.Address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", server.Address, err)
	}
	defer listener.Close()
	log.Printf("Server listening on %s", server.Address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go func() {
			defer conn.Close()
			reader := bufio.NewReader(conn)

			message, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Failed to read from connection: %v", err)
				return
			}

			var msg JsonMessage
			if err := json.Unmarshal([]byte(message), &msg); err != nil {
				log.Printf("Failed to unmarshal JsonMessage: %v", err)
				return
			}

			handleMessage(msg)
		}()
	}
}
