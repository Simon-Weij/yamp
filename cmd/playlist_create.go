package cmd

import (
	"fmt"
	"yamp/internal/playlist"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := playlist.CreatePlaylist(args[0]); err != nil {
			return fmt.Errorf("could not create playlist: %w", err)
		}
		fmt.Printf("successfully created playlist %s \n", args[0])
		return nil
	},
}

func init() {
	playlistCmd.AddCommand(createCmd)
}
