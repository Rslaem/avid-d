package network_test

import (
	"TMAABE/network"
	"testing"
)

func TestSR(T *testing.T) {
	network.TestTCPNodeCommunication(T)
}

func TestConn(T *testing.T) {
	network.TestTCPNode_Concurrency(T)
}
