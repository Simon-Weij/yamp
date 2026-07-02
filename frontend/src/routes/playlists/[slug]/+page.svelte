<script lang="ts">
  import { page } from "$app/state";
  import { PlaylistItem, Song } from "../../../../bindings/yamp/models";
  import {
    ListSongsInPlaylist,
    AddSongToPlaylist,
  } from "../../../../bindings/yamp/playlistrepository";
  import { SearchSong } from "../../../../bindings/yamp/browserrepository";
  import * as Table from "$lib/components/ui/table/index.js";
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import Button from "$lib/components/ui/button/button.svelte";
  import { Plus } from "@lucide/svelte";
  import Input from "$lib/components/ui/input/input.svelte";

  let searchSongs: Song[] = $state([]);
  let open = $state(false);
  let searchSongPrompt: string = $state("");

  let playlist = $derived(page.params.slug);
  let songs = $state<PlaylistItem[]>([]);

  function msToMinutesSeconds(ms: number): string {
    const totalSeconds = Math.floor(ms / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;

    const paddedSeconds = seconds.toString().padStart(2, "0");

    return `${minutes}:${paddedSeconds}`;
  }

  async function addSongToPlaylist(song: Song) {
    await AddSongToPlaylist(song, playlist ?? "");
    songs = await ListSongsInPlaylist(playlist ?? "");
    open = false;
  }

  $effect(() => {
    if (!searchSongPrompt) return;

    const id = setTimeout(async () => {
      searchSongs = await SearchSong(searchSongPrompt);
    }, 300);

    return () => clearTimeout(id);
  });

  $effect(() => {
    ListSongsInPlaylist(playlist ?? "").then((res) => {
      songs = res;
    });
  });
</script>

<Dialog.Root bind:open>
  <main class="w-full flex flex-col justify-center">
    <h1 class="font-bold text-3xl text-center mt-5">{playlist}</h1>
    <Dialog.Trigger>
      <Button class="mx-auto my-5 cursor-pointer"><Plus /> Add Song</Button>
    </Dialog.Trigger>
    <Table.Root class="w-full">
      <Table.Header>
        <Table.Row>
          <Table.Head class="w-15 text-xs text-gray-500">#</Table.Head>
          <Table.Head class="w-15 text-xs text-gray-500">cover</Table.Head>
          <Table.Head class="text-xs text-gray-500">Song</Table.Head>
          <Table.Head class="w-15 text-xs text-gray-500">Duration</Table.Head>
          <Table.Head />
          <Table.Head />
        </Table.Row>
      </Table.Header>
      <Table.Body class="cursor-pointer">
        {#each songs as song, i (song)}
          <Table.Row>
            <Table.Cell>{i + 1}</Table.Cell>
            <Table.Cell>
              <img src={song.Cover} alt={song.Album} class="rounded" />
            </Table.Cell>
            <Table.Cell>
              <div class="flex flex-col">
                {song.Title}
                <div>
                  <span class="text-gray-400">{song.Artist} - {song.Album}</span
                  >
                </div>
              </div></Table.Cell
            >
            <Table.Cell>{msToMinutesSeconds(song.Duration)}</Table.Cell>
            <Table.Cell />
            <Table.Cell />
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
  </main>

  <Dialog.Content class="sm:max-w-2xl max-h-[80vh] flex flex-col">
    <Dialog.Header>
      <Dialog.Title>Add song</Dialog.Title>

      <Input placeholder="Search..." bind:value={searchSongPrompt} />

      <div class="max-h-96 overflow-y-auto">
        <Table.Root class="w-full">
          <Table.Header class="sticky top-0">
            <Table.Row>
              <Table.Head class="w-15 text-xs text-gray-500">Cover</Table.Head>
              <Table.Head class="text-xs text-gray-500">Song</Table.Head>
              <Table.Head class="w-15 text-xs text-gray-500"
                >Duration</Table.Head
              >
              <Table.Head />
              <Table.Head />
            </Table.Row>
          </Table.Header>

          <Table.Body class="cursor-pointer">
            {#each searchSongs as song (song)}
              <Table.Row
                onclick={() => {
                  addSongToPlaylist(song);
                }}
              >
                <Table.Cell>
                  <img src={song.artworkUrl100} alt={song.collectionName} />
                </Table.Cell>

                <Table.Cell>
                  <div class="flex flex-col">
                    {song.trackName}
                    <div>
                      <span class="text-gray-400">
                        {song.artistName} - {song.collectionName}
                      </span>
                    </div>
                  </div>
                </Table.Cell>

                <Table.Cell>
                  {msToMinutesSeconds(song.trackTimeMillis)}
                </Table.Cell>

                <Table.Cell />
                <Table.Cell />
              </Table.Row>
            {/each}
          </Table.Body>
        </Table.Root>
      </div>
    </Dialog.Header>
  </Dialog.Content>
</Dialog.Root>
