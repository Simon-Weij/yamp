<script lang="ts">
  import { page } from "$app/state";
  import { PlaylistItem } from "../../../../bindings/yamp/models";
  import { ParsePlaylistFile } from "../../../../bindings/yamp/playlistrepository";
  import { Plus } from "@lucide/svelte";
  import Modal from "../../components/Modal.svelte";
  let open = $state(false);

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
        <Modal />
      </div>
    </div>
  {/if}
</main>
