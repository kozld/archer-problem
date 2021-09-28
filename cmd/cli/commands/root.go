package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	// length is number of archers
	length int
	// timeout in seconds after program started to run messaging
	timeout int

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "archers",
		Short: "Short about archers",
		Long:  `Long about archers`,
		//Run:   initArchers,
	}
)

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
