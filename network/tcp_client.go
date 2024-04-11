package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

type TCPClient struct {
	ServerAddress string
	conn          net.Conn
	isConnected   bool
}

// create a client instance
func NewTCPClient(serverAddress string) *TCPClient {
	return &TCPClient{ServerAddress: serverAddress}
}

func (c *TCPClient) IsConnected() bool {
	return c.isConnected
}

// create a connection
func (c *TCPClient) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.ServerAddress)
	if err != nil {
		return fmt.Errorf("failed to connect to server at %s: %v", c.ServerAddress, err)
	}
	c.isConnected = true
	return nil
}

// close the connection
func (c *TCPClient) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			return err
		}
		c.isConnected = false
		return nil
	}
	return fmt.Errorf("connection does not exist")
}

// sendMessage to server
func (c *TCPClient) SendMessage(msg JsonMessage) error {

	if c.conn == nil {
		if err := c.Connect(); err != nil {
			return err
		}
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("序列化JsonMessage时出错: %v", err)
	}
	msgBytes = append(msgBytes, '\n')
	_, err = c.conn.Write(msgBytes)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}
