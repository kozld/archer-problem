package rpc

import "context"

const (
	MessageService     = "MessageRPCAPI"
	MessageServiceFunc = "Message"
)

type MessageRPCAPI struct {
	service *RPCService
}

//type Envelope struct {
//	Message string
//}

type Message string

func (m *MessageRPCAPI) Message(ctx context.Context, in Message, out *Message) error {
	*out = m.service.ReceiveMessage(in)
	return nil
}
