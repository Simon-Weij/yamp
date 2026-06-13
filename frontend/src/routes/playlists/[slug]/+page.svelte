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

<main class="flex flex-col min-h-screen p-6">
  <div class="w-full flex flex-col gap-3">
    {#each playlists as item, i (item.Title)}
      <section class="cursor-pointer rounded-lg p-4 hover:bg-button-nav-hover">
        <div class="font-semibold">
          {i + 1}. {item.Title}
        </div>

        <div class="text-sm opacity-80 mt-1">
          {item.Album} - {item.Artist}
        </div>
      </section>
    {/each}
  </div>
</main>
