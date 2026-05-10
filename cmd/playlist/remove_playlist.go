package playlistcmd

import (
	"fmt"
	"yamp/internal/playlist"

	"github.com/spf13/cobra"
)

var removePlaylistCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm", "rm-playlist"},
	Args:    cobra.ExactArgs(1),
	Short:   "Remove a playlist",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := playlist.DeletePlaylist(args[0]); err != nil {
			return fmt.Errorf("could not remove playlist: %w", err)
		}
		fmt.Printf("successfully removed playlist %s \n", args[0])
		return nil
	},
}

func init() {
	playlistCmd.AddCommand(removePlaylistCmd)
}
