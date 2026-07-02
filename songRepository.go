package main

type SongRepository struct{}

func NewSongRepository() *SongRepository {
	return &SongRepository{}
}

func (sr *SongRepository) CurrentlyPlaying() PlaylistItem {
	return PlaylistItem{
		Artist:   "Artist 1",
		Album:    "Album 1",
		Title:    "Title 1",
		Cover:    "https://placehold.co/400x400",
		Duration: 5000,
	}
}
