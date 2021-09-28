package rpc

import (
	"fmt"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p-core/protocol"
	rpc "github.com/libp2p/go-libp2p-gorpc"
)

type RPCService struct {
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

func (s *RPCService) Setup() error {
	s.rpcServer = rpc.NewServer(s.host, s.protocol)

	messageRPCAPI := MessageRPCAPI{service: s}
	err := s.rpcServer.Register(&messageRPCAPI)
	if err != nil {
		return err
	}

	s.rpcClient = rpc.NewClientWithServer(s.host, s.protocol, s.rpcServer)
	return nil
}

func (s *RPCService) Message(dest peer.ID, message string) {
	var reply = make([]*Envelope, 1)
	err := s.rpcClient.Call(
		dest,
		MessageService,
		MessageServiceFunc,
		Envelope{Message: message},
		CopyEnvelopesToIfaces(reply),
	)

	if err != nil {
		fmt.Printf("Peer %s returned error: %-v\n", dest, err)
	} else {
		fmt.Printf("Peer %s echoed: %s\n", dest, reply[0].Message)
	}
}

func (s *RPCService) ReceiveMessage(envelope Envelope) Envelope {
	return Envelope{Message: fmt.Sprintf("Peer %s echoing: %s", s.host.ID(), envelope.Message)}
}

func CopyEnvelopesToIfaces(in []*Envelope) []interface{} {
	ifaces := make([]interface{}, len(in))
	for i := range in {
		in[i] = &Envelope{}
		ifaces[i] = in[i]
	}
	return ifaces
}
