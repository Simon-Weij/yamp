package playlistcmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"yamp/internal/musicdiscovery"
	"yamp/internal/play"
	"yamp/internal/playlist"

	"github.com/adrg/xdg"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: []string{"a"},
	Args:    cobra.ExactArgs(2),
	Short:   "Add songs to your playlist",
	RunE: func(cmd *cobra.Command, args []string) error {
		playlistName := args[0]
		songName := args[1]

		isInternal := false
		playlistExists, err := playlist.PlaylistExists(playlistName, isInternal)
		if err != nil {
			return fmt.Errorf("something went wrong while checking is playlist exists: %w", err)
		}
		if !playlistExists {
			return fmt.Errorf("playlist doesn't exist")
		}

		id := uuid.NewString()

		songPath := filepath.Join(os.TempDir(), "yamp", id, "song.mp3")

		if err := play.DownloadSong(songName, songPath); err != nil {
			return fmt.Errorf("something went wrong while downloading song: %w", err)
		}

		songLocation, metadata, err := initaliseSong(songName, songPath)
		if err != nil {
			return fmt.Errorf("could not initialise song: %w", err)
		}

		if err := playlist.AddItemToPlaylist(playlistName, metadata.Artist, metadata.Title, songLocation, isInternal); err != nil {
			return fmt.Errorf("could not add song to playlist: %w", err)
		}

		fmt.Printf("Successfully initialised %s \n", songName)

		return nil
	},
}

func initaliseSong(song string, path string) (string, *musicdiscovery.Metadata, error) {
	metadata, err := musicdiscovery.ExtractMetadata(song)
	if err != nil {
		return "", nil, fmt.Errorf("failed to extract metadata: %s ", err)
	}

	songpath := filepath.Join(xdg.UserDirs.Music, "yamp", metadata.Artist, metadata.Album)

	if err := os.MkdirAll(songpath, 0755); err != nil {
		return "", nil, fmt.Errorf("something went wrong while creating directory %s: %w", songpath, err)
	}

	destination := filepath.Join(songpath, fmt.Sprintf("%s.mp3", metadata.Title))

	if err := moveFile(path, destination); err != nil {
		return "", nil, fmt.Errorf("something went wrong while moving file from %s to %s: %w", path, destination, err)
	}

	return destination, metadata, nil
}

func moveFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %v", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Couldn't open dest file: %v", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	if err != nil {
		return fmt.Errorf("Couldn't copy to dest from source: %v", err)
	}

	inputFile.Close()

	err = os.Remove(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't remove source file: %v", err)
	}
	return nil
}

func init() {
	playlistCmd.AddCommand(addCmd)
}

