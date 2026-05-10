package playlistcmd

import (
	"fmt"
	"yamp/internal/playlist"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Aliases: []string{"new", "mk"},
	Short: "Create a new playlist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		isInternal := false
		if err := createPlaylistFn(args[0], isInternal); err != nil {
			return fmt.Errorf("could not create playlist: %w", err)
		}
		fmt.Printf("successfully created playlist %s \n", args[0])
		return nil
	},
}

var createPlaylistFn = playlist.CreatePlaylist

func init() {
	playlistCmd.AddCommand(createCmd)
}
