package src

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/protocol"
	"github.com/libp2p/go-libp2p-gorpc"
)

type Service struct {
	rpcServer *rpc.Server
	rpcClient *rpc.Client
	host      host.Host
	protocol  protocol.ID
	counter   int
}

func NewService(host host.Host, protocol protocol.ID) *Service {
	return &Service{
		host:     host,
		protocol: protocol,
	}
}

func (s *Service) SetupRPC() error {
	s.rpcServer = rpc.NewServer(s.host, s.protocol)

	messageRPCAPI := MessageRPCAPI{service: s}
	err := s.rpcServer.Register(&messageRPCAPI)
	if err != nil {
		return err
	}

	s.rpcClient = rpc.NewClientWithServer(s.host, s.protocol, s.rpcServer)
	return nil
}

func (s *Service) StartMessaging(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.counter++
			s.Message(fmt.Sprintf("Message (%d): Hello from %s", s.counter, s.host.ID().Pretty()))
		}
	}
}

func (s *Service) Message(message string) {
	peers := s.host.Peerstore().Peers()
	var replies = make([]*Envelope, len(peers))

	errs := s.rpcClient.MultiCall(
		Ctxts(len(peers)),
		peers,
		MessageService,
		MessageServiceFunc,
		Envelope{Message: message},
		CopyEnvelopesToIfaces(replies),
	)

	for i, err := range errs {
		if err != nil {
			fmt.Printf("Peer %s returned error: %-v\n", peers[i].Pretty(), err)
		} else {
			fmt.Printf("Peer %s echoed: %s\n", peers[i].Pretty(), replies[i].Message)
		}
	}
}

func (s *Service) ReceiveMessage(envelope Envelope) Envelope {
	return Envelope{Message: fmt.Sprintf("Peer %s echoing: %s", s.host.ID(), envelope.Message)}
}

func Ctxts(n int) []context.Context {
	ctxs := make([]context.Context, n)
	for i := 0; i < n; i++ {
		ctxs[i] = context.Background()
	}
	return ctxs
}

func CopyEnvelopesToIfaces(in []*Envelope) []interface{} {
	ifaces := make([]interface{}, len(in))
	for i := range in {
		in[i] = &Envelope{}
		ifaces[i] = in[i]
	}
	return ifaces
}
