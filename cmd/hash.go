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
	Long: "Takes a URL as an argument and outputs its hash. " +
		"You can use this to determine what the filenames of your static content should be.",
	Run: func(cmd *cobra.Command, args []string) {
		hash, err := utils.RawURLToHash(args[0])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Printf("%s: %s\n", args[0], hash)
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
}
