package models

import (
	"container/list"
)

///////////////////////////////////////
///////// ARCHERS LINKED LIST /////////
///////////////////////////////////////

const (
	ArcherIdKey      = "id"
	ArcherIdFirst    = 0
	ArcherAddressKey = "address"
)

// NewLineOfArchers can construct new instance of LineOfArchers
func NewLineOfArchers() *LineOfArchers {
	return &LineOfArchers{list.New()}
}

// LineOfArchers is a linked list of Archer
type LineOfArchers struct {
	items *list.List
}

// AddArcher can add Archer to list
func (l *LineOfArchers) AddArcher(archer *Archer) {
	l.items.PushBack(archer)
}

func (l *LineOfArchers) FireEveryoneSync() {
	///////////////////////////////////////
	/// STAGE 1
	///////////////////////////////////////

	//el := l.items.Front()
	//firstLeftArcher, ok := el.Value.(Archer)
	//if !ok {
	//	panic("error interface cast")
	//}
	//
	//firstLeftArcher.SaveToMemory(ArcherIdKey, ArcherIdFirst)

	//firstLeftArcher.MessageTo()

	n := ArcherIdFirst
	for e := l.items.Front(); e != nil; e = e.Next() {
		archer, ok := e.Value.(Archer)
		if !ok {
			panic("error interface cast")
		}

		archer.SaveToMemory(ArcherIdKey, n)
		//dataFromMemory := archer.GetFromMemory(ArcherIdKey)
		//archer.MessageTo()

		n += 1
	}

	///////////////////////////////////////
	/// STAGE 2
	///////////////////////////////////////
}

func (l *LineOfArchers) linkArchers() {

}

// InitArcher initializes archer
//func initArcher() {
//	//ctx, cancel := context.WithCancel(context.Background())
//	//
//	//host, err := libp2p.New(ctx)
//	//if err != nil {
//	//	panic(err)
//	//}
//	//defer host.Close()
//
//	//fmt.Println("Addresses:", host.Addrs())
//	//fmt.Println("ID:", host.ID())
//
//	addr := fmt.Sprintf("%s/p2p/%s", host.Addrs()[0].String(), host.ID().String())
//	multiAddr, err := multiaddr.NewMultiaddr(addr)
//	if err != nil {
//		panic(err)
//	}
//
//	//addrInfo, _ := peer.AddrInfoFromP2pAddr(multiAddr)
//	//addrInfo.ID
//
//	this := a.nodes.PushBack(multiAddr)
//	neighbor := this.Prev()
//
//	if neighbor != nil {
//		neighborAddr, ok := neighbor.Value.(multiaddr.Multiaddr)
//		if !ok {
//			panic("error interface cast")
//		}
//
//		dht, err := src.NewDHT(ctx, host, neighborAddr)
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		service := src.NewService(host, protocol.ID("/p2p/rpc/archers"))
//		err = service.SetupRPC()
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		go src.Discover(ctx, host, dht, "archers/msg")
//		//go service.Start(ctx)
//	}
//
//	run(host, cancel)
//}
