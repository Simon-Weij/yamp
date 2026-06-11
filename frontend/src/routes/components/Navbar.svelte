<script lang="ts">
  import { onMount } from "svelte";
  import { ListPlaylists } from "../../../bindings/yamp/playlistrepository";
  import { goto } from "$app/navigation";
  import { resolve } from "$app/paths";

  let playlists: string[] = [];

  onMount(async () => {
    playlists = await ListPlaylists();
  });
</script>

<nav class="w-60 bg-nav-bg text-nav-text h-full">
  <div class="ml-5 mt-5 mr-5">
    {#each playlists as playlist (playlist)}
      <button
        class="cursor-pointer hover:bg-button-nav-hover w-full px-4 py-2 text-left rounded"
        onclick={() => {
          goto(resolve("/playlists/[slug]", { slug: playlist }));
        }}
      >
        {playlist}
      </button>
    {/each}
  </div>
</nav>
