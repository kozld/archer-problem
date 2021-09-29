package rpc

import (
	"context"
)

const (
	MessageService     = "MessageRPCAPI"
	MessageServiceFunc = "Message"
)

type MessageRPCAPI struct {
	service *RPCService
}

type Command struct {
	Cmd  string
	Args []interface{}
}

func (m *MessageRPCAPI) Message(ctx context.Context, in Command, out *Command) error {
	*out = m.service.ReceiveMessage(in)
	return nil
}
