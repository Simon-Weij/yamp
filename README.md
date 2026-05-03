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

Currently, we don't package our project yet, since the project is in very early
development, if you want to install the project anyways you can install it with
the following steps

```bash
# Clone the project 
git clone https://github.com/Simon-Weij/yamp.git

# Build the project
go build -o yamp

# Move the binary to a location in your $PATH
mv yamp /some/other/location/in/path
```

## Roadmap

- More flexible playlists, currently they're really limited

- Export playlists from popular platforms e.g Spotify/yt-music

- Write tests

## How does it work?

#### The play command

Currently this project heavily depends on yt-dlp for its core functionality,
when a song is played using for example yamp play "Artist - Song name" yt-dlp is
asked to download the song to a temporary directory, and is then played with
mpv.

#### Playlists

For playlists the song is first also downloaded to a temporary directory, very
similar to the play command, but we collect some metadata, by getting the
initial song request (Artist - Song name) and querying
[musicbrainz](https://musicbrainz.org/) with it. After we collect the metadata,
the songs get added to ~/Music/yamp/Artist/Album/songname.mp3, then this
metadata and song gets added to the specified playlist

## Contributing

Contributions are welcome! As long as:

- For big changes e.g implementing a new feature, please open an issue first.
  Same for bugs

- Smaller changes like typoes don't need an issue, just open a PR directly

- Getting help from AI in contributions is allowed, except if the PR is fully
  automated with for example openclaw. Or if you can't reason over your changes
  or don't understand them, only make a PR if you know what you changed and why

- Please make sure your code is formatted correctly, and actually runs as
  expected

- Adding libraries is typically fine, as long as it makes sense to add them,
  don't add another cli like urfave/cli when we already use cobra for example

## License

[MIT](https://choosealicense.com/licenses/mit/)
