package cmd

import (
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
)

var (
	log = logging.MustGetLogger("cmd")
	// VERSION is the current version of the bot
	VERSION string
)

var rootCmd = &cobra.Command{
	Use:   "mu2",
	Short: "Mu2 is a discord bot.",
	Long:  "",
}

// Execute runs the bot
func Execute(version string) error {
	VERSION = version
	return rootCmd.Execute()
}
