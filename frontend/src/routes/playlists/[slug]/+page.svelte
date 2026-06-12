<script lang="ts">
  import { page } from "$app/state";
  import { PlaylistItem } from "../../../../bindings/yamp/models";
  import { ParsePlaylistFile } from "../../../../bindings/yamp/playlistrepository";

  let slug = $derived(page.params.slug);

  let playlists: PlaylistItem[] = $state([]);

  $effect(() => {
    const fetchPlaylists = async () => {
      const data = await ParsePlaylistFile(slug ?? "");
      playlists = data || [];
    };

    fetchPlaylists();
  });
</script>

<h1>Playlist: {slug}</h1>

{#each playlists as item (item.Title)}
  <p>Title: {item.Title}</p>
  <p>Album: {item.Album}</p>
  <p>Artist: {item.Artist}</p>
{/each}
