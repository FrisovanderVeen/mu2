package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Gives the version",
	Long:  "Gives the current version of mu2",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VERSION)
	},
}
