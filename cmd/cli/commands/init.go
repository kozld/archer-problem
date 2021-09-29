package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"

	"github.com/stdi0/archer-problem/src/models"
	"github.com/stdi0/archer-problem/src/p2p_controller"
)

const (
	// below are the keys that will be recorded in the archer's memory:

	// ArcherAddressKey is archer's own address
	ArcherAddressKey = "address"
	// LNeighborAddressKey is archer's left neighbor address
	LNeighborAddressKey = "left_neighbor"
	// RNeighborAddressKey is archer's right neighbor address
	RNeighborAddressKey = "right_neighbor"
)

// init
func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a squad of archers",
	Run: func(cmd *cobra.Command, args []string) {
		var archers = make([]*models.Archer, 0, length)
		var peerAddresses = make([]multiaddr.Multiaddr, 0, length)

		// create a squad of archers
		var neighborPeer multiaddr.Multiaddr
		for i := 0; i < length; i++ {
			c, h := p2p_controller.NewP2PController(neighborPeer)
			archer := models.NewArcher(c)
			go archer.Burn()
			archers = append(archers, archer)

			addr := fmt.Sprintf("%s/p2p/%s", h.Addrs()[0].String(), h.ID().String())
			fmt.Printf("[DEBUG] %s\n", addr)
			multiAddr, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				panic(err)
			}
			peerAddresses = append(peerAddresses, multiAddr)
			neighborPeer = multiAddr
		}

		// get to know the neighbors
		for id, archer := range archers {
			archer.SaveToMemory(ArcherAddressKey, peerAddresses[id])
			if id > 0 {
				archer.SaveToMemory(LNeighborAddressKey, peerAddresses[id-1])
			}
			if id < len(archers)-1 {
				archer.SaveToMemory(RNeighborAddressKey, peerAddresses[id+1])
			}
		}

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		<-ch
	},
}
