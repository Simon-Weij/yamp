<script lang="ts">
  import { onMount } from "svelte";
  import {
    ListPlaylists,
    CreatePlaylist,
  } from "../../../bindings/yamp/playlistrepository";
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";
  import { Plus } from "@lucide/svelte";

  let playlists: string[] = [];
  let showInput = false;
  let playlistName = "";
  let inputEl: HTMLInputElement;

  onMount(async () => {
    playlists = await ListPlaylists();
  });

  function openInput() {
    playlistName = "";
    showInput = true;
    setTimeout(() => inputEl?.focus(), 0);
  }

  function closeInput() {
    showInput = false;
    playlistName = "";
  }

  function handleCreate() {
    if (!playlistName.trim()) return;
    CreatePlaylist(playlistName);
    playlists = [...playlists, playlistName.trim()];
    closeInput();
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === "Escape") closeInput();
    if (e.key === "Enter") handleCreate();
  }
</script>

<nav class="w-60 bg-nav-bg text-nav-text h-full">
  <div class="ml-5 mt-5 mr-5">
    {#each playlists as playlist (playlist)}
      <button
        class="cursor-pointer hover:bg-button-nav-hover w-full px-4 py-2 text-left rounded"
        onclick={() => goto(resolve("/playlists/[slug]", { slug: playlist }))}
      >
        {playlist}
      </button>
    {/each}

    <button
      class="cursor-pointer flex gap-1 hover:bg-button-nav-hover w-full px-4 py-2 text-left rounded"
      onclick={openInput}
    >
      <Plus size={18} /> Create playlist
    </button>

    {#if showInput}
      <div class="mt-1 px-2 py-2 rounded bg-white/5 flex flex-col gap-2">
        <input
          bind:this={inputEl}
          bind:value={playlistName}
          type="text"
          placeholder="Playlist name"
          onkeydown={handleKeydown}
          class="bg-transparent border border-white/20 rounded px-2 py-1 text-sm outline-none focus:border-white/50 placeholder:text-nav-text/40 w-full"
        />
        <div class="flex gap-1 justify-end">
          <button
            onclick={closeInput}
            class="text-xs px-2 py-1 rounded hover:bg-button-nav-hover"
          >
            Cancel
          </button>
          <button
            onclick={handleCreate}
            disabled={!playlistName.trim()}
            class="text-xs px-2 py-1 rounded bg-white/10 hover:bg-white/20 disabled:opacity-40 disabled:cursor-not-allowed"
          >
            Create
          </button>
        </div>
      </div>
    {/if}
  </div>
</nav>
