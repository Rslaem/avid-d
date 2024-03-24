package vid_test

import (
	. "github.com/QinYuuuu/avid-d/erasurecode"
	"github.com/QinYuuuu/avid-d/hasher"
	. "github.com/QinYuuuu/avid-d/vid"
	"fmt"
	"testing"
)

func sendData(ch chan Message, msgs Message) {
	ch <- msgs
}

func receiveData(ch <-chan Message) Message {
	msg := <-ch
	return msg
}

//func TestInit(t *testing.T) {
//	N := 4 //"number of servers in the cluster"
//	F := 1 //"number of faulty servers to tolerate"
//	param := &ProtocolParams{
//		N:  N,
//		F:  F,
//		ID: 0,
//	}
//	codec := NewReedSolomonCode(N-2*F, N)
//	v := NewVID(0, 0, *param, codec)
//	//msgs, _ := v.Init()
//	fmt.Printf("vid=%v\n", *v)
//}

func TestRecv(t *testing.T) {
	N := 4
	F := 1

	var paramsArray [4]*ProtocolParams
	for i := 0; i < len(paramsArray); i++ {
		paramsArray[i] = &ProtocolParams{N: N, F: F, ID: i}
	}
	codec := NewReedSolomonCode(N-2*F, N)

	var v [4]*VID

	str := "Hello world!"
	data := []byte(str)

	for i := 0; i < len(v); i++ {
		v[i] = NewVID(0, 0, *paramsArray[i], codec)
	}
	v[0].SetPayload(hasher.SHA256Hasher(data))
	var msg []Message
	msg, _ = v[0].Init()

	var channels [4]chan Message
	for i := range channels {
		channels[i] = make(chan Message)
	}

	// Recv--init & Disperse
	for len(msg) > 0 {
		go sendData(channels[msg[0].Dest()], msg[0])
		tmp1 := receiveData(channels[msg[0].Dest()])
		tmp2, _ := v[msg[0].Dest()].Recv(tmp1)

		msg = append(msg, tmp2...)
		msg = msg[1:]
	}

	// Recv--Retrieve
	msg = v[0].RequestPayload()
	for i := 0; i < N; i++ {
		v[i].InitRetrieve()
	}

	for len(msg) > 0 {
		go sendData(channels[msg[0].Dest()], msg[0])
		tmp1 := receiveData(channels[msg[0].Dest()])
		tmp2, _ := v[msg[0].Dest()].Recv(tmp1)

		msg = append(msg, tmp2...)
		msg = msg[1:]

		msg = append(msg, v[0].RequestPayload()...)

		if v[0].IfCanceled() {
			break
		}
	}
	//msg = append(msg, v[0].RequestPayload()...)
	//for len(msg) > 0 {
	//	go sendData(channels[msg[0].Dest()], msg[0])
	//	tmp1 := receiveData(channels[msg[0].Dest()])
	//	tmp2, _ := v[msg[0].Dest()].Recv(tmp1)
	//
	//	msg = append(msg, tmp2...)
	//	msg = msg[1:]
	//
	//	msg = append(msg, v[0].RequestPayload()...)
	//
	//	if v[0].IfCanceled() {
	//		break
	//	}
	//}
	fmt.Println(v[0].IfCanceled())
	fmt.Println("ttttt")

}
