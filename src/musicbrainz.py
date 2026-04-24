from dataclasses import asdict, dataclass

import certifi
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


class Api(QObject):
    @Slot(str, result=list)
    def searchRecordings(self, query: str):
        return [asdict(r) for r in search_recordings(query)]
