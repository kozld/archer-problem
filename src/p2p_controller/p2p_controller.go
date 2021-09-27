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

	return &P2PController{host, cancel, neighborPeer}, host
}

type P2PController struct {
	Host         host.Host
	Cancel       context.CancelFunc
	NeighborPeer multiaddr.Multiaddr
}

func (c *P2PController) MessageTo(neighbor models.Archer) {

}

func (c *P2PController) Fire() {

}

func (c *P2PController) Start() {
	//addr := fmt.Sprintf("%s/p2p/%s", c.N.Addrs()[0].String(), c.Host.ID().String())
	//multiAddr, err := multiaddr.NewMultiaddr(addr)
	//if err != nil {
	//	panic(err)
	//}

	//addrInfo, _ := peer.AddrInfoFromP2pAddr(multiAddr)
	//addrInfo.ID

	//this := a.nodes.PushBack(multiAddr)
	//neighbor := this.Prev()

	if c.NeighborPeer != nil {
		//neighborAddr, ok := neighbor.Value.(multiaddr.Multiaddr)
		//if !ok {
		//	panic("error interface cast")
		//}

		dht, err := NewDHT(context.Background(), c.Host, c.NeighborPeer)
		if err != nil {
			log.Fatal(err)
		}

		//service := src.NewService(host, protocol.ID("/p2p/rpc/archers"))
		//err = service.SetupRPC()
		//if err != nil {
		//	log.Fatal(err)
		//}

		go Discover(context.Background(), c.Host, dht, "archers/msg")
		//go service.Start(ctx)
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
