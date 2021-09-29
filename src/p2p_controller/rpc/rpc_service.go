package rpc

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	rpc "github.com/libp2p/go-libp2p-gorpc"
)

type RPCService struct {
	handler   *Handler
	rpcServer *rpc.Server
	rpcClient *rpc.Client
	host      host.Host
	protocol  protocol.ID
}

func NewRPCService(host host.Host, protocol protocol.ID) *RPCService {
	return &RPCService{
		host:     host,
		protocol: protocol,
	}
}

func (s *RPCService) Setup(handler *Handler) error {
	s.handler = handler
	s.rpcServer = rpc.NewServer(s.host, s.protocol)

	messageRPCAPI := MessageRPCAPI{service: s}
	err := s.rpcServer.Register(&messageRPCAPI)
	if err != nil {
		return err
	}

	s.rpcClient = rpc.NewClientWithServer(s.host, s.protocol, s.rpcServer)
	return nil
}

func (s *RPCService) Message(dest peer.ID, message Command) string {
	var msg Command
	err := s.rpcClient.Call(
		dest,
		MessageService,
		MessageServiceFunc,
		message,
		&msg,
	)

	if err != nil {
		return fmt.Sprintf("Peer %s returned error: %-v\n", dest, err)
	}

	return fmt.Sprintf("Peer %s echoed: %s\n", dest, msg)
}

func (s *RPCService) ReceiveMessage(msg Command) Command {
	s.handler.Handle(msg)
	return msg //Message(fmt.Sprintf("Peer %s echoing: %s", s.host.ID(), msg))
}

//func CopyEnvelopesToIfaces(in []*Envelope) []interface{} {
//	ifaces := make([]interface{}, len(in))
//	for i := range in {
//		in[i] = &Envelope{}
//		ifaces[i] = in[i]
//	}
//	return ifaces
//}
