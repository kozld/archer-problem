package src

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	disc "github.com/libp2p/go-libp2p-discovery"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/multiformats/go-multiaddr"
)

var wg sync.WaitGroup

func NewKDHT(ctx context.Context, host host.Host, bootstrapPeers []multiaddr.Multiaddr) (*disc.RoutingDiscovery, error) {
	var options []dht.Option

	if len(bootstrapPeers) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
	}

	kdht, err := dht.New(ctx, host, options...)
	if err != nil {
		return nil, err
	}

	if err = kdht.Bootstrap(ctx); err != nil {
		return nil, err
	}

	for _, peerAddr := range bootstrapPeers {
		fmt.Println(peerAddr)
		peerinfo, err := peer.AddrInfoFromP2pAddr(peerAddr)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := host.Connect(ctx, *peerinfo); err != nil {
				log.Printf("Error while connecting to node %q: %-v", peerinfo, err)
			} else {
				log.Printf("Connection established with bootstrap node: %q", *peerinfo)
			}
		}()
	}
	wg.Wait()

	return disc.NewRoutingDiscovery(kdht), nil
}
