package cmd

import (
	"fmt"
	"os"

	"github.com/CHTJonas/httkey-server/pkg/utils"
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Hash a URL",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		rawurl := os.Args[2]
		hash, err := utils.RawURLToHash(rawurl)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("%s: %s\n", rawurl, hash)
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateManifestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateManifestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
