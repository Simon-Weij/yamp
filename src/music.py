import subprocess
import time


def download_song_to_tmp(artist: str, song: str):
    timestamp = int(time.time())
    output_path = f"/tmp/yamp/{artist}-{song}-{timestamp}.%(ext)s"

    subprocess.run(
        [
            "yt-dlp",
            f"ytsearch1:{artist} - {song}",
            "-x",
            "--audio-format",
            "mp3",
            "-o",
            output_path,
        ]
    )
    return output_path
