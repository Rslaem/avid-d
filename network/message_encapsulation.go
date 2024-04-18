package network

import (
	"log"

	"github.com/QinYuuuu/avid-d/protobuf"
	"google.golang.org/protobuf/proto"
)

func Encapsulation(messageType string, ID []byte, sender uint32, payloadMessage any) *protobuf.Message {
	var data []byte
	var err error
	switch messageType {
	case "Z":
		data, err = proto.Marshal(&protobuf.Message.(*protobuf.Z))
	}
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}
	return &protobuf.Message{
		MessageType: messageType,
		ID:          ID,
		Sender:      sender,
		Data:        data,
	}
}
func Decapsulation(message *protobuf.Message) any {
	var payloadMessage any
	switch message.MessageType {
	case "Z":
		payloadMessage = &protobuf.Z{}
		proto.Unmarshal(message.Data, payloadMessage.(*protobuf.Z))
	}
	return payloadMessage
}
