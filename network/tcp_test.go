package network_test

import (
	"TMAABE/network"
	"testing"
)

func TestSR(t *testing.T) {
	network.TestTCPNodeCommunication(t)
}

func TestConn(t *testing.T) {
	network.TestTCPNode_Concurrency(t)
}

func TestRpc(t *testing.T) {
	network.TestRPCNodeCommunication(t)
}
