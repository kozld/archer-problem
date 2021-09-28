package p2p_controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/stdi0/archer-problem/src/models"
)

///////////////////////////////////////
/////////// P2P CONTROLLER ////////////
///////////////////////////////////////

func NewP2PController(neighborPeer multiaddr.Multiaddr) (models.Controller, host.Host) {
	ctx, cancel := context.WithCancel(context.Background())
	host, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}
	rpc := NewRPCService(host, "/p2p/rpc/archers")

	return &P2PController{host, cancel, neighborPeer, rpc}, host
}

type P2PController struct {
	Host         host.Host
	Cancel       context.CancelFunc
	NeighborPeer multiaddr.Multiaddr
	RPC          *RPCService
}

func (c *P2PController) MessageTo(neighbor models.Archer, message string) {
	multiAddrIface := neighbor.GetFromMemory("address")
	if multiAddr, ok := multiAddrIface.(multiaddr.Multiaddr); !ok {
		panic("error interface cast")
	} else {
		addrInfo, err := peer.AddrInfoFromP2pAddr(multiAddr)
		if err != nil {
			panic(err)
		}
		c.RPC.Message(addrInfo.ID, message)
	}
}

func (c *P2PController) Fire() {

}

func (c *P2PController) Start() {
	if c.NeighborPeer != nil {
		dht, err := NewDHT(context.Background(), c.Host, c.NeighborPeer)
		if err != nil {
			log.Fatal(err)
		}
		go Discover(context.Background(), c.Host, dht, "archers/msg")
	}

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	fmt.Printf("\rExiting...\n")

	c.Cancel()

	if err := c.Stop; err != nil {
		panic(err)
	}
	os.Exit(0)
}

func (c *P2PController) Stop() error {
	return c.Host.Close()
}
