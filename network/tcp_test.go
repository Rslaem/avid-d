package network_test

import (
	"github.com/QinYuuuu/avid-d/network"
	"testing"
)

func TestSR(t *testing.T) {
	network.TestTCPNodeCommunication(t)
}

func TestConn(t *testing.T) {
	network.TestTCPNode_Concurrency(t)
}
