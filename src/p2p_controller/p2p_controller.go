package p2p_controller

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/stdi0/archer-problem/src/models"
	"github.com/stdi0/archer-problem/src/p2p_controller/dht"
	"github.com/stdi0/archer-problem/src/p2p_controller/discover"
	"github.com/stdi0/archer-problem/src/p2p_controller/rpc"
)

///////////////////////////////////////
// P2P CONTROLLER IMPLEMENTATION
///////////////////////////////////////

// NewP2PController can construct new instance of Controller
func NewP2PController(neighborPeer multiaddr.Multiaddr) (models.Controller, host.Host) {
	ctx, cancel := context.WithCancel(context.Background())
	host, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}

	r := rpc.NewRPCService(host, "/p2p/rpc/archers")
	ctrl := &P2PController{
		Host:         host,
		Cancel:       cancel,
		NeighborPeer: neighborPeer,
		RPC:          r,
	}

	handler := rpc.NewHandler(ctrl)
	if err = r.Setup(handler); err != nil {
		panic(err)
	}

	return ctrl, host
}

// P2PController implements the Controller interface
type P2PController struct {
	Host         host.Host
	Cancel       context.CancelFunc
	NeighborPeer multiaddr.Multiaddr
	RPC          *rpc.RPCService
}

// Message can send command to another Archer
func (c *P2PController) Message(to models.Archer, message interface{}) {
	command, ok := message.(rpc.Command)
	if !ok {
		panic("error interface cast")
	}
	valueFromMemory := to.GetFromMemory("address")
	if multiAddr, ok := valueFromMemory.(multiaddr.Multiaddr); !ok {
		panic("error interface cast")
	} else {
		addrInfo, err := peer.AddrInfoFromP2pAddr(multiAddr)
		if err != nil {
			panic(err)
		}
		out := c.RPC.Message(addrInfo.ID, command)
		fmt.Printf("Got out: %s", out)
	}
}

// Fire can to fire
func (c *P2PController) Fire() {
	fmt.Println("Fire from controller")
}

// Start can start controller
func (c *P2PController) Start() {
	if c.NeighborPeer != nil {
		dht, err := dht.NewDHT(context.Background(), c.Host, c.NeighborPeer)
		if err != nil {
			log.Fatal(err)
		}
		go discover.Discover(context.Background(), c.Host, dht, "archers/msg")
	}
}

// Stop can stop controller
func (c *P2PController) Stop() error {
	c.Cancel()
	return c.Host.Close()
}
