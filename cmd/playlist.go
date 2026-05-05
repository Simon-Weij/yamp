package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var playlistCmd = &cobra.Command{
	Use:   "playlist",
	Aliases: []string{"pl"},
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("playlist called")
	},
}

func init() {
	rootCmd.AddCommand(playlistCmd)

}
