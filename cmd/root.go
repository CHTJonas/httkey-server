package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var path string

var rootCmd = &cobra.Command{
	Use:   "httkey",
	Short: "HTTP server that serves files based on the hash of the URL",
	Long: "httkey is a web server that serves static files in response to HTTP requests. " +
		"The static files are all stored in a single directory with filenames that are a hash of their URL." +
		"The hash used is a 128-bit MurmurHash3 of the request path, keyed on the HTTP Host.",
}

func Execute(v string) {
	version = v
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "directory from which to serve files (default current directory)")
}

func initConfig() {
	if len(path) == 0 {
		var err error
		path, err = os.Getwd()
		if err != nil {
			log.Panicln(err)
		}
	}
}
