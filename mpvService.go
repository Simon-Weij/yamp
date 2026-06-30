package main

import (
	"fmt"

	"github.com/gen2brain/go-mpv"
)

type MpvService struct{}

func NewMpvService() *MpvService {
	return &MpvService{}
}

func PlaySong(song PlaylistItem) error {
	m := mpv.New()
	defer m.TerminateDestroy()

	if err := m.Initialize(); err != nil {
		return err
	}

	// No video
	if err := m.SetOptionString("vo", "null"); err != nil {
		return err
	}

	// Best audio quality
	if err := m.SetOptionString("ytdl-format", "bestaudio/best"); err != nil {
		return err
	}

	// Normalize volume
	if err := m.SetOptionString("af", "lavfi=[loudnorm=I=-16:TP=-1.5:LRA=11]"); err != nil {
		return err
	}

	if err := m.RequestLogMessages("info"); err != nil {
		return err
	}
	if err := m.Command([]string{"loadfile", fmt.Sprintf("ytdl://ytsearch:%s - %s", song.Artist, song.Title)}); err != nil {
		return err
	}

	for {
		e := m.WaitEvent(1)
		switch e.EventID {
		case mpv.EventLogMsg:
			msg := e.LogMessage()
			fmt.Printf("[%s] %s", msg.Prefix, msg.Text)
		case mpv.EventFileLoaded:
			fmt.Println("Now playing:", m.GetPropertyString("media-title"))
		case mpv.EventEnd:
			ef := e.EndFile()
			if ef.Reason != mpv.EndFileRedirect {
				return nil
			}
		case mpv.EventShutdown:
			return nil
		}
	}
}
