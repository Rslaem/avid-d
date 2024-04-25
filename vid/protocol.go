<<<<<<< HEAD
package vid

import (
	//"bytes"
	//"encoding/binary"
	//"fmt"
	"log"
	//"strconv"
	//"strings"
	//"sync"
)

// Sizer is the interface that wraps the Size method.
type Sizer interface {
	Size() int // Size returns the size of the object in the emulator.
}

// Event is an update from the execution of a protocol. Protocols represents the
// updates by setting the individual bits.
type Event int

const (
	Terminate    Event = 1 << iota // the protocol terminated
	StateChanged                   // the internal state of the protocol is changed
	FirstCustomEvent
)

// ProtocolParams contains common parameters of an asynchronous distributed algorithm.
type ProtocolParams struct {
	*log.Logger
	//logPrefix string  // prefix of log messages printed by this instance of protocol
	N        int     // the number of servers in the cluster
	F        int     // the number of faulty servers to tolerate
	ID       int     // the ID of the server running this instance of the protocol
	DBPath   KVStore // the leveldb handler which the protocol can use to store long-term data
	DBPrefix []byte  // the prefix of the database key used by this instance of the protocol
}

// Priority is an enumerator that represents the priority of a message.
// Possible values are "High" and "Low".
type Priority int

const (
	High   Priority = iota // signaling
	Medium                 // dispersal
	Low                    // retrieval
)

// Message is the interface that groups the methods that a message should support.
type Message interface {
	Index() int // the index of VID instance
	Dest() int  // Dest returns the destination of the message
	From() int  // From returns the source of the messsge
	Size() int  // Size returns the size of the object in the emulator.
}

// Protocol is the interface that groups the methods that a protocol implements.
// An asynchronous distributed system is modelled as a bunch of IO automata.
// Each instance of the protocol is an IO automaton, which reacts to messages
// coming from other instances, and optionally sends out new messages and updates
// its internal state. A protocol also supports an Init function, which triggers
// the specific initialization procedure for the protocol. Note that a protocol
// may (and should) react to incoming messages even before the initialization.
type Protocol interface {
	Recv(m Message) ([]Message, Event) // Recv handles an incoming message and returns new messages and updates
	Init() ([]Message, Event)          // Init executes the initialization procedure of the protocol and returns new messages and updates
}

// KVStore is the interface that a key-value storage should implement. The storage is used by the protocols
// to dump large execution states that are rarely used.
type KVStore interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
}
=======
package vid

import (
	//"bytes"
	//"encoding/binary"
	//"fmt"
	"log"
	//"strconv"
	//"strings"
	//"sync"
)

// Sizer is the interface that wraps the Size method.
type Sizer interface {
	Size() int // Size returns the size of the object in the emulator.
}

// Event is an update from the execution of a protocol. Protocols represents the
// updates by setting the individual bits.
type Event int

const (
	Terminate    Event = 1 << iota // the protocol terminated
	StateChanged                   // the internal state of the protocol is changed
	FirstCustomEvent
)

// ProtocolParams contains common parameters of an asynchronous distributed algorithm.
type ProtocolParams struct {
	*log.Logger
	//logPrefix string  // prefix of log messages printed by this instance of protocol
	N        int     // the number of servers in the cluster
	F        int     // the number of faulty servers to tolerate
	ID       int     // the ID of the server running this instance of the protocol
	DBPath   KVStore // the leveldb handler which the protocol can use to store long-term data
	DBPrefix []byte  // the prefix of the database key used by this instance of the protocol
}

// Priority is an enumerator that represents the priority of a message.
// Possible values are "High" and "Low".
type Priority int

const (
	High   Priority = iota // signaling
	Medium                 // dispersal
	Low                    // retrieval
)

// Message is the interface that groups the methods that a message should support.
type Message interface {
	Index() int // the index of VID instance
	Dest() int  // Dest returns the destination of the message
	From() int  // From returns the source of the messsge
	Size() int  // Size returns the size of the object in the emulator.
}

// Protocol is the interface that groups the methods that a protocol implements.
// An asynchronous distributed system is modelled as a bunch of IO automata.
// Each instance of the protocol is an IO automaton, which reacts to messages
// coming from other instances, and optionally sends out new messages and updates
// its internal state. A protocol also supports an Init function, which triggers
// the specific initialization procedure for the protocol. Note that a protocol
// may (and should) react to incoming messages even before the initialization.
type Protocol interface {
	Recv(m Message) ([]Message, Event) // Recv handles an incoming message and returns new messages and updates
	Init() ([]Message, Event)          // Init executes the initialization procedure of the protocol and returns new messages and updates
}

// KVStore is the interface that a key-value storage should implement. The storage is used by the protocols
// to dump large execution states that are rarely used.
type KVStore interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
}
>>>>>>> e982a5d3560d233384b7cc8b8a3b52c93986a5ee
