package main

import (
	"TMAABE/batchdkg"
	. "TMAABE/erasurecode"
	. "TMAABE/hasher"
	. "TMAABE/network"
	"TMAABE/tmaabe"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/Nik-U/pbc"
	//"net/http"
)

type DKGPayload1 struct {
	Z [][][]byte
	R [][]byte
}

type DKGPayload2 struct {
	Ga [][][]byte
}

type DKGPayload3 struct {
	Asum [][]byte
}

type ChunkMessage struct {
	Index int
	Chunk ReedSolomonChunk
	Check []byte
}

func generateIPList(n int, str string) {
	file, err := os.OpenFile(str+"/iplist", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}

	//ips := []string{"18.135.17.154", "18.133.238.135", "18.171.177.78",
	//	"18.171.177.233", "18.170.224.66", "3.10.198.201", "3.8.48.231", "18.133.244.123"}
	peers := []struct {
		Id   int
		Addr string
	}{}
	//for j, ip := range ips {
	for i := 0; i < n; i++ {
		peer := struct {
			Id   int
			Addr string
		}{i, fmt.Sprintf("%s:%d", "127.0.0.1", 5000+i)}
		peers = append(peers, peer)
	}
	//}

	fmt.Println(peers)
	data, _ := json.Marshal(peers)
	_, err = file.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func ReadGlobalParameters(metadataPath string) *tmaabe.GlobalParameters {
	paramData, err := ioutil.ReadFile(metadataPath + "/a_1.properties")
	if err != nil {
		log.Fatalf("node failed to read paramData %v\n", err)
		return nil
	}
	params := string(paramData)

	generatorData, err := ioutil.ReadFile(metadataPath + "/Generators")
	if err != nil {
		log.Fatalf("node failed to read generatorData %v\n", err)
		return nil
	}
	var generator map[string][]byte
	err = json.Unmarshal(generatorData, &generator)
	if err != nil {
		log.Fatalf("Error during Unmarshal() generator: %v\n", err)
		return nil
	}
	gp := tmaabe.NewGlobalParametersFromString(params, generator["g"], generator["h"])
	return gp
}

func test(nodeNum, secretNum, f, id int, path string) {
	s := NewServer(id, nodeNum, fmt.Sprintf("%s/iplist_%d", path, nodeNum/8))
	s.Register("/test", s.Incoming.ReceivePost)
	s.Init()
	var mutex sync.RWMutex

	disperse1 := make([]bool, nodeNum)
	disperse2 := make([]bool, nodeNum)
	disperse3 := make([]bool, nodeNum)
	ndisperse1s := make([]int, nodeNum)
	ndisperse2s := make([]int, nodeNum)
	ndisperse3s := make([]int, nodeNum)
	receiveChunk1 := make([][]ErasureCodeChunk, nodeNum)
	receiveChunk2 := make([][]ErasureCodeChunk, nodeNum)
	receiveChunk3 := make([][]ErasureCodeChunk, nodeNum)
	for i := 0; i < nodeNum; i++ {
		disperse1[i] = false
		disperse2[i] = false
		disperse3[i] = false
		ndisperse1s[i] = 0
		ndisperse2s[i] = 0
		ndisperse3s[i] = 0
		receiveChunk1[i] = make([]ErasureCodeChunk, nodeNum)
		receiveChunk2[i] = make([]ErasureCodeChunk, nodeNum)
		receiveChunk3[i] = make([]ErasureCodeChunk, nodeNum)
	}
	go func() {
		for m := range s.RecvChan() {
			//log.Printf("recv message type:%v\n", m.DataType)
			if m.DataType == "test" {
				log.Printf("recv message:%v\n", string(m.Content))
			} else if m.DataType == "disperse1" || m.DataType == "retrieve1" {
				var chunkmessage ChunkMessage
				err := json.Unmarshal(m.Content, &chunkmessage)
				if err != nil {
					log.Printf("json Unmarshal error: %v", err)
				}
				i := chunkmessage.Index
				j := chunkmessage.Chunk.Idx
				if !disperse1[i] {
					mutex.Lock()
					receiveChunk1[i][j] = &chunkmessage.Chunk
					ndisperse1s[i]++
					//log.Printf("[node %d] receive chunk from node %d in disperse1, chunk number %d", id, i, ndisperse1s[i])
					mutex.Unlock()
				}

			} else if m.DataType == "disperse2" || m.DataType == "retrieve2" {
				var chunkmessage ChunkMessage
				err := json.Unmarshal(m.Content, &chunkmessage)
				if err != nil {
					log.Printf("json Unmarshal error: %v", err)
				}
				i := chunkmessage.Index
				j := chunkmessage.Chunk.Idx
				if !disperse2[i] {
					mutex.Lock()
					receiveChunk2[i][j] = &chunkmessage.Chunk
					ndisperse2s[i]++
					//log.Printf("[node %d] receive chunk from node %d in disperse2, chunk number %d", id, i, ndisperse2s[i])
					mutex.Unlock()
				}
			} else if m.DataType == "disperse3" || m.DataType == "retrieve3" {
				var chunkmessage ChunkMessage
				err := json.Unmarshal(m.Content, &chunkmessage)
				if err != nil {
					log.Printf("json Unmarshal error: %v", err)
				}
				i := chunkmessage.Index
				j := chunkmessage.Chunk.Idx
				if !disperse3[i] {
					mutex.Lock()
					receiveChunk3[i][j] = &chunkmessage.Chunk
					ndisperse3s[i]++
					//log.Printf("[node %d] receive chunk from node %d in disperse2, chunk number %d", id, i, ndisperse3s[i])
					mutex.Unlock()
				}
			} else {
				continue
			}
		}
	}()

	n0, _ := new(big.Int).SetString("1363895147340162124487750544377566700025348452567", 10)
	n1, _ := new(big.Int).SetString("1257354545315887944833595666025792933231792977521", 10)
	n2, _ := new(big.Int).SetString("1296657106138026641358592699056954007605324218609", 10)
	n := new(big.Int)
	n.Mul(n0, n1)
	n.Mul(n, n2)

	log.Printf("[node %d] init tmaabe GlobalParameters", id)
	gp := ReadGlobalParameters(path)
	pairing := gp.GetPairing()
	g := gp.GetGenerateG()
	dkgParam := batchdkg.NewParam(pairing, g, n)
	start := time.Now()
	dkg := batchdkg.NewDKG(dkgParam, secretNum, nodeNum, f)
	dkg.ShareGen()

	var dkgPayload1 DKGPayload1
	dkgPayload1.Z = make([][][]byte, nodeNum)
	dkgPayload1.R = make([][]byte, secretNum)

	for i := 0; i < nodeNum; i++ {
		dkgPayload1.Z[i] = make([][]byte, secretNum+1)
		shares := dkg.GetShares()[i]
		log.Printf("[node %d] generat shares to node %d\n %v \n", s.ID, i, shares)
		for k := 0; k < secretNum+1; k++ {
			dkgPayload1.Z[i][k] = shares[k].Bytes()
		}
	}
	for k := 0; k < secretNum; k++ {
		dkgPayload1.R[k] = dkg.Getr()[k].Bytes()
	}

	codec := NewReedSolomonCode(nodeNum-2*f, nodeNum)
	/*
		go func() {
			for i := 0; i < nodeNum; i++ {
				if i == id {
					continue
				}
				respond, _ := s.Outgoing.SendPost(i, "test", "test", []byte(fmt.Sprintf("hello from node %d", s.ID)))
				log.Printf(string(respond))
				//fmt.Println(err)
			}
		}()*/

	var payload Payload
	jsonP, _ := json.Marshal(&dkgPayload1)
	payload = jsonP

	eschunk, err := codec.Encode(payload)
	if err != nil {
		log.Fatal("encode wrong!", err)
	}
	rschunk := make([]ReedSolomonChunk, nodeNum)
	for i := 0; i < nodeNum; i++ {
		if i == id {
			mutex.Lock()
			receiveChunk1[i] = eschunk
			ndisperse1s[i] = nodeNum
			disperse1[i] = true
			mutex.Unlock()
			continue
		}
		if tmp, ok := eschunk[i].(*ReedSolomonChunk); ok {
			go func(i int) {
				rschunk[i] = *tmp
				//fmt.Println("Ok Value =", rschunk, "Ok =", ok)
				data := ChunkMessage{
					Index: id,
					Chunk: rschunk[i],
					Check: SHA256Hasher(rschunk[i].GetData()),
				}
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Printf("json Marshal error: %v", err)
				}
				//log.Printf("json data: %v", jsonData)
				respond, _ := s.Outgoing.SendPost(i, "disperse1", "test", jsonData)
				log.Printf("[node %d] disperse1 send to node %d respond: %s", id, i, string(respond))
			}(i)
		} else {
			fmt.Println("Failed Value =", rschunk, "Ok =", ok)
		}
	}

	var wg sync.WaitGroup
	wg.Add(nodeNum - 1)
	for i := 0; i < nodeNum; i++ {
		if i == id {
			continue
		}
		go func(i int) {
			for {
				mutex.RLock()
				n := receiveChunk1[i][id] == nil
				mutex.RUnlock()
				if !n {
					break
				}
				time.Sleep(2 * time.Second)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	for index := 0; index < nodeNum; index++ {
		go func(index int) {
			//log.Printf("[node %d] index %d chunks %v\n", id, index, receiveChunk1[index])
			retrievechunk := receiveChunk1[index][id].(*ReedSolomonChunk)
			data := ChunkMessage{
				Index: index,
				Chunk: *retrievechunk,
				Check: SHA256Hasher(receiveChunk1[index][id].GetData()),
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("json Marshal error: %v", err)
			}
			//log.Printf("json data: %v", jsonData)
			for j := 0; j < nodeNum; j++ {
				if index == j || j == id {
					continue
				}
				//log.Printf("[node %d] retrieve1 send to node %d msg: %s", id, j, jsonData)
				s.Outgoing.SendPost(j, "retrieve1", "test", jsonData)
				//log.Printf("[node %d] retrieve1 send to node %d respond: %s", id, j, string(respond))
			}
		}(index)
	}

	wg.Add(nodeNum - 1)
	for i := 0; i < nodeNum; i++ {
		//log.Printf("[node %d] index %d chunks %v\n", id, i, receiveChunk1[i])
		if i == id {
			continue
		}
		go func(i int) {
			for {
				mutex.RLock()
				n := ndisperse1s[i]
				mutex.RUnlock()
				//log.Printf("[node %d] index %d chunks number %d", id, i, n)
				if n > nodeNum-2*f {
					disperse1[i] = true
					break
				}
				//log.Printf("[node %d] index %d chunks %v\n", id, i, receiveChunk1[i])
				time.Sleep(2 * time.Second)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	wg.Add(nodeNum)
	for i := 0; i < nodeNum; i++ {
		go func(i int) {
			var msg Payload
			tmp := make([]ErasureCodeChunk, 0)
			mutex.RLock()
			for _, chunk := range receiveChunk1[i] {
				if chunk != nil {
					tmp = append(tmp, chunk)
				}
				if len(tmp) == nodeNum-2*f {
					break
				}
			}
			mutex.RUnlock()
			codec.Decode(tmp, &msg)
			dp1 := DKGPayload1{}
			json.Unmarshal(msg.([]byte), &dp1)
			//log.Printf("[node %d] retrieve: %v \n", s.ID, dp1)
			shares := make([]*big.Int, secretNum+1)
			for k := 0; k < secretNum+1; k++ {
				shares[k] = new(big.Int).SetBytes(dp1.Z[s.ID][k])
			}
			mutex.Lock()
			dkg.SetShares(i, shares)
			mutex.Unlock()
			//log.Printf("[node %d] receive shares from node: %d:\n %v \n", s.ID, i, shares)
			r := make([]*big.Int, secretNum)
			for k := 0; k < secretNum; k++ {
				r[k] = new(big.Int).SetBytes(dp1.R[k])
			}
			mutex.Lock()
			dkg.ReceiveR(i, r)
			mutex.Unlock()
			wg.Done()
		}(i)
	}
	wg.Wait()
	dkg.FaultDetectPhase1()
	///////////////////////////////////////
	//                                   //
	// start disperse the second payload //
	//                                   //
	///////////////////////////////////////
	log.Printf("2")
	var dkgPayload2 DKGPayload2
	dkgPayload2.Ga = make([][][]byte, nodeNum)
	for i := 0; i < nodeNum; i++ {
		dkgPayload2.Ga[i] = make([][]byte, secretNum+1)
		comms := dkg.GetComms()[i]
		for k := 0; k < secretNum+1; k++ {
			dkgPayload2.Ga[i][k] = comms[k].Bytes()
		}
	}

	jsonP, _ = json.Marshal(&dkgPayload2)
	payload = jsonP

	eschunk, err = codec.Encode(payload)
	if err != nil {
		log.Fatal("encode wrong!", err)
	}
	rschunk = make([]ReedSolomonChunk, nodeNum)
	for i := 0; i < nodeNum; i++ {
		if i == id {
			mutex.Lock()
			receiveChunk2[i] = eschunk
			ndisperse2s[i] = nodeNum
			disperse2[i] = true
			mutex.Unlock()
			continue
		}
		if tmp, ok := eschunk[i].(*ReedSolomonChunk); ok {
			go func(i int) {
				rschunk[i] = *tmp
				//fmt.Println("Ok Value =", rschunk, "Ok =", ok)
				data := ChunkMessage{
					Index: id,
					Chunk: rschunk[i],
					Check: SHA256Hasher(rschunk[i].GetData()),
				}
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Printf("json Marshal error: %v", err)
				}
				//log.Printf("json data: %v", jsonData)
				s.Outgoing.SendPost(i, "disperse2", "test", jsonData)
				//log.Printf("[node %d] disperse2 send to node %d respond: %s", id, i, string(respond))
			}(i)
		} else {
			fmt.Println("Failed Value =", rschunk, "Ok =", ok)
		}
	}

	wg.Add(nodeNum - 1)
	for i := 0; i < nodeNum; i++ {
		if i == id {
			continue
		}
		go func(i int) {
			for {
				mutex.RLock()
				n := receiveChunk2[i][id] == nil
				mutex.RUnlock()
				if !n {
					break
				}
				time.Sleep(2 * time.Second)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	for index := 0; index < nodeNum; index++ {
		go func(index int) {
			//log.Printf("[node %d] index %d chunks %v\n", id, index, receiveChunk2[index])
			retrievechunk := receiveChunk2[index][id].(*ReedSolomonChunk)
			data := ChunkMessage{
				Index: index,
				Chunk: *retrievechunk,
				Check: SHA256Hasher(receiveChunk2[index][id].GetData()),
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("json Marshal error: %v", err)
			}
			//log.Printf("json data: %v", jsonData)
			for j := 0; j < nodeNum; j++ {
				if index == j || j == id {
					continue
				}
				//log.Printf("[node %d] retrieve2 send to node %d msg: %s", id, j, jsonData)
				s.Outgoing.SendPost(j, "retrieve2", "test", jsonData)
				//log.Printf("[node %d] retrieve2 send to node %d respond: %s", id, j, string(respond))
			}
		}(index)
	}

	wg.Add(nodeNum - 1)
	for i := 0; i < nodeNum; i++ {
		//log.Printf("[node %d] index %d chunks %v\n", id, i, receiveChunk2[i])
		if i == id {
			continue
		}
		go func(i int) {
			for {
				mutex.RLock()
				n := ndisperse2s[i]
				mutex.RUnlock()
				//log.Printf("[node %d] index %d chunks number %d", id, i, n)
				if n > nodeNum-2*f {
					disperse2[i] = true
					break
				}
				//log.Printf("[node %d] index %d chunks %v\n", id, i, receiveChunk2[i])
				time.Sleep(2 * time.Second)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	wg.Add(nodeNum)
	for i := 0; i < nodeNum; i++ {
		go func(i int) {
			var msg Payload
			tmp := make([]ErasureCodeChunk, 0)
			mutex.RLock()
			for _, chunk := range receiveChunk2[i] {
				if chunk != nil {
					tmp = append(tmp, chunk)
				}
				if len(tmp) == nodeNum-2*f {
					break
				}
			}
			mutex.RUnlock()
			codec.Decode(tmp, &msg)
			dp2 := DKGPayload2{}
			json.Unmarshal(msg.([]byte), &dp2)

			ga := make([][]*pbc.Element, nodeNum)
			for j := 0; j < nodeNum; j++ {
				ga[j] = make([]*pbc.Element, secretNum+1)
				for k := 0; k < secretNum+1; k++ {
					ga[j][k] = pairing.NewG1().SetBytes(dp2.Ga[j][k])
				}
			}
			mutex.Lock()
			//log.Printf("[node %d] receive commitment from node %d", s.ID, i)
			dkg.SetComms(i, ga)
			mutex.Unlock()
			wg.Done()
		}(i)
	}
	///////////////////////////////////////
	//                                   //
	// start disperse the third payload  //
	//                                   //
	///////////////////////////////////////
	log.Printf("3")
	var dkgPayload3 DKGPayload3
	dkgPayload3.Asum = make([][]byte, nodeNum)
	aijsum := dkg.Getaijsum()

	for i := 0; i < nodeNum; i++ {
		dkgPayload3.Asum[i] = aijsum[i].Bytes()
	}

	jsonP, _ = json.Marshal(&dkgPayload3)
	payload = jsonP

	eschunk, err = codec.Encode(payload)
	if err != nil {
		log.Fatal("encode wrong!", err)
	}
	rschunk = make([]ReedSolomonChunk, nodeNum)

	for i := 0; i < nodeNum; i++ {
		if i == id {
			mutex.Lock()
			receiveChunk3[i] = eschunk
			ndisperse3s[i] = nodeNum
			disperse3[i] = true
			mutex.Unlock()
			continue
		}
		if tmp, ok := eschunk[i].(*ReedSolomonChunk); ok {
			go func(i int) {
				rschunk[i] = *tmp
				//fmt.Println("Ok Value =", rschunk, "Ok =", ok)
				data := ChunkMessage{
					Index: id,
					Chunk: rschunk[i],
					Check: SHA256Hasher(rschunk[i].GetData()),
				}
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Printf("json Marshal error: %v", err)
				}
				//log.Printf("json data: %v", jsonData)
				s.Outgoing.SendPost(i, "disperse3", "test", jsonData)
				//log.Printf("[node %d] disperse3 send to node %d respond: %s", id, i, string(respond))
			}(i)
		} else {
			fmt.Println("Failed Value =", rschunk, "Ok =", ok)
		}
	}

	wg.Add(nodeNum - 1)
	for i := 0; i < nodeNum; i++ {
		if i == id {
			continue
		}
		go func(i int) {
			for {
				mutex.Lock()
				n := receiveChunk3[i][id] == nil
				mutex.Unlock()
				if !n {
					break
				}
				time.Sleep(2 * time.Second)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	for index := 0; index < nodeNum; index++ {
		go func(index int) {
			log.Printf("[node %d] index %d chunks %v\n", id, index, receiveChunk3[index])
			retrievechunk := receiveChunk3[index][id].(*ReedSolomonChunk)
			data := ChunkMessage{
				Index: index,
				Chunk: *retrievechunk,
				Check: SHA256Hasher(receiveChunk3[index][id].GetData()),
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				log.Printf("json Marshal error: %v", err)
			}
			//log.Printf("json data: %v", jsonData)
			for j := 0; j < nodeNum; j++ {
				if index == j || j == id {
					continue
				}
				log.Printf("[node %d] retrieve3 send to node %d msg: %s", id, j, jsonData)
				respond, _ := s.Outgoing.SendPost(j, "retrieve3", "test", jsonData)
				log.Printf("[node %d] retrieve3 send to node %d respond: %s", id, j, string(respond))
			}
		}(index)
	}

	wg.Add(nodeNum - 1)
	for i := 0; i < nodeNum; i++ {
		log.Printf("[node %d] index %d chunks %v\n", id, i, receiveChunk3[i])
		if i == id {
			continue
		}
		go func(i int) {
			for {
				mutex.RLock()
				n := ndisperse3s[i]
				mutex.RUnlock()
				log.Printf("[node %d] index %d chunks number %d", id, i, n)
				if n > nodeNum-2*f {
					disperse3[i] = true
					break
				}
				log.Printf("[node %d] index %d chunks %v\n", id, i, receiveChunk3[i])
				time.Sleep(2 * time.Second)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	wg.Add(nodeNum)
	for i := 0; i < nodeNum; i++ {
		go func(i int) {
			var msg Payload
			tmp := make([]ErasureCodeChunk, 0)
			mutex.RLock()
			for _, chunk := range receiveChunk3[i] {
				if chunk != nil {
					tmp = append(tmp, chunk)
				}
				if len(tmp) == nodeNum-2*f {
					break
				}
			}
			mutex.RUnlock()
			codec.Decode(tmp, &msg)
			dp3 := DKGPayload3{}

			json.Unmarshal(msg.([]byte), &dp3)
			//log.Printf("[node %d] retrieve: %v \n", s.ID, dp1)
			aijsumReceived := make([]*big.Int, nodeNum)
			for i := 0; i < nodeNum; i++ {
				aijsumReceived[i] = new(big.Int).SetBytes(dp3.Asum[i])
			}
			dkg.Receiveaijsum(i, aijsumReceived)
			wg.Done()
		}(i)
	}
	wg.Wait()
	for i := 0; i < nodeNum; i++ {
		flag := dkg.FaultDetectPhase2(i)
		log.Printf("[node %d] verify node %v\n", s.ID, flag)
	}
/*
	wg.Add(nodeNum)
	for i := 0; i < nodeNum; i++ {
		go func(i int) {
			mutex.Lock()
			flag := dkg.PKRecVerify(i)
			mutex.Unlock()
			log.Printf("[node %d] verify node %v pk\n", s.ID, flag)
			wg.Done()
		}(i)
	}
	wg.Wait()*/
	dkg.PKRecStep3()
	end := time.Now()
	if !isExist("/home/ubuntu/testdata/") {
		os.MkdirAll("/home/ubuntu/testdata/",os.ModePerm)
	}
	os.Mkdir(fmt.Sprintf("/home/ubuntu/testdata/%d_%d",nodeNum,secretNum), 0777)
	output := fmt.Sprintf("/home/ubuntu/testdata/%d_%d/node%d", nodeNum, secretNum, s.ID)
	file2, _ := os.OpenFile(output, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	file2.Write([]byte(fmt.Sprintf("send %d bytes\n", s.GetAmount())))
	file2.Write([]byte(fmt.Sprintf("bandwidth %d bytes\n", s.GetBandwidth())))
	file2.Write([]byte(fmt.Sprintf("cost time %vs", end.Sub(start).Seconds())))
	file2.Close()
	time.Sleep(2 * time.Second)
}

func isExist(path string) bool {
	_, err := os.Stat(path)    //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func main() {

	N := flag.Int("n", 4, "total nodes number")
	F := flag.Int("f", 1, "number of faulty nodes")
	S := flag.Int("s", 2, "number of secret")
	ID := flag.Int("id", 0, "server id")
	Path := flag.String("path", "/home/ubuntu/", "node info path")
	flag.Parse()
	//generateIPList(*N, *Path)
	if *F == 0 {
		log.Fatalln("F must be greater than 0")
	}
	if *N < *F*3+1 {
		log.Fatalln("N must be greater or equal to 3F+1")
	}
	if *Path == "" {
		log.Fatalln("path of node information is empty")
	}
	var wg sync.WaitGroup
	wg.Add(*N/8)
	for i := 0; i < *N/8; i++ {
		go func(i int){
			test(*N, *S, *F, i+(*ID)*(*N/8), *Path)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
