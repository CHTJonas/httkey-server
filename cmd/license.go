package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var licenseCmd = &cobra.Command{
	Use:   "license",
	Short: "Print copyright license",
	Long:  "Prints httkey's copyright and licensing information.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Copyright (c) 2021 Charlie Jonas.")
		fmt.Println("This software is released under the BSD 2-Clause License.")
		fmt.Println("Please visit https://github.com/CHTJonas/httkey-server for more information.")
	},
}

func init() {
	rootCmd.AddCommand(licenseCmd)
}
