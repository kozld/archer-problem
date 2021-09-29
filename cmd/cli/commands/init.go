package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"

	"github.com/stdi0/archer-problem/src/models"
	"github.com/stdi0/archer-problem/src/p2p_controller"
)

const (
	// below are the keys that will be recorded in the archer's memory:

	// ArcherAddressKey is archer's own address
	ArcherAddressKey = "address"
	// LNeighborKey is archer's left neighbor
	LNeighborKey = "left_neighbor"
	// RNeighborKey is archer's right neighbor
	RNeighborKey = "right_neighbor"
)

// init
func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "start",
	Short: "Init a squad of archers & doing sync fire",
	Run: func(cmd *cobra.Command, args []string) {
		var archers = make([]*models.Archer, 0, length)

		// create a squad of archers
		var neighborPeer multiaddr.Multiaddr
		for i := 0; i < length; i++ {
			// use p2p controller implementation
			c, h := p2p_controller.NewP2PController(neighborPeer)
			archer := models.NewArcher(c)

			fmt.Println("Burning new archer...")
			archer.Burn()
			fmt.Println("Sleep 2 seconds...")
			time.Sleep(2 * time.Second)

			archers = append(archers, archer)

			addr := fmt.Sprintf("%s/p2p/%s", h.Addrs()[0].String(), h.ID().String())
			fmt.Printf("[DEBUG] %s\n", addr)
			multiAddr, err := multiaddr.NewMultiaddr(addr)
			if err != nil {
				panic(err)
			}
			// save own address to memory
			archer.SaveToMemory(ArcherAddressKey, multiAddr)
			neighborPeer = multiAddr
		}

		// get to know the neighbors
		for id, archer := range archers {
			if id > 0 {
				archer.SaveToMemory(LNeighborKey, archers[id-1])
			}
			if id < len(archers)-1 {
				archer.SaveToMemory(RNeighborKey, archers[id+1])
			}
		}

		// doing sync fire!
		firstLeftArcher := archers[0]
		valueFromMemory := firstLeftArcher.GetFromMemory(RNeighborKey)
		rightNeighbor, ok := valueFromMemory.(*models.Archer)
		if !ok {
			panic("error interface cast")
		}
		firstLeftArcher.Message(*rightNeighbor, "hello")

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		<-ch
	},
}
