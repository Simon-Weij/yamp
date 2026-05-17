package cmd

import (
	"fmt"
	"os"
	playlistcmd "yamp/cmd/playlist"
	"yamp/tui"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "yamp",
	Short:        "Yet another music player",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := tui.RunTUI()
		if err != nil {
			return fmt.Errorf("could not run tui %w", err)
		}
		return nil
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
