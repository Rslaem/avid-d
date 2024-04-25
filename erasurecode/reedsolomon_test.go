<<<<<<< HEAD
package erasurecode_test

import (
	"fmt"
	"log"
	"testing"

	. "github.com/QinYuuuu/avid-d/erasurecode"
)

func TestReedSolomonCode(t *testing.T) {
	N := 4 //"number of servers in the cluster"
	F := 1 //"number of faulty servers to tolerate"
	rscode := NewReedSolomonCode(N-2*F, N)
	var codec ErasureCode = rscode
	data := []byte("a test message")

	var payload Payload = data

	eschunk, err := codec.Encode(payload)
	if err != nil {
		log.Fatal("encode wrong!", err)
	}
	rschunk := make([]ReedSolomonChunk, N)
	for i := 0; i < N; i++ {
		if tmp, ok := eschunk[i].(*ReedSolomonChunk); ok {
			rschunk[i] = *tmp
			//fmt.Println("Ok Value =", rschunk, "Ok =", ok)
		} else {
			fmt.Println("Failed Value =", rschunk, "Ok =", ok)
		}
	}
	for i := 0; i < N; i++ {
		fmt.Println(rschunk[i])
	}

	eschunk2 := make([]ErasureCodeChunk, N-F)
	for i := 0; i < N-F; i++ {
		eschunk2[i] = &rschunk[i]
	}

	var message Payload
	codec.Decode(eschunk2, &message)
	fmt.Println(string(message.([]byte)))
}
=======
package erasurecode_test

import (
	"fmt"
	"log"
	"testing"
	. "TMAABE/erasurecode"
)

func TestReedSolomonCode(t *testing.T) {
	N := 4 //"number of servers in the cluster"
	F := 1 //"number of faulty servers to tolerate"
	rscode := NewReedSolomonCode(N-2*F, N)
	var codec ErasureCode = rscode
	data := []byte("a test message")

	var payload Payload = data

	eschunk, err := codec.Encode(payload)
	if err != nil {
		log.Fatal("encode wrong!", err)
	}
	rschunk := make([]ReedSolomonChunk, N)
	for i := 0; i < N; i++ {
		if tmp, ok := eschunk[i].(*ReedSolomonChunk); ok {
			rschunk[i] = *tmp
			//fmt.Println("Ok Value =", rschunk, "Ok =", ok)
		} else {
			fmt.Println("Failed Value =", rschunk, "Ok =", ok)
		}
	}
	for i := 0; i < N; i++ {
		fmt.Println(rschunk[i])
	}

	eschunk2 := make([]ErasureCodeChunk, N-F)
	for i:=0;i<N-F;i++{
		eschunk2[i] = &rschunk[i]
	}

	var message Payload
	codec.Decode(eschunk2, &message)
	fmt.Println(string(message.([]byte)))
}
>>>>>>> e982a5d3560d233384b7cc8b8a3b52c93986a5ee
