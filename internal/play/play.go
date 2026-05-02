package play

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/lrstanley/go-ytdlp"
)

func DownloadSong(name string, output string) error {
	ytdlp.MustInstall(context.TODO(), nil)
	ytdlp.MustInstallFFmpeg(context.TODO(), nil)

	dl := ytdlp.New().ExtractAudio().AudioFormat("mp3").Verbose().ParseMetadata("title:%(artist)s - %(title)s").EmbedMetadata().Output(filepath.Join(output))

	out, err := dl.Run(context.TODO(), "ytsearch1:"+name)

	if err != nil {
		return fmt.Errorf("something went wrong while downloading song! %w with the following logs: %s", err, out)
	}
	return nil
}

func PlaySong(path string) error {
	cmd := exec.Command("mpv", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not start mpv %w", err)
	}
	return nil
}
