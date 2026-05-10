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
	info, err := os.Stat(output)
	if err == nil {
		if !info.IsDir() {
			return nil
		}
		return fmt.Errorf("output path is a directory: %s", output)
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("could not check output path %s: %w", output, err)
	}

	ytdlpInstall()
	ytdlpInstallFFmpeg()

	out, err := ytdlpDownload(name, output)

	if err != nil {
		logOutput := ""
		if out != nil {
			logOutput = out.Stdout
		}
		return fmt.Errorf("something went wrong while downloading song! %w with the following logs: %s", err, logOutput)
	}
	return nil
}

func PlaySong(path string) error {
	cmd := execCommand("mpv", path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("could not start mpv %w", err)
	}
	return nil
}

var (
	ytdlpInstall = func() {
		ytdlp.MustInstall(context.TODO(), nil)
	}
	ytdlpInstallFFmpeg = func() {
		ytdlp.MustInstallFFmpeg(context.TODO(), nil)
	}
	ytdlpDownload = func(name, output string) (*ytdlp.Result, error) {
		dl := ytdlp.New().ExtractAudio().AudioFormat("mp3").Verbose().ParseMetadata("title:%(artist)s - %(title)s").EmbedMetadata().Output(filepath.Join(output))
		return dl.Run(context.TODO(), "ytsearch1:"+name)
	}
	execCommand = exec.Command
)