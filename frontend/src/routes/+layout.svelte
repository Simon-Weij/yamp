<script lang="ts">
  let { children } = $props();
  import "../app.css";
  import { onMount } from "svelte";

  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import AppSidebar from "$lib/components/app-sidebar.svelte";
  function disableContextMenu(e: MouseEvent) {
    e.preventDefault();
  }

  onMount(() => {
    document.documentElement.classList.add("dark");
  });
</script>

<svelte:window on:contextmenu={disableContextMenu} />
<Sidebar.Provider>
  <AppSidebar />
  <main class="w-screen">
    {@render children?.()}
  </main>
</Sidebar.Provider>>
