package vid

import (
	//"encoding/json"
	"fmt"
	"log"

	//"encoding/binary"
	//"encoding/gob"
	//"bytes"

	. "github.com/QinYuuuu/avid-d/erasurecode"
	"github.com/QinYuuuu/avid-d/hasher"
)

// VIDPayload is the interface that a payload of the VID protocol should implement.
type VIDPayload interface{}

// We are going type embedding here. This is to help us remove unneeded states when the time
// comes. For example, we can remove all received chunks after decoding the payload.

// VIDOutput is the output of a VID instance.
type VIDOutput struct {
	payload           VIDPayload // the object being dispersed in the VID, nil if it is not yet retrieved or decoded
	Checksum          Checksum
	ourChunk          ErasureCodeChunk // the erasure coded chunk of the dispersed file that should be held by us, nil if not received
	requestUnanswered []bool           // nil if we are answering chunk requests right away; otherwise, true if we have not answered the chunk
	// request previously sent by the corresponding node
	canReleaseChunk bool // if we are allowed to release the payload chunk
	canceled        bool
}

/*
type VIDRetrieve struct {
	sentRetrieve   bool
	retrieveChunks []StoredErasureCodeChunk // the chunks that we have received, or nil if we have not received the chunk from that server
}
*/
// VIDPayloadState is the execution state of the VID that is related to decoding the dispersed file.
type VIDPayloadState struct {
	chunks           []ErasureCodeChunk // the chunks that we have received, or nil if we have not received the chunk from that server
	nChunks          int                // the number of chunks that we have received, which should equal the number of non-nil items in chunks
	payloadScheduled bool               // if we want to request chunks
	sentRequest      bool               // if we have requested payload chunks
}

// VIDDisperseState is the core execution state of the VID Disperse.
type VIDDisperseState struct {
	echos []bool // the chunks of the broadcasted file that we have received in Echo messages,
	// or nil if we have not received Echo from that server
	nEchos    int    // the number of echos that we have received, which should equal the number of non-nil items in echos
	gots      []bool // the number of gots we have received
	readys    []bool // if we received ready from that server
	nReadys   int    // the number of readys we received, which should equal the number of true's in readys
	sentEcho  bool   // if we have sent out Echo
	sentReady bool   // if we have sent out Ready
}

type VID struct {
	initID  int // ID of the node which should disperse the file
	IndexID int
	*VIDPayloadState
	*VIDDisperseState
	*VIDOutput
	codec ErasureCode // the codec we want to use
	ProtocolParams
}

// NewVID constructs a new VID instance.
func NewVID(initID int, indexID int, p ProtocolParams, codec ErasureCode) *VID {
	v := VID{
		ProtocolParams:   p,
		IndexID:          indexID,
		VIDDisperseState: &VIDDisperseState{},
		VIDPayloadState:  &VIDPayloadState{},
		VIDOutput:        &VIDOutput{},
		codec:            codec,
	}
	v.initID = initID
	v.chunks = make([]ErasureCodeChunk, p.N)
	v.gots = make([]bool, p.N)
	v.readys = make([]bool, p.N)
	v.echos = make([]bool, p.N)
	for i := 0; i < p.N; i++ {
		v.echos[i] = false
		v.readys[i] = false
	}
	v.nEchos = 0
	v.nReadys = 0
	v.requestUnanswered = make([]bool, p.N)

	// if we are supposed to disperse, set payload to default to nothing
	if v.ID == v.initID {
		v.payload = nil
	}
	return &v
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
	/*
		content := make([][]byte, len(pldChunks))
		for i, c := range pldChunks {
			log.Printf("[node %d] generate chunks %v", v.ID, c)
			content[i] = hasher.SHA256Hasher(c.GetData())
		}*/
	checksum := Checksum{
		Value: [][]byte{hasher.SHA256Hasher([]byte("A"))},
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
		msg.Checksum = checksum
		msgs = append(msgs, msg)
	}
	v.sentReady = true
	log.Printf("[node %d] dispersing chunks", v.ID)
	return msgs, 0
}

// handleEcho processes an Echo message with the given source and broadcasted chunk. It panics if the broadcasted chunk is nil.
// It is a nop if VIDCoreState is nil.
/*
func (v *VID) handleEcho(from int, c ErasureCodeChunk, ad Checksum) {
	log.Printf("[node %d] handling Echo from node %d", v.ID, from)
	if v.VIDDisperseState == nil {
		log.Printf("[node %d] VIDDisperseState is nil", v.ID)
		// if the core state is dropped, it means that we no longer need to handle echo
		return
	}
	if c == nil {
		panic("handling echo message with empty chunk")
	}
	// record the message, and we only take the first message
	//fmt.Printf("v.echos %v \n", v.echos)
	if !v.echos[from] {
		v.echos[from] = true
		v.nEchos += 1
		log.Printf("[node %d] nEchos %d\n", v.ID, v.nEchos)
	}
	//fmt.Printf("v.chunks %v \n", v.chunks)
	if !v.chunks[from].IsStored {
		v.nChunks += 1
		v.chunks[from].Store(c, v.DBPath)
		log.Printf("[node %d] store chunk from %d\n", v.ID, from)
		log.Printf("[node %d] nChunks %d\n", v.ID, v.nChunks)
	}
}
*/
// handleDisperse processes a Disperse message from the given source with the given dispersed and broadcasted chunk. It panics if
// either chunk is nil. It is a nop if VIDCoreState is nil, or if the source is not the initID of the VID.
func (v *VID) handleDisperse(from int, c ErasureCodeChunk, checksum Checksum) {
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
			//fmt.Printf("has stored chunks %v\n", v.chunks)
		}
		v.ourChunk = c
		log.Printf("[node %d] has stored its own chunk from node %d", v.ID, from)
	}
}

// handleGot processes a Got message. If we have received n-t Got, broadcast Ready
func (v *VID) handleGot(from int) {
	if v.VIDDisperseState == nil {
		// if the core state is dropped, it means that we no longer need to handle ready
		return
	}
	if !v.gots[from] {
		log.Printf("[node %d] receive ready message from node %d\n", v.ID, from)

		v.readys[from] = true
		v.nReadys += 1
		log.Printf("[node %d] nReadys is %d\n", v.ID, v.nReadys)
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

// handleChunkResponse handles a Response message from the given source and dispersed chunk. It is a nop if VIDPayloadState
// is nil, and panics if the dispersed chunk is nil.
func (v *VID) handleChunkResponse(from int, c ErasureCodeChunk) {
	if v.VIDPayloadState == nil {
		return
	}
	if c == nil {
		panic("handling chunk response message with nil payloadChunk")
	}

	// record the chunk and we only take the first message
	if !v.chunks[from].Stored() {
		v.nChunks += 1
		v.chunks[from].Store(c, v.DBPath)
	}
	v.Printf("receiving chunk from node %v\n", from)
}

// respondRequest handles a Request message from the given source and returns a slice of messages to be sent as the response.
// If we are allowed to respond to the request, the response is sent right away. Otherwise, we record the request and return.

func (v *VID) respondRequest(from, ourid int) []Message {
	var msgs []Message
	log.Printf("[node %d] handling request from node %d", v.ID, from)
	// if we can respond to chunk requests
	if v.canReleaseChunk && v.ourChunk.Stored() {
		msg := &VIDMessage{}
		msg.RespondChunk = true
		msg.IndexID = v.IndexID
		msg.FromID = ourid
		msg.DestID = from
		msg.PayloadChunk = v.ourChunk
		msgs = append(msgs, msg)
	} else {
		// otherwise buffer the response
		v.requestUnanswered[from] = true
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
	var msgs []Message

	// handle the message
	/*
		if m.Echo {
			v.handleEcho(m.FromID, m.PayloadChunk, m.Checksum)

			// if we have requested the dispersed file, but have not sent the requests, see if we can send
			// we send requests when we have got N-F Echos
			if v.payload == nil {
				if v.payloadScheduled {
					msg := &VIDMessage{}
					msg.RequestChunk = true
					msg.IndexID = v.IndexID
					msg.DestID = m.FromID
					msg.FromID = v.ID
					msgs = append(msgs, msg)
					v.Printf("requesting chunk from node %v\n", m.FromID)
				}
			}
		}*/

	if m.Got {
		v.handleGot(m.FromID)
	}
	if m.Ready {
		v.handleReady(m.FromID)
	}
	if m.Disperse {
		v.handleDisperse(m.FromID, m.PayloadChunk, m.Checksum)
	}

	if m.RespondChunk {
		v.handleChunkResponse(m.FromID, m.PayloadChunk)
	}
	if m.RequestChunk {
		msgs = append(msgs, v.respondRequest(m.FromID, v.ID)...)
		return msgs, 0
	}
	// from now on, the message is not used anymore

	// logics that happens when we are not terminated
	if v.VIDDisperseState != nil {
		// if we have received 2F+1 Echos, send out Ready
		if v.nEchos >= (2*v.F+1) && !v.sentReady {
			for i := 0; i < v.N; i++ {
				msg := &VIDMessage{}
				msg.Ready = true
				msg.IndexID = v.IndexID
				msg.FromID = v.ID
				msg.DestID = i

				msgs = append(msgs, msg)
			}
			v.sentReady = true
			v.Println("sending out Ready due to enough Echos")
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
			v.Println("sending out Ready due to enough Readys")
		}

		// if we have got our chunks, send out Echo
		if !v.sentEcho && v.ourChunk.Stored() && v.ourEcho.Stored() {
			for i := 0; i < v.N; i++ {
				msg := &VIDMessage{}
				msg.Echo = true
				msg.FromID = v.ID
				msg.IndexID = v.IndexID
				msg.DestID = i
				msg.PayloadChunk = v.ourChunk
				msg.Checksum = v.ourEcho.Load()
				msgs = append(msgs, msg)
			}
			v.sentEcho = true
			v.Printf("[node %d] sending out Echos", v.ID)
		}
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
	// if we have got N-2F chunks, decode the dispersed file
	if v.VIDPayloadState != nil {
		if v.payload == nil && v.nChunks > v.N-v.F*2 {
			// collect the chunks
			chunks := make([]ErasureCodeChunk, v.N-v.F*2)
			collected := 0
			for _, val := range v.chunks {
				log.Printf("chunk:%v \n", val)
				if val.Stored() {
					chunks[collected] = val.Load()
					collected += 1
				}
				if collected >= v.N-v.F*2 {
					break
				}
			}
			if collected < v.N-v.F*2 {
				log.Printf("echos:%v \n", v.echos)
				log.Printf("chunks:%v \n", v.chunks)
				panic("insufficient shards")
			}
			// decode the dispersed file
			var espayload Payload
			err := v.codec.Decode(chunks, &espayload)
			if err != nil {
				panic(err)
			}
			v.payload = espayload.(VIDPayload)
			v.Printf("decoding payload %s", v.payload)
		}
		// delete payload state now that we have decoded the payload
		// note that we can't move this into the IF above, because the initiating node will never enter the IF above,
		// because it does not need to decode in order to obtain the payload
		if v.payload != nil && v.nReadys >= (2*v.F+1) {
			// on the disk
			v.VIDPayloadState = nil
			if !v.canceled {
				log.Printf("[node %d] sending out cancel", v.ID)
				msgs = append(msgs, v.sendOutCancel()...)
				v.canceled = true
			}
		}
	}

	// See if we can remove the Disperse state. We can do it when the dispersed file have been requested (we need
	// the core state to know who have sent us Echos, so that we know who to request chunks from), and after
	// the protocol is terminated.
	// TODO: currently, we are removing it only after decoding the payload
	if (v.VIDPayloadState == nil) && v.Terminated() {
		// BUG(leiy): We are not deleting the echos (StoredErasureCodeChunk)
		// on the disk
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
	return v.nReadys >= v.F*2+1
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
	if v.requestUnanswered != nil && v.ourChunk.Stored() {
		for from, t := range v.requestUnanswered {
			if t {
				msg := &VIDMessage{}
				msg.RespondChunk = true
				msg.FromID = v.ID
				msg.IndexID = v.IndexID
				msg.DestID = from
				msg.PayloadChunk = v.ourChunk.Load()
				msgs = append(msgs, msg)
			}
		}
		// we don't need the buffer anymore
		v.requestUnanswered = nil
	}
	return msgs
}

// RequestPayload schedules the VID to request the dispersed file, and returns a slice of messages to be sent. If more than N-F
// nodes have sent us Echo, it sends out these requests right away. Otherwise, the request will be sent upon receiving N-F Echos.
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
		v.Printf("requesting chunk from node %v\n", i)
	}
	return msgs
}
