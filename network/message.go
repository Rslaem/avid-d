package network

type Message interface {
	Dest() int // Dest returns the destination of the message
	From() int // From returns the source of the messsge
	Size() int // Size returns the size of the object in the emulator.
}

type HttpMessage struct {
	DataType string
	Content  []byte
}
