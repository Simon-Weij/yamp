package playlistcmd

import (
	"fmt"
	"os"
	"yamp/internal/playlist"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove-song",
	Aliases: []string{"rm-song", "remove-track"},
	Args:    cobra.ExactArgs(2),
	Short:   "Remove a song from a playlist",
	Run: func(cmd *cobra.Command, args []string) {
		artist := args[0]
		title := args[1]
		playlistName, err := cmd.Flags().GetString("playlist")
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := playlist.RemoveItemFromPlaylist(playlistName, artist, title); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
}

func init() {
	playlistCmd.AddCommand(removeCmd)

	removeCmd.Flags().String("playlist", "p", "Playlist name")
	_ = removeCmd.MarkFlagRequired("playlist")
}
