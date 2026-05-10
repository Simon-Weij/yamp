package playlistcmd

import (
	"fmt"
	"yamp/internal/playlist"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove-song",
	Aliases: []string{"rm-song", "remove-track"},
	Args:    cobra.ExactArgs(2),
	Short:   "Remove a song from a playlist",
	RunE: func(cmd *cobra.Command, args []string) error {
		artist := args[0]
		title := args[1]
		playlistName, err := cmd.Flags().GetString("playlist")
		if err != nil {
			return err
		}

		if err := removeItemFromPlaylistFn(playlistName, artist, title); err != nil {
			return fmt.Errorf("could not remove song: %w", err)
		}
		return nil
	},
}

func init() {
	playlistCmd.AddCommand(removeCmd)

	removeCmd.Flags().String("playlist", "p", "Playlist name")
	_ = removeCmd.MarkFlagRequired("playlist")
}

var removeItemFromPlaylistFn = playlist.RemoveItemFromPlaylist
