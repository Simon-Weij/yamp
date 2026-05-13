package cmd

import (
	"os"
	playlistcmd "yamp/cmd/playlist"
	"yamp/tui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "yamp",
	Short:        "Yet another music player",
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		tui.RunTUI()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(playlistcmd.Command())
}
