<script lang="ts">
  import {
    SearchSong,
    FetchImageAsBase64,
    AddSongToPlaylist,
  } from "../../../bindings/yamp/browserrepository";
  import { Song } from "../../../bindings/yamp/models";
  export let onclose: () => void;
  export let onsongadded: () => void = () => {};

  let timeout: ReturnType<typeof setTimeout>;
  let songs: Song[] = [];

  export let playlist: string;

  function handleInput(e: Event) {
    clearTimeout(timeout);
    timeout = setTimeout(async () => {
      const query = (e.target as HTMLInputElement).value;
      songs = await SearchSong(query);
    }, 300);
  }

  async function addSongToPlaylist(song: Song) {
    console.log(song);
    await AddSongToPlaylist(song, playlist);
    onsongadded();
    onclose();
  }

  function formatDuration(ms: number): string {
    if (!ms) return "--:--";
    const totalSeconds = Math.floor(ms / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
  }
</script>

<div
  class="rounded-xl shadow-xl bg-modal-bg p-6 flex flex-col w-md h-128"
  role="presentation"
  onclick={(e) => e.stopPropagation()}
>
  <input
    type="text"
    placeholder="Search songs..."
    class="bg-bg px-4 py-3 shrink-0"
    oninput={handleInput}
  />
  <ul class="mt-4 overflow-y-auto flex-1 flex flex-col gap-2">
    {#each songs as song (song)}
      <button
        type="button"
        class="flex flex-row items-center gap-3 hover:bg-bg py-2 px-3 rounded text-left w-full"
        onclick={() => addSongToPlaylist(song)}
      >
        {#await FetchImageAsBase64(song.artworkUrl100) then src}
          <img
            {src}
            alt={song.collectionName}
            class="w-10 h-10 rounded shrink-0"
          />
        {/await}
        <div class="flex flex-col flex-1 min-w-0">
          <span class="font-medium truncate">{song.trackName}</span>
          <span class="text-sm text-gray-400 truncate"
            >{song.artistName} - {song.collectionName}</span
          >
        </div>
        <span class="text-sm text-gray-400 shrink-0 tabular-nums">
          {formatDuration(song.trackTimeMillis)}
        </span>
      </button>
    {/each}
  </ul>
</div>
