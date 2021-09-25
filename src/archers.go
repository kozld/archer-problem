package src

import (
	"context"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/libp2p/go-libp2p"
	"github.com/multiformats/go-multiaddr"
)

type Archers struct {
	//mtx *sync.Mutex
	nodes addrList
}

func NewArchers(len int8) *Archers {
	return &Archers{
		nodes: make(addrList, 0, len),
	}
}

// addArcher
//func (a *Archers) addArcher(addr *peer.AddrInfo) {
//	a.mtx.Lock()
//	defer a.mtx.Unlock()
//
//	a.nodes = append(a.nodes, addr)
//}

// InitArcher initializes archer
func (a *Archers) InitArcher() {
	ctx, cancel := context.WithCancel(context.Background())

	host, err := libp2p.New(ctx)
	if err != nil {
		panic(err)
	}
	defer host.Close()

	fmt.Println("Addresses:", host.Addrs())
	fmt.Println("ID:", host.ID())

	err = a.nodes.Set(fmt.Sprintf("%s/p2p/%s", host.Addrs()[0].String(), host.ID().String()))
	if err != nil {
		panic(err)
	}

	fmt.Println(a.nodes)

	dht, err := NewKDHT(ctx, host, a.nodes)
	if err != nil {
		log.Fatal(err)
	}

	go Discover(ctx, host, dht, "ldej/echo")

	//addr, err := multiaddr.NewMultiaddr(host.Addrs()[1].String())
	//if err != nil {
	//	panic(err)
	//}
	//addrInfo, err := peer.AddrInfoFromP2pAddr(addr)
	//if err != nil {
	//	panic(err)
	//}
	//a.addArcher(addrInfo)

	//sigCh := make(chan os.Signal)
	//signal.Notify(sigCh, syscall.SIGKILL, syscall.SIGINT)
	//<-sigCh

	run(host, cancel)
}

func run(h host.Host, cancel func()) {
	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-c

	fmt.Printf("\rExiting...\n")

	cancel()

	if err := h.Close(); err != nil {
		panic(err)
	}
	os.Exit(0)
}

//func (a *Archers) GetNodes() []*peer.AddrInfo {
//	return a.nodes
//}

type addrList []multiaddr.Multiaddr

func (al *addrList) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (al *addrList) Set(value string) error {
	addr, err := multiaddr.NewMultiaddr(value)
	if err != nil {
		return err
	}
	*al = append(*al, addr)
	return nil
}
