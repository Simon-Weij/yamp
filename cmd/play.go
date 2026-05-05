package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"yamp/internal/play"

	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:   "play",
	Aliases: []string{"p"},
	Short: "Play a song by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		songName := args[0]
		tempdir := filepath.Join(os.TempDir(), "yamp")
		if err := os.MkdirAll(tempdir, 0700); err != nil {
			return fmt.Errorf("couldn't create directory %s", tempdir)
		}
		finalPath := filepath.Join(tempdir, "song.mp3")
		play.DownloadSong(songName, finalPath)

		play.PlaySong(finalPath)

		if err := os.RemoveAll(tempdir); err != nil {
			return fmt.Errorf("could not clean up properly %w ", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
}
