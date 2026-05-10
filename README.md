# Yamp

An easy to use and simple yet powerful music player and browser, yamp should
support all features you expect from a modern music player, including playlists,
music browsing and more

## Prerequisites

This project is relatively light on dependencies, but for building the project
as explained in [#installation](#installation) you need a couple things

- [Git](#https://git-scm.com/install/)
- [Go](https://go.dev/doc/install)

For runtime dependencies, you only need [mpv](https://mpv.io/installation/) this
project also depends on [yt-dlp](https://github.com/yt-dlp/yt-dlp) and
[ffmpeg](https://ffmpeg.org/download.html), but these should automatically
resolve themselves

## Installation

Currently, we don't package project yet, this is because the project is very
early in development if you want to install the project anyways, you can install
it using the following steps

```bash
# Clone the project 
git clone https://github.com/Simon-Weij/yamp.git

# Build the project
go build -o yamp

# Move the binary to a location in your $PATH
mv yamp /some/other/location/in/path
```

## Usage

- `yamp play <song>` (aliases: `p`): Play a song by name
- `yamp playlist` (aliases: `pl`): Manage playlists
- `yamp playlist add <playlist> <song>` (aliases: `a`): Add songs to your playlist
- `yamp playlist create <name>` (aliases: `new`, `mk`): Create a new playlist
- `yamp playlist list` (aliases: `ls`): List playlists
- `yamp playlist songs <playlist>` (aliases: `tracks`, `items`): List songs in a playlist
- `yamp playlist delete <name>` (aliases: `remove`, `rm`, `rm-playlist`): Remove a playlist
- `yamp playlist remove-song <artist> <title> --playlist <name>` (aliases: `rm-song`, `remove-track`): Remove a song from a playlist
- `yamp playlist similar-to <song>`: Make a playlist with similar songs

## Roadmap

- More flexible playlists, currently they're really limited

- Export playlists from popular platforms for example Spotify or youtube-music

- Write tests

## How Does It Work?

### The Play Command

When the play command is hit yt-dlp is asked to download the song to a temporary
directory, after that the song is played with mpv

### Playlists

For playlists the song is first downloaded to a temporary directory then we
collect some metadata, by getting the initial song request and querying
[musicbrainz](https://musicbrainz.org/) with it. After we collect the metadata,
the songs get added to a directory like ~/Music/yamp/Artist/Album/songname.mp3,
then ~/Music/playlists/playlistname.m3u gets written with the new song

## Contributing

Contributions are welcome! As long as:

- For big changes e.g implementing a new feature, please open an issue first.
  Same for bugs

- Smaller changes like typoes don't need an issue, just open a PR directly

- Getting help from AI in contributions is allowed, but don't fully automate the
  PR with tools like openclaw. Or use ai in a way that you can't reason over
  your changes, only open a PR if you know what changes you made and why

- Please make sure your code is formatted correctly, and actually runs as
  expected

- Adding libraries is typically fine, as long as it makes sense to add them,
  don't add another cli like urfave/cli when we already use cobra for example

## License

[MIT](https://choosealicense.com/licenses/mit/)
