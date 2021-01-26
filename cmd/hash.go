package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/CHTJonas/httkey-server/pkg/utils"
	"github.com/spf13/cobra"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: "Hash a URL",
	Long: "Takes a URL as an argument and outputs its hash. " +
		"You can use this to determine what the filenames of your static content should be.",
	Run: func(cmd *cobra.Command, args []string) {
		s := args[0]
		if _, err := utils.ParseURL(s); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		s = strings.TrimPrefix(s, "http://")
		s = strings.TrimPrefix(s, "https://")
		s = "http://" + s
		hash, err := utils.RawURLToHash(s)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		s = fmt.Sprintf("%s\t%s", s, hash)
		fmt.Println(s)
	},
}

func init() {
	rootCmd.AddCommand(hashCmd)
}
