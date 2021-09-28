package p2p_controller

import "context"

const (
	MessageService     = "MessageRPCAPI"
	MessageServiceFunc = "Message"
)

type MessageRPCAPI struct {
	service *RPCService
}

type Envelope struct {
	Message string
}

func (m *MessageRPCAPI) Message(ctx context.Context, in Envelope, out *Envelope) error {
	*out = m.service.ReceiveMessage(in)
	return nil
}
