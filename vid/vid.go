package vid

import (
	"encoding/json"

	"github.com/QinYuuuu/avid-d/commit/merklecommitment"
	"github.com/QinYuuuu/avid-d/hasher"

	//"encoding/json"
	"fmt"
	"log"

	//"encoding/binary"
	//"encoding/gob"
	//"bytes"
	"github.com/QinYuuuu/avid-d/erasurecode"
	escode "github.com/QinYuuuu/avid-d/erasurecode"
)

// TODO: we should pass the pointer to the on-disk shard in the message, and let the
// message encoder to retrieve the data.

// VIDPayload is the interface that a payload of the VID protocol should implement.
type VIDPayload interface{}

// We are going type embedding here. This is to help us remove unneeded states when the time
// comes. For example, we can remove all received chunks after decoding the payload.

// VIDOutput is the output of a VID instance.
type VIDOutput struct {
	payload           VIDPayload              // the object being dispersed in the VID, nil if it is not yet retrieved or decoded
	ourChunk          escode.ErasureCodeChunk // the erasure coded chunk of the dispersed file that should be held by us, nil if not received
	ourChecksum       StoredChecksum          // the erasure coded chunk of the broadcasted file that should be held by us, nil if not received
	requestUnanswered []bool                  // nil if we are answering chunk requests right away; otherwise, true if we have not answered the chunk
	// request previously sent by the corresponding node
	canReleaseChunk bool // if we are allowed to release the payload chunk
	canceled        bool
}

// VIDPayloadState is the execution state of the VID that is related to decoding the dispersed file.
type VIDPayloadState struct {
	chunks           []escode.ErasureCodeChunk // the chunks that we have received, or nil if we have not received the chunk from that server
	nChunks          int                       // the number of chunks that we have received, which should equal the number of non-nil items in chunks
	payloadScheduled bool                      // if we want to request chunks
	sentRequest      bool                      // if we have requested payload chunks
}

// VIDDisperseState is the core execution state of the VID Disperse.
type VIDDisperseState struct {
	gots []bool // the chunks of the broadcasted file that we have received in Echo messages,
	// or nil if we have not received Echo from that server
	nGots     int    // the number of gots that we have received, which should equal the number of non-nil items in gots
	readys    []bool // if we received ready from that server
	nReadys   int    // the number of readys we received, which should equal the number of true's in readys
	sentGot   bool   // if we have sent out Echo
	sentReady bool   // if we have sent out Ready
}

type VIDRetrieveState struct {
	flag          bool
	setRoot       bool                         // if we have set the merkle tree root
	Root          *merklecommitment.MerkleTree // the root hash
	proof         merklecommitment.Witness     // the proof of the chunk with root R
	rootGot       []bool                       // if we received Root(r) from that server
	nRootGots     int                          // the number of Root(r) we have received
	rootReady     []bool                       // if we received Ready(r) from that server
	nRootReadys   int                          // the number of Ready(r) we have received
	sentRoot      bool                         // if we have sent out Root(r)
	sentRootReady bool                         // if we have sent out Ready(r)
}

type VID struct {
	initID  int // ID of the node which should disperse the file
	IndexID int // ID of the vid instance
	*VIDPayloadState
	*VIDDisperseState
	*VIDRetrieveState
	*VIDOutput
	codec erasurecode.ErasureCode // the codec we want to use
	ProtocolParams
}

// SerializePayload serializes a Payload into a []byte
func SerializePayload(p VIDPayload) []byte {
	b, err := json.Marshal(p)
	if err != nil {
		log.Fatalf("Failed to serialize payload: %v", err)
	}
	return b
}

// NewVID constructs a new VID instance.
func NewVID(initID int, indexID int, p ProtocolParams, codec erasurecode.ErasureCode) *VID {
	v := VID{
		ProtocolParams:   p,
		IndexID:          indexID,
		VIDDisperseState: &VIDDisperseState{},
		VIDRetrieveState: &VIDRetrieveState{},
		VIDPayloadState:  &VIDPayloadState{},
		VIDOutput:        &VIDOutput{},
		codec:            codec,
	}
	v.initID = initID
	v.chunks = make([]escode.ErasureCodeChunk, p.N)
	v.readys = make([]bool, p.N)
	v.gots = make([]bool, p.N)
	for i := 0; i < p.N; i++ {
		v.gots[i] = false
		v.readys[i] = false
	}
	v.nGots = 0
	v.nReadys = 0

	v.ourChunk = nil
	v.requestUnanswered = make([]bool, p.N)

	v.sentReady = false
	v.sentGot = false

	// if we are supposed to disperse, set payload to default to nothing
	if v.ID == v.initID {
		v.payload = nil
	}
	return &v
}
func (v *VID) InitRetrieve() {
	v.rootReady = make([]bool, v.N)
	v.rootGot = make([]bool, v.N)
	for i := 0; i < v.N; i++ {
		v.rootGot[i] = false
		v.rootReady[i] = false
	}
	v.nRootGots = 0
	v.nRootReadys = 0

	v.sentRoot = false
	v.sentRootReady = false
	v.setRoot = false

	v.VIDRetrieveState.flag = true
	v.canReleaseChunk = true
}

// Init executes the initialization procedure of the protocol and returns new messages and updates
// Init starts the dispersion of this VID. It is a nop if the caller is not the node which is supposed to initiate the VID.
func (v *VID) Init() ([]Message, Event) {
	log.Printf("[node %d] Init a VID instance %d", v.ID, v.IndexID)
	var msgs []Message
	// do nothing if we are not supposed to disperse
	if v.initID != v.ID {
		log.Printf("[node %d] do nothing because not supposed to disperse", v.ID)
		return msgs, 0
	}
	// encode the payload and the associated data
	pldChunks, err := v.codec.Encode(v.payload)
	log.Printf("[node %d] encode VID payload %s", v.ID, v.payload)
	if err != nil {
		panic("error encoding payload " + err.Error())
	}

	// Serialize the payload and create the checksum
	var data [][]byte
	for _, p := range pldChunks {
		data = append(data, SerializePayload(p))
	}
	mt, err := merklecommitment.NewMerkleTree(data, hasher.SHA256Hasher)
	if err != nil {
		log.Fatalf("Failed to create Merkle tree: %v", err)
	}
	var witnesses []merklecommitment.Witness
	for i := 0; i < len(data); i++ {
		witness, err := merklecommitment.CreateWitness(mt, i)
		if err != nil {
			log.Fatalf("Failed to create Merkle witness: %v", err)
		}
		witnesses = append(witnesses, witness)
	}
	// send out Disperse and Ready messages
	// can't do Echo here because both Echo and Disperse use the AssociatedChunk field
	for i := 0; i < v.N; i++ {
		msg := &VIDMessage{}
		msg.Disperse = true
		msg.Ready = true
		msg.IndexID = v.IndexID
		msg.FromID = v.ID
		msg.DestID = i
		msg.PayloadChunk = pldChunks[i]
		msg.Checksum = Checksum{Value: witnesses[i], Root: mt}
		msgs = append(msgs, msg)
	}
	v.sentReady = true
	log.Printf("[node %d] dispersing chunks", v.ID)
	return msgs, 0
}

// handleEcho processes an Echo message with the given source and broadcasted chunk. It panics if the broadcasted chunk is nil.
// It is a nop if VIDCoreState is nil.
func (v *VID) handleGot(from int) {
	log.Printf("[node %d] handling Got from node %d", v.ID, from)
	if v.VIDDisperseState == nil {
		log.Printf("[node %d] VIDDisperseState is nil", v.ID)
		// if the core state is dropped, it means that we no longer need to handle echo
		return
	}

	// record the message, and we only take the first message
	if !v.gots[from] {
		v.gots[from] = true
		v.nGots += 1
		log.Printf("[node %d] nGots %d\n", v.ID, v.nGots)
	}
}

// handleReady processes a Ready message from the given source. It is a nop if VIDCoreState is nil.
func (v *VID) handleReady(from int) {
	if v.VIDDisperseState == nil {

		// if the core state is dropped, it means that we no longer need to handle ready
		return
	}
	// record the message
	if !v.readys[from] {
		log.Printf("[node %d] receive ready message from node %d\n", v.ID, from)
		v.readys[from] = true
		v.nReadys += 1
		log.Printf("[node %d] nReadys is %d\n", v.ID, v.nReadys)
	}
}

// handleDisperse processes a Disperse message from the given source with the given dispersed and broadcasted chunk. It panics if
// either chunk is nil. It is a nop if VIDCoreState is nil, or if the source is not the initID of the VID.
func (v *VID) handleDisperse(from int, c erasurecode.ErasureCodeChunk, checksum Checksum) (flag int) {
	log.Printf("[node %d] handling disperse message from %d\n", v.ID, from)
	if v.VIDDisperseState == nil {
		log.Printf("VIDDisperseState is nil \n")
		// it means that we have terminated, and decoded associated data
		return
	}
	if from != v.initID {
		// not the one who is supposed to send chunk
		return
	}

	if c == nil {
		panic("handling disperse message with nil payloadChunk")
	}
	/*
		if checksum.Size() == 0 {
			panic("handling disperse message with empety checksum")
		}*/

	// record the message, and we only take the first message
	if v.ourChunk == nil {
		// we need to check VIDPayloadState here, because the initiating node will hear its own Disperse message, and at
		// that time, the VIDPayloadState may already been marked as nil
		if v.VIDPayloadState != nil {
			v.chunks[v.ID] = c
			v.nChunks += 1
			log.Printf("[node %d] nChunks is %d\n", v.ID, v.nChunks)
			fmt.Printf("has stored chunks %v\n", v.chunks)
		}

		v.ourChunk = c
		log.Printf("[node %d] has stored its own chunk from node %d", v.ID, from)
	}
	if !v.ourChecksum.IsStored {
		// see the comment above for why we need to check for VIDCoreState here
		if v.VIDDisperseState != nil {
			v.gots[v.ID] = true
			v.nGots += 1
			log.Printf("[node %d] nGots is %d\n", v.ID, v.nGots)
			//fmt.Printf("stored checksums %v\n", v.gots)
		}
		result, _ := merklecommitment.Verify(checksum.Root.Root().Hash(), checksum.Value, SerializePayload(v.ourChunk), hasher.SHA256Hasher)
		if result != false {
			v.ourChecksum.Store(checksum)
			v.proof = checksum.Value

			flag = 1
			log.Printf("[node %d] has stored its own checksum from node %d", v.ID, from)
		} else {
			flag = 0
			log.Printf("[node %d] has received wrong checksum from node %d", v.ID, from)
		}
	}
	return flag
}

func (v *VID) handleRootReady(from int) {
	if v.VIDRetrieveState == nil {
		return
	}

	if !v.rootReady[from] {
		log.Printf("[node %d] receive Ready(r) from node %d\n", v.ID, from)
		v.rootReady[from] = true
		v.nRootReadys += 1
		log.Printf("[node %d] nRootReadys is %d\n", v.ID, v.nRootReadys)

	}
}

func (v *VID) handleRoot(from int) {
	log.Printf("[node %d] handling Root(r) from node %d", v.ID, from)
	if v.VIDRetrieveState == nil {
		log.Printf("[node %d] VIDTrieveState is nil", v.ID)
		// if the core state is dropped, it means that we no longer need to handle echo
		return
	}
	//if c == nil {
	//	panic("handling echo message with empty chunk")
	//}

	if !v.rootGot[from] {
		v.rootGot[from] = true
		v.nRootGots += 1
		log.Printf("[node %d] nRootGots %d\n", v.ID, v.nRootGots)
	}
	//fmt.Printf("v.chunks %v \n", v.chunks)
}

// handleChunkResponse handles a Response message from the given source and dispersed chunk. It is a nop if VIDPayloadState
// is nil, and panics if the dispersed chunk is nil.
func (v *VID) handleChunkResponse(from int, c erasurecode.ErasureCodeChunk, witness Checksum) {
	if v.VIDPayloadState == nil {
		return
	}
	if c == nil {
		panic("handling chunk response message with nil payloadChunk")
	}

	// record the chunk and we only take the first message
	// TODO: check the merkle proof ...
	judge, _ := merklecommitment.Verify(witness.Root.Root().Hash(), witness.Value, SerializePayload(c), hasher.SHA256Hasher)

	if (v.chunks[from] == nil) && judge {
		v.nChunks += 1
		v.chunks[from] = c
		log.Printf("[node %d] store chunk from %d\n", v.ID, from)
		log.Printf("[node %d] nChunks %d\n", v.ID, v.nChunks)
	}
	fmt.Printf("receiving chunk from node %v\n", from)
}

// respondRequest handles a Request message from the given source and returns a slice of messages to be sent as the response.
// If we are allowed to respond to the request, the response is sent right away. Otherwise, we record the request and return.

func (v *VID) respondRequest(from int) []Message {
	var msgs []Message
	log.Printf("[node %d] handling request from node %d", v.ID, from)
	// if we can respond to chunk requests
	judge, _ := merklecommitment.Verify(v.ourChecksum.Root.Root().Hash(), v.ourChecksum.Value, SerializePayload(v.ourChunk), hasher.SHA256Hasher)

	if !v.setRoot && judge == true {
		if !v.sentRoot {
			for i := 0; i < v.N; i++ {
				msg := &VIDMessage{}
				msg.RootGot = true
				msg.IndexID = v.IndexID
				msg.FromID = v.ID
				msg.DestID = i
				//msg.PayloadChunk = v.ourChunk
				//msg.Checksum = Checksum{v.ourChecksum.Value, v.ourChecksum.Root}
				msgs = append(msgs, msg)
			}
			v.sentRoot = true
		} else {
			//向invoker发送元组消息
			msg := &VIDMessage{}
			msg.ToInvoker = from
			msg.RespondChunk = true
			msg.IndexID = v.IndexID
			msg.FromID = v.ID
			msg.DestID = from
			msg.PayloadChunk = v.ourChunk
			msg.Checksum = Checksum{v.ourChecksum.Value, v.ourChecksum.Root}
			msgs = append(msgs, msg)
		}
	} else if judge == false {
		for i := 0; i < v.N; i++ {
			msg := &VIDMessage{}
			msg.RootGotPerp = true
			msg.IndexID = v.IndexID
			msg.FromID = v.ID
			msg.DestID = i
			msgs = append(msgs, msg)
		}
	}
	if v.setRoot && judge == true {
		//向invoker发送元组消息
		msg := &VIDMessage{}
		msg.ToInvoker = from
		msg.RespondChunk = true
		msg.IndexID = v.IndexID
		msg.FromID = v.ID
		msg.DestID = from
		msg.PayloadChunk = v.ourChunk
		msg.Checksum = Checksum{v.ourChecksum.Value, v.ourChecksum.Root}
		msgs = append(msgs, msg)
	}
	return msgs
}

func (v *VID) sendOutCancel() []Message {
	var msgs []Message
	for i := 0; i < v.N; i++ {
		msg := &VIDMessage{}
		msg.IndexID = v.IndexID
		msg.FromID = v.ID
		msg.DestID = i
		msg.Cancel = true
		msgs = append(msgs, msg)
	}
	return msgs
}

func (v *VID) IfCanceled() bool {
	return v.canceled
}

// Recv handles a VID message of type *VIDMessage and returns a slice of messages to send and the execution result.
func (v *VID) Recv(mg Message) ([]Message, Event) {
	m, ok := mg.(*VIDMessage)
	fmt.Printf("recv message %v\n", mg)
	if !ok {
		log.Println("interface message convert to VIDMessage failed")
	}
	if m.IndexID != v.IndexID {
		log.Panic("wrong vid instance")
	}
	//fmt.Printf("message %v %v %v\n", m, m.PayloadChunk, m.Checksum)
	var msgs []Message

	// handle the message
	if m.Got {
		v.handleGot(m.FromID)

		// if we have requested the dispersed file, but have not sent the requests, see if we can send
		// we send requests when we have got N-F gots
		/*if v.payload == nil {
			if v.payloadScheduled {
				msg := &VIDMessage{}
				msg.RequestChunk = true
				msg.IndexID = v.IndexID
				msg.DestID = m.FromID
				msg.FromID = v.ID
				msgs = append(msgs, msg)

				fmt.Printf("requesting chunk from node %v\n", m.FromID)
			}
		}*/
	}
	if m.Ready {
		v.handleReady(m.FromID)
	}
	if m.Disperse {
		result := v.handleDisperse(m.FromID, m.PayloadChunk, m.Checksum)
		if result == 1 {
			for i := 0; i < v.N; i++ {
				msg := &VIDMessage{}
				msg.Got = true
				msg.IndexID = v.IndexID
				msg.FromID = v.ID
				msg.DestID = i
				msgs = append(msgs, msg)
			}
		}
	}
	if m.RootGot {
		v.handleRoot(m.FromID)
	}
	if m.RootReady {
		v.handleRootReady(m.FromID)
	}

	if m.RespondChunk {
		v.handleChunkResponse(m.FromID, m.PayloadChunk, m.Checksum)
		//if v.payload == nil {
		//	if v.payloadScheduled {
		//		msg := &VIDMessage{}
		//		msg.RequestChunk = true
		//		msg.IndexID = v.IndexID
		//		msg.DestID = m.FromID
		//		msg.FromID = v.ID
		//		msgs = append(msgs, msg)
		//		fmt.Printf("requesting chunk from node %v\n", m.FromID)
		//	}
		//}

	}
	if m.RequestChunk {
		msgs = append(msgs, v.respondRequest(m.FromID)...)
		return msgs, 0
	}

	// from now on, the message is not used anymore

	// logics that happens when we are not terminated
	if v.VIDDisperseState != nil {
		// if we have received 2F+1 gots, send out Ready
		if v.nGots >= (2*v.F+1) && !v.sentReady {
			for i := 0; i < v.N; i++ {
				msg := &VIDMessage{}
				msg.Ready = true
				msg.IndexID = v.IndexID
				msg.FromID = v.ID
				msg.DestID = i
				msgs = append(msgs, msg)
			}
			v.sentReady = true
			fmt.Println("sending out Ready due to enough gots")
		}

		// if we have received F+1 Readys, send out Ready (if we have not done so)
		if v.nReadys >= (v.F+1) && !v.sentReady {
			for i := 0; i < v.N; i++ {
				msg := &VIDMessage{}
				msg.Ready = true
				msg.FromID = v.ID
				msg.IndexID = v.IndexID
				msg.DestID = i
				msgs = append(msgs, msg)
			}
			v.sentReady = true
			fmt.Println("sending out Ready due to enough Readys")
		}

		// if we have got our chunks, send out Echo
		/*
			if !v.sentGot && v.ourChunk != nil && v.ourChecksum.Stored() {
				for i := 0; i < v.N; i++ {
					msg := &VIDMessage{}
					msg.Got = true
					msg.FromID = v.ID
					msg.IndexID = v.IndexID
					msg.DestID = i
					msgs = append(msgs, msg)
				}
				v.sentGot = true
				fmt.Printf("[node %d] sending out gots", v.ID)
			}*/
	}
	/*	// if we can answer chunk requests, answer the recorded ones now
		if v.requestUnanswered != nil && v.canReleaseChunk && v.ourChunk.Stored() {
			for from, t := range v.requestUnanswered {
				if t {
					msg := &VIDMessage{}
					msg.RespondChunk = true
					msg.FromID = v.ID
					msg.DestID = from
					msg.PayloadChunk = v.ourChunk.LoadPointer(v.DBPath)
					msgs = append(msgs, msg)
				}
			}
			// now that we have answered the requests, we don't need the record anymore
			v.requestUnanswered = nil
		}
	*/

	if v.VIDRetrieveState.flag != false {
		if v.nRootGots >= (2*v.F+1) && !v.sentRootReady {
			for i := 0; i < v.N; i++ {
				msg := &VIDMessage{}
				msg.RootReady = true
				msg.IndexID = v.IndexID
				msg.FromID = v.ID
				msg.DestID = i
				msgs = append(msgs, msg)
			}
			v.sentRootReady = true
			fmt.Println("sending out Root due to enough Ready")
		}

		if v.nRootReadys >= (v.F+1) && !v.sentRootReady {
			for i := 0; i < v.N; i++ {
				msg := &VIDMessage{}
				msg.RootReady = true
				msg.FromID = v.ID
				msg.IndexID = v.IndexID
				msg.DestID = i
				msgs = append(msgs, msg)
			}
			v.sentRootReady = true
			fmt.Println("sending out Ready due to enough RootReadys")
		}
		if v.nRootReadys >= 2*v.F+1 && !v.setRoot {
			v.setRoot = true
		}
	}

	// if we have got N-2F chunks, decode the dispersed file
	if v.VIDRetrieveState.flag != false {
		fmt.Println("*********************************************************")
		//v.payload = nil
		if v.payload == nil && v.nChunks > v.N-v.F*2 {
			// collect the chunks
			chunks := make([]erasurecode.ErasureCodeChunk, v.N-v.F*2)
			collected := 0
			for _, val := range v.chunks {
				log.Printf("chunk:%v \n", val)
				if val != nil {
					chunks[collected] = val
					collected += 1
				}
				if collected >= v.N-v.F*2 {
					break
				}
			}
			if collected < v.N-v.F*2 {
				log.Printf("gots:%v \n", v.gots)
				log.Printf("chunks:%v \n", v.chunks)
				panic("insufficient shards")
			}
			// decode the dispersed file
			var espayload erasurecode.Payload
			err := v.codec.Decode(chunks, &espayload)
			if err != nil {
				panic(err)
			}
			v.payload = espayload.(VIDPayload)
			fmt.Printf("decoding payload %s", v.payload)
			// on the disk
			v.VIDPayloadState = nil
			if !v.canceled {
				log.Printf("[node %d] sending out cancel", v.ID)
				msgs = append(msgs, v.sendOutCancel()...)
				v.canceled = true
			}

		}
		// delete payload state now that we have decoded the payload
		// note that we can't move this into the IF above, because the initiating node will never enter the IF above,
		// because it does not need to decode in order to obtain the payload
		if v.payload == nil && v.nRootReadys >= (2*v.F+1) {
			v.Root = v.ourChecksum.Root
		}
	}

	// See if we can remove the Disperse state. We can do it when the dispersed file have been requested (we need
	// the core state to know who have sent us gots, so that we know who to request chunks from), and after
	// the protocol is terminated.
	// TODO: currently, we are removing it only after decoding the payload
	if (v.VIDPayloadState == nil) && v.Terminated() {
		// BUG(leiy): We are not deleting the gots (StoredErasureCodeChunk)
		v.VIDDisperseState = nil
	}
	if v.Terminated() {
		log.Printf("[node %d] when terminated the VIDPayloadState %v the VIDCoreState %v \n", v.ID, v.VIDPayloadState == nil, v.VIDDisperseState == nil)
		return msgs, Terminate
	} else {
		return msgs, 0
	}
}

// Terminated checks if the file is successfully dispersed and we have decoded the broadcasted file.
func (v *VID) Terminated() bool {
	dispersed := v.PayloadDispersed()
	return dispersed
}

// PayloadDispersed checks if the files are successfully dispersed, i.e. all nodes will eventually be able to retrieve and files.
// It does not check if the broadcasted file is already decoded by us. For that, the user should use Terminated. It always returns
// true if VIDCoreState is nil.
func (v *VID) PayloadDispersed() bool {
	if v.VIDDisperseState == nil {
		return true
	}
	return v.nRootReadys >= v.F*2+1
}

// Payload returns the dispersed file and true if it is decoded, or nil and false if not.
func (v *VID) Payload() (VIDPayload, bool) {
	pThere := v.payload != nil
	return v.payload, pThere
}

// SetPayload sets the dispersed file of the VID instance. It panics if the caller is not the node supposed to inititate the VID.
func (v *VID) SetPayload(p VIDPayload) {
	if v.ID != v.initID {
		panic("cannot set VID payload when not being the init ID")
	}
	v.payload = p
}

// ReleaseChunk enables the VID to answer requests for its dispersed chunk. If we have got our dispersed chunk,
// it sends out responses to those who have sent us a request before and returns a slice of these messages.
// It is a nop if this function is called before.
func (v *VID) ReleaseChunk() []Message {
	var msgs []Message
	if v.canReleaseChunk {
		return msgs
	}
	v.canReleaseChunk = true
	// if we can answer to the requests
	if v.requestUnanswered != nil && v.ourChunk != nil {
		for from, t := range v.requestUnanswered {
			if t {
				msg := &VIDMessage{}
				msg.RespondChunk = true
				msg.FromID = v.ID
				msg.IndexID = v.IndexID
				msg.DestID = from
				msg.PayloadChunk = v.ourChunk
				msgs = append(msgs, msg)
			}
		}
		// we don't need the buffer anymore
		v.requestUnanswered = nil
	}
	return msgs
}

// RequestPayload schedules the VID to request the dispersed file, and returns a slice of messages to be sent. If more than N-F
// nodes have sent us Echo, it sends out these requests right away. Otherwise, the request will be sent upon receiving N-F gots.
// It is a nop if VIDPayloadState is nil, or if we have requested the dispersed file before.
func (v *VID) RequestPayload() []Message {
	log.Printf("[node %d] generate retrieve request", v.ID)
	// if we have requested before, do nothing
	var msgs []Message
	// if we have not got the payload...
	msgs = append(msgs, v.sendPayloadRequests()...)
	return msgs
}

// sendPayloadRequests send out requests for the dispersed file. It searches by the increasing order of node IDs,
// and sends the requests to the first N-F nodes it found to have sent us Echo. It panics when less than N-F nodes
// have sent us Echo.
func (v *VID) sendPayloadRequests() []Message {
	var msgs []Message
	nr := 0
	for i := 0; i < v.N; i++ {
		// if we have received Echo from idx
		msg := &VIDMessage{}
		msg.RequestChunk = true
		msg.DestID = i
		msg.FromID = v.ID
		msg.IndexID = v.IndexID
		msgs = append(msgs, msg)
		nr += 1
		fmt.Printf("requesting chunk from node %v\n", i)
	}
	return msgs
}
