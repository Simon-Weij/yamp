<script lang="ts">
  import { page } from "$app/state";
  //import { onMount } from "svelte";
  import { PlaylistItem } from "../../../../bindings/yamp/models";
  //import { ListSongsInPlaylist } from "../../../../bindings/yamp/playlistrepository";
  import * as Table from "$lib/components/ui/table/index.js";

  export function msToMinutesSeconds(ms: number): string {
    const totalSeconds = Math.floor(ms / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;

    const paddedSeconds = seconds.toString().padStart(2, "0");

    return `${minutes}:${paddedSeconds}`;
  }

  let playlist = $derived(page.params.slug);
  let songs = $state<PlaylistItem[]>([
    {
      Artist: "Artist 1",
      Album: "Album 1",
      Title: "Song 1",
      Cover: "https://placehold.co/400x400",
      Duration: 240000,
    },
    {
      Artist: "Artist 2",
      Album: "Album 2",
      Title: "Song 2",
      Cover: "https://placehold.co/400x400",
      Duration: 180000,
    },
    {
      Artist: "Artist 3",
      Album: "Album 1",
      Title: "Song 3",
      Cover: "https://placehold.co/400x400",
      Duration: 210000,
    },
    {
      Artist: "Artist 1",
      Album: "Album 3",
      Title: "Song 4",
      Cover: "https://placehold.co/400x400",
      Duration: 195000,
    },
    {
      Artist: "Artist 4",
      Album: "Album 4",
      Title: "Song 5",
      Cover: "https://placehold.co/400x400",
      Duration: 225000,
    },
  ]);
  /*  onMount(async () => {
    songs = await ListSongsInPlaylist(playlist ?? "");
  });
  */
</script>

<main class="w-full">
  <h1 class="font-bold text-3xl text-center">{playlist}</h1>
  <Table.Root class="w-full">
    <Table.Header>
      <Table.Row>
        <Table.Head class="w-15 text-xs text-gray-500">#</Table.Head>
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
            <div class="flex flex-col">
              {song.Title}
              <div>
                <span class="text-gray-400">{song.Artist} - {song.Album}</span>
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
