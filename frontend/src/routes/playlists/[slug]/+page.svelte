<script lang="ts">
  import { page } from "$app/state";
  import { PlaylistItem } from "../../../../bindings/yamp/models";
  import { ParsePlaylistFile } from "../../../../bindings/yamp/playlistrepository";
  import { GetAlbumCoverBase64 } from "../../../../bindings/yamp/musicservice";
  import { Plus } from "@lucide/svelte";
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
      <section class="cursor-pointer rounded-lg p-4 hover:bg-button-nav-hover flex items-center gap-4">
        {#await GetAlbumCoverBase64(item.Artist, item.Album) then src}
          <img {src} alt="cover" class="w-16 h-16 object-cover rounded-md shrink-0" />
        {/await}
        <div class="flex flex-col min-w-0">
          <div class="font-semibold truncate">{i + 1}. {item.Title}</div>
          <div class="text-sm opacity-80 mt-1 truncate">{item.Album} - {item.Artist}</div>
        </div>
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
