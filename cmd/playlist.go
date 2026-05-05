package cmd

import (
	"github.com/spf13/cobra"
)

var playlistCmd = &cobra.Command{
	Use:     "playlist",
	Aliases: []string{"pl"},
	Short:   "Manage playlists",
}

func init() {
	rootCmd.AddCommand(playlistCmd)

}
