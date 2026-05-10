package playlistcmd

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
		playlists, err := listPlaylistsFn()
		if err != nil {
			return err
		}
		for _, name := range playlists {
			fmt.Println(name)
		}
		return nil
	},
}

var listPlaylistsFn = playlist.ListPlaylists

func init() {
	playlistCmd.AddCommand(playlistListCmd)
}
