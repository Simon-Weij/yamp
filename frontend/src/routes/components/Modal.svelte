<script lang="ts">
  import {
    SearchSong,
    FetchImageAsBase64,
  } from "../../../bindings/yamp/browserrepository";
  import { Song } from "../../../bindings/yamp/models";

  let timeout: ReturnType<typeof setTimeout>;
  let songs: Song[] = [];

  function handleInput(e: Event) {
    clearTimeout(timeout);
    timeout = setTimeout(async () => {
      const query = (e.target as HTMLInputElement).value;
      songs = await SearchSong(query);
    }, 300);
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
      <li class="flex flex-row items-center gap-3">
        {#await FetchImageAsBase64(song.artworkUrl100) then src}
          <img
            {src}
            alt={song.collectionName}
            class="w-10 h-10 rounded shrink-0"
          />
        {/await}
        <div class="flex flex-col">
          <span class="font-medium">{song.trackName}</span>
          <span class="text-sm text-gray-400"
            >{song.artistName} - {song.collectionName}</span
          >
        </div>
      </li>
    {/each}
  </ul>
</div>
