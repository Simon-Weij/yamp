package playlistcmd

import (
	"fmt"
	"regexp"
	"strings"
	"yamp/internal/musicdiscovery"
	"yamp/internal/play"
	"yamp/internal/playlist"

	"github.com/spf13/cobra"
)

var playlistSimilarToCmd = &cobra.Command{
	Use:   "similar-to",
	Short: "A brief description of your command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		songs, err := getSimilarSongsFn(args[0])
		if err != nil {
			return fmt.Errorf("could not find similar songs: %w", err)
		}
		var songsLines []string
		for _, line := range strings.Split(strings.TrimSpace(songs.Stdout), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			}
			songsLines = append(songsLines, line)
			fmt.Println(cleanLines(line))
		}
		playlistName := "similar-to-" + args[0]
		if err := addSongsToPlaylistFn(songsLines, playlistName); err != nil {
			return fmt.Errorf("could not add songs to playlist %w", err)
		}

		fmt.Printf("finished playlist %s", playlistName)

		return nil
	},
}

func addSongsToPlaylist(songs []string, playlistName string) error {
	isInternal := true
	if err := createPlaylistSimilarToFn(playlistName, isInternal); err != nil {
		return fmt.Errorf("could not create playlist: %w", err)
	}
	for _, song := range songs {
		cleaned := cleanLines(song)
		metadata, err := extractMetadataFn(cleaned)
		if err != nil {
			return fmt.Errorf("could not get metadata for %s", cleaned)
		}
		filepath := convertSongMetadataToFilePathFn(metadata.Artist, metadata.Album, metadata.Title)
		fmt.Printf("downloading %s \n", cleaned)
		if err := downloadSongFn(cleaned, filepath); err != nil {
			return fmt.Errorf("could not download %s: %w", cleaned, err)
		}
		fmt.Printf("finished downloading %s \n", cleaned)

		location := convertSongMetadataToFilePathFn(metadata.Artist, metadata.Album, metadata.Title)
		if err := addItemToPlaylistFn(playlistName, metadata.Artist, metadata.Title, location, isInternal); err != nil {
			return fmt.Errorf("could not add %s to playlist: %w", cleaned, err)
		}
	}
	return nil
}

func cleanLines(value string) string {
	re := regexp.MustCompile(`\([^)]*\)|\[[^\]]*\]`)
	stripped := re.ReplaceAllString(value, "")
	stripped = strings.ReplaceAll(stripped, "\"", "")
	return strings.TrimSpace(stripped)
}

func init() {
	playlistCmd.AddCommand(playlistSimilarToCmd)
}

var (
	getSimilarSongsFn              = musicdiscovery.GetSimilarSongs
	addSongsToPlaylistFn           = addSongsToPlaylist
	createPlaylistSimilarToFn      = playlist.CreatePlaylist
	extractMetadataFn              = musicdiscovery.ExtractMetadata
	convertSongMetadataToFilePathFn = playlist.ConvertSongMetadataToFilePath
	downloadSongFn                 = play.DownloadSong
	addItemToPlaylistFn            = playlist.AddItemToPlaylist
)
