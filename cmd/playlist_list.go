package cmd

import (
	"fmt"
	"yamp/internal/playlist"

	"github.com/spf13/cobra"
)

var playlistListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List playlists",
	RunE: func(cmd *cobra.Command, args []string) error {
		playlists, err := playlist.ListPlaylists()
		if err != nil {
			return err
		}
		for _, name := range playlists {
			fmt.Println(name)
		}
		return nil
	},
}

func init() {
	playlistCmd.AddCommand(playlistListCmd)
}
