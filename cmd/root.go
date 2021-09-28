package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/stdi0/archer-problem/src/models"
	"github.com/stdi0/archer-problem/src/p2p_controller"
)

const (
	// DefaultLength is default number of archers
	DefaultLength = 10
	// DefaultTimeout is default timeout in seconds after program
	// started to run messaging
	DefaultTimeout = 3

	///////////////////////////////////////////////////////////////////
	// below are the keys that will be recorded in the archer's memory
	///////////////////////////////////////////////////////////////////

	// ArcherAddressKey is archer's own address
	ArcherAddressKey = "address"
	// LNeighborAddressKey is archer's left neighbor address
	LNeighborAddressKey = "left_neighbor"
	// RNeighborAddressKey is archer's right neighbor address
	RNeighborAddressKey = "right_neighbor"
)

var (
	cfgFile string
	// length is number of archers
	length int
	// timeout in seconds after program started to run messaging
	timeout int

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "archer",
		Short: "Short about archer",
		Long:  `Long about archer`,
		Run:   run,
	}
)

func run(cmd *cobra.Command, args []string) {
	//lineOfArchers := models.NewLineOfArchers()
	var neighborPeer multiaddr.Multiaddr
	for i := 0; i < length; i++ {
		c, h := p2p_controller.NewP2PController(neighborPeer)

		addr := fmt.Sprintf("%s/p2p/%s", h.Addrs()[0].String(), h.ID().String())
		fmt.Printf("[DEBUG] %s\n", addr)

		multiAddr, err := multiaddr.NewMultiaddr(addr)
		if err != nil {
			panic(err)
		}

		a := models.NewArcher(c)
		go a.Burn()

		a.SaveToMemory(ArcherAddressKey, multiAddr)

		//lineOfArchers.AddArcher(a)
		neighborPeer = multiAddr
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().IntVarP(&length, "length", "l", DefaultLength, "number of archers")
	rootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", DefaultTimeout, "timeout in seconds after program started to run messaging")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".archer-problem" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".archer-problem")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
