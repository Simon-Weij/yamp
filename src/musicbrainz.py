# pyright: reportArgumentType=false

from dataclasses import dataclass

import requests
from PySide6.QtCore import QObject, Slot
from typing_extensions import Optional


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


@dataclass
class Release:
    id: str
    title: str
    artist: str
    date: str
    release_id: str


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
    )


def search_recordings(query: str) -> list[Recording]:
    response = requests.get(
        "https://musicbrainz.org/ws/2/recording",
        params={"query": query, "fmt": "json"},
        headers={"User-Agent": "yamp/0.1 (Simon-Weij/yamp)"},
    )
    response.raise_for_status()
    data = response.json()
    return [parse_recording(r) for r in data["recordings"]]


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
    )


def search_releases(query: str) -> list[Release]:
    response = requests.get(
        "https://musicbrainz.org/ws/2/release",
        params={"query": query, "fmt": "json"},
        headers={"User-Agent": "yamp/0.1 (Simon-Weij/yamp)"},
    )
    response.raise_for_status()
    data = response.json()
    return [parse_release(r) for r in data["releases"]]


def serialize_song(recording: Recording) -> dict:
    return {
        "recording_id": recording.id,
        "title": recording.title,
        "artist": recording.artist,
        "album": recording.album,
        "date": recording.release_date,
        "release_id": recording.release_id,
        "duration_ms": recording.duration_ms,
    }


def serialize_album(release: Release) -> dict:
    return {
        "title": release.title,
        "artist": release.artist,
        "album": release.title,
        "date": release.date,
        "release_id": release.release_id,
    }


def dedupe_results_by_artist(results: list[dict]) -> list[dict]:
    seen_artists = set()
    unique_results = []

    for result in results:
        artist = (result.get("artist") or "").strip().casefold()
        if artist in seen_artists:
            continue

        seen_artists.add(artist)
        unique_results.append(result)

    return unique_results


class Api(QObject):
    @Slot(str, result=list)
    def searchSongs(self, query: str):
        return [serialize_song(r) for r in search_recordings(query)]

    @Slot(str, result=list)
    def searchAlbums(self, query: str):
        return dedupe_results_by_artist(
            [serialize_album(r) for r in search_releases(query)]
        )

    @Slot("QVariant")
    def playSong(self, data):
        print(data)
