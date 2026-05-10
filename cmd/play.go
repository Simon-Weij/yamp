package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"yamp/internal/play"

	"github.com/spf13/cobra"
)

var playCmd = &cobra.Command{
	Use:     "play",
	Aliases: []string{"p"},
	Short:   "Play a song by name",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		songName := args[0]
		tempdir := filepath.Join(tempDir(), "yamp")
		if err := mkdirAll(tempdir, 0700); err != nil {
			return fmt.Errorf("couldn't create directory %s", tempdir)
		}
		finalPath := filepath.Join(tempdir, "song.mp3")
		if err := downloadSong(songName, finalPath); err != nil {
			return fmt.Errorf("could not download song: %w", err)
		}

		if err := playSong(finalPath); err != nil {
			return fmt.Errorf("could not play song: %w", err)
		}

		if err := removeAll(tempdir); err != nil {
			return fmt.Errorf("could not clean up properly %w ", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(playCmd)
}

var (
	tempDir      = os.TempDir
	mkdirAll     = os.MkdirAll
	removeAll    = os.RemoveAll
	downloadSong = play.DownloadSong
	playSong     = play.PlaySong
)
