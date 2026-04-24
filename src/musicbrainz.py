# pyright: reportArgumentType=false

from dataclasses import dataclass

import requests
from PySide6.QtCore import QObject, Slot
from typing_extensions import Optional

from music import download_song_to_tmp


@dataclass
class Recording:
    id: str
    title: str
    artist: str
    album: str
    duration_ms: Optional[int]
    release_id: str
    release_group_id: str
    release_date: str
    score: int


@dataclass
class Release:
    id: str
    title: str
    artist: str
    date: str
    release_id: str
    score: int


def parse_recording(data: dict) -> Recording:
    artist = "".join(
        ac.get("joinphrase", "") + ac["name"] for ac in data.get("artist-credit", [])
    )
    release = data.get("releases", [{}])[0]

    return Recording(
        id=data["id"],
        title=data["title"],
        artist=artist,
        album=release.get("title", ""),
        release_id=release.get("id", ""),
        release_group_id=release.get("release-group", {}).get("id", ""),
        duration_ms=data.get("length") or 0,
        release_date=data.get("first-release-date", ""),
        score=int(data.get("score") or 0),
    )


def search_recordings(query: str) -> list[Recording]:
    response = requests.get(
        "https://musicbrainz.org/ws/2/recording",
        params={"query": query, "fmt": "json", "limit": 50},
        headers={"User-Agent": "yamp/0.1 (Simon-Weij/yamp)"},
    )
    response.raise_for_status()
    data = response.json()
    return sorted(
        (parse_recording(r) for r in data["recordings"]),
        key=lambda recording: (
            recording.score,
            recording.release_date,
            recording.title,
        ),
        reverse=True,
    )


def parse_release(data: dict) -> Release:
    artist = "".join(
        ac.get("joinphrase", "") + ac["name"] for ac in data.get("artist-credit", [])
    )

    return Release(
        id=data["id"],
        title=data["title"],
        artist=artist,
        date=data.get("date", ""),
        release_id=data["id"],
        score=int(data.get("score") or 0),
    )


def search_releases(query: str) -> list[Release]:
    response = requests.get(
        "https://musicbrainz.org/ws/2/release",
        params={"query": query, "fmt": "json", "limit": 50},
        headers={"User-Agent": "yamp/0.1 (Simon-Weij/yamp)"},
    )
    response.raise_for_status()
    data = response.json()
    return sorted(
        (parse_release(r) for r in data["releases"]),
        key=lambda release: (release.score, release.date, release.title),
        reverse=True,
    )


def serialize_song(recording: Recording) -> dict:
    return {
        "recording_id": recording.id,
        "title": recording.title,
        "artist": recording.artist,
        "album": recording.album,
        "date": recording.release_date,
        "release_id": recording.release_id,
        "duration_ms": recording.duration_ms,
        "score": recording.score,
    }


def serialize_album(release: Release) -> dict:
    return {
        "title": release.title,
        "artist": release.artist,
        "album": release.title,
        "date": release.date,
        "release_id": release.release_id,
        "score": release.score,
    }


class Api(QObject):
    @Slot(str, result=list)
    def searchSongs(self, query: str):
        return [serialize_song(r) for r in search_recordings(query)]

    @Slot(str, result=list)
    def searchAlbums(self, query: str):
        return [serialize_album(r) for r in search_releases(query)]

    @Slot("QVariant")
    def playSong(self, data):
        download_song_to_tmp(data["artist"], data["title"])
        print(data)
