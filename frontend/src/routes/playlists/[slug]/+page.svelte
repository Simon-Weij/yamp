<script lang="ts">
  import { page } from "$app/state";
  import { PlaylistItem } from "../../../../bindings/yamp/models";
  import { ParsePlaylistFile } from "../../../../bindings/yamp/playlistrepository";
  import { RemoveSongFromPlaylist } from "../../../../bindings/yamp/browserrepository";
  import { GetAlbumCoverPath } from "../../../../bindings/yamp/musicservice";
  import { Plus, Trash2 } from "@lucide/svelte";
  import Modal from "../../components/Modal.svelte";
  let open = $state(false);

  let slug = $derived(page.params.slug);

  let playlists: PlaylistItem[] = $state([]);

  async function fetchPlaylists() {
    const currentSlug = slug;
    const data = await ParsePlaylistFile(currentSlug ?? "");
    playlists = data || [];
  }

  $effect(() => {
    fetchPlaylists();
  });
</script>

<main class="flex flex-col min-h-screen p-6">
  <div class="w-full flex flex-col gap-3">
    {#each playlists as item, i (item.Title)}
      <section
        class="cursor-pointer rounded-lg p-4 hover:bg-button-nav-hover flex items-center gap-4"
      >
        {#await GetAlbumCoverPath(item.Artist, item.Album) then path}
          <img
            src={path}
            alt="cover"
            class="w-16 h-16 object-cover rounded-md shrink-0"
            loading="lazy"
          />
        {/await}
        <div class="flex flex-col min-w-0 flex-1">
          <div class="font-semibold truncate">{i + 1}. {item.Title}</div>
          <div class="text-sm opacity-80 mt-1 truncate">
            {item.Album} - {item.Artist}
          </div>
        </div>
        <button
          class="shrink-0 p-1.5 rounded-md opacity-40 hover:opacity-100 hover:text-red-500 hover:bg-red-500/10 transition-all cursor-pointer"
          onclick={async (e) => {
            e.stopPropagation();
            await RemoveSongFromPlaylist(
              item.Title,
              item.Album,
              item.Artist,
              slug ?? "",
            );
            await fetchPlaylists();
          }}
          aria-label="Remove"
        >
          <Trash2 size={16} />
        </button>
      </section>
    {/each}
    <button
      class="cursor-pointer flex flex-row items-center rounded-lg p-4 hover:bg-button-nav-hover font-bold"
      onclick={() => (open = true)}
    >
      <Plus /> Add song
    </button>
  </div>
  {#if open}
    <div
      class="fixed inset-0 z-40 bg-modal-overlay"
      role="presentation"
      onclick={() => (open = false)}
    >
      <div
        class="fixed inset-0 z-50 flex items-center justify-center"
        role="dialog"
        aria-modal="true"
      >
        <Modal
          playlist={slug ?? ""}
          onclose={() => (open = false)}
          onsongadded={fetchPlaylists}
        />
      </div>
    </div>
  {/if}
</main>
