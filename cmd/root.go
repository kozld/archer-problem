package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"

	"github.com/stdi0/archer-problem/src"
)

const (
	// DefaultLength is default number of archers
	DefaultLength = 10
	// DefaultTimeout is default timeout in seconds after program
	// started to run messaging
	DefaultTimeout = 3
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
	length  int8
	// timeout in seconds after program started to run messaging
	timeout int8
)

func run(cmd *cobra.Command, args []string) {
	archers := src.NewArchers(length)
	go archers.InitArcher()
	archers.InitArcher()
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

	rootCmd.PersistentFlags().Int8VarP(&length, "length", "l", DefaultLength, "number of archers")
	rootCmd.PersistentFlags().Int8VarP(&timeout, "timeout", "t", DefaultTimeout, "timeout in seconds after program started to run messaging")
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
