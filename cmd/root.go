package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "httkey",
	Short: "HTTP server that serves files based on the hash of the URL",
	Long:  `httkey is a web server that serves files with a filename based on a 128-bit MurmurHash3 of the URL.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	// TODO
}
