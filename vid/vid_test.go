package vid_test

import (
	"fmt"
	"testing"
	. "TMAABE/vid"
	. "TMAABE/erasurecode"
)

func TestInit(t *testing.T) {
	N := 4 //"number of servers in the cluster"
	F := 1 //"number of faulty servers to tolerate"
	param := &ProtocolParams{
		N:  N,
		F:  F,
		ID: 0,
	}
	codec := NewReedSolomonCode(N-2*F, N)
	v := NewVID(0,0, *param, codec)
	//msgs, _ := v.Init()
	fmt.Printf("vid=%v\n",*v)
}
