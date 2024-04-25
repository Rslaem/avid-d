<<<<<<< HEAD
package vid

import (
	//"bytes"
	//"encoding/binary"
	//"encoding/gob"
	"fmt"

	. "github.com/QinYuuuu/avid-d/erasurecode"
)

// VIDMessage is the message emitted and handled by the VID.
type VIDMessage struct {
	IndexID      int  // index of vid instance
	Got          bool // true if this is a Got message
	Ready        bool // true if this is a Ready message
	Disperse     bool // true if this is a Disperse message; a Disperse message contains the dispersed chunk
	RequestChunk bool // true if this message requests a chunk of the dispersed file
	RespondChunk bool // true if this message responds with a chunk request; such a message contains a dispersed chunk
	Cancel       bool
	PayloadChunk ErasureCodeChunk // the chunk of the dispersed file
	Checksum     Checksum         // the checksum
	DestID       int              // destination of the message
	FromID       int              // source of the message
}

func (m VIDMessage) Index() int {
	return m.IndexID
}

// Dest returns the destination of the message.
func (m VIDMessage) Dest() int {
	return m.DestID
}

// From returns the source of the message.
func (m VIDMessage) From() int {
	return m.FromID
}

// Size returns the size of the message in the emulator. It is equal to the size of the PayloadChunk plus AssociatedChunk.
func (m VIDMessage) Size() int {
	totSize := 0
	if m.PayloadChunk != nil {
		totSize += m.PayloadChunk.(Sizer).Size()
	}
	if m.Checksum.Size() != 0 {
		totSize += m.Checksum.Size()
	}
	return totSize
}

// String formats the VIDMessage for debug output.
func (m VIDMessage) String() string {
	t := ""
	if m.Got {
		t += "Got"
	}
	if m.Ready {
		t += "Ready"
	}
	if m.Disperse {
		t += "Disperse"
	}
	if m.Cancel {
		t += "Cancel"
	}
	if m.RequestChunk {
		t += "Request"
	}
	if m.RespondChunk {
		t += "Respond"
	}
	return fmt.Sprintf("%v from node %d in vid %d", t, m.FromID, m.IndexID)
}
=======
package vid

import (
	//"bytes"
	//"encoding/binary"
	//"encoding/gob"
	"fmt"
	"TMAABE/erasurecode"
)

// VIDMessage is the message emitted and handled by the VID.
type VIDMessage struct {
	IndexID      int  // index of vid instance
	Echo         bool // true if this is an Echo message; an Echo message contains the broadcasted chunk
	Ready        bool // true if this is a Ready message
	Disperse     bool // true this is a Disperse message; a Disperse message contains the dispersed chunk
	RequestChunk bool // true if this message requests a chunk of the dispersed file
	RespondChunk bool // true if this message responds with a chunk request; such a message contains a dispersed chunk
	Cancel       bool
	PayloadChunk erasurecode.ErasureCodeChunk // the chunk of the dispersed file
	Checksum     Checksum                     // the checksum
	DestID       int                          // destination of the message
	FromID       int                          // source of the message
}

func (m VIDMessage) Index() int {
	return m.IndexID
}

// Dest returns the destination of the message.
func (m VIDMessage) Dest() int {
	return m.DestID
}

// From returns the source of the message.
func (m VIDMessage) From() int {
	return m.FromID
}

// Size returns the size of the message in the emulator. It is equal to the size of the PayloadChunk plus AssociatedChunk.
func (m VIDMessage) Size() int {
	totSize := 0
	if m.PayloadChunk != nil {
		totSize += m.PayloadChunk.(Sizer).Size()
	}
	if m.Checksum.Size() != 0 {
		totSize += m.Checksum.Size()
	}
	return totSize
}

// String formats the VIDMessage for debug output.
func (m VIDMessage) String() string {
	t := ""
	if m.Echo {
		t += "Echo"
	}
	if m.Ready {
		t += "Ready"
	}
	if m.Disperse {
		t += "Disperse"
	}
	if m.Cancel {
		t += "Cancel"
	}
	if m.RequestChunk {
		t += "Request"
	}
	if m.RespondChunk {
		t += "Respond"
	}
	return fmt.Sprintf("%v from node %d in vid %d", t, m.FromID, m.IndexID)
}
>>>>>>> e982a5d3560d233384b7cc8b8a3b52c93986a5ee
