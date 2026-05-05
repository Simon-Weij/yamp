package cmd

import (
	"fmt"
	"yamp/internal/playlist"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Aliases: []string{"songs", "tracks", "items"},
	Args:  cobra.ExactArgs(1),
	Short: "List all songs in an playlist",
	RunE: func(cmd *cobra.Command, args []string) error {
		items, err := playlist.ListPlaylistItems(args[0])
		if err != nil {
			return fmt.Errorf("could not list playlist items: %w", err)
		}
		for i, item := range items {
			fmt.Printf("%d. %s - %s \n", i+1, item.Artist, item.Title)
		}
		return nil
	},
}

func init() {
	playlistCmd.AddCommand(listCmd)
}
