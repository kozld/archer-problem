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

	ArcherAddressKey = "address"
)

var (
	cfgFile string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "archer-problem",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: run,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
	}

	// length is number of archers
	length int
	// timeout in seconds after program started to run messaging
	timeout int
)

func run(cmd *cobra.Command, args []string) {
	lineOfArchers := models.NewLineOfArchers()

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
		a.SaveToMemory(ArcherAddressKey, multiAddr)
		go a.Burn()

		lineOfArchers.AddArcher(a)
		neighborPeer = multiAddr
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	//archers := models.NewArchers(length)
	//go archers.InitArcher()
	//archers.InitArcher()
	//time.Sleep(3 * time.Second)
	//for _, addr := range archers.GetNodes() {
	//	fmt.Printf("ADDR: %s\n", addr)
	//}
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
