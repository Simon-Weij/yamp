<script lang="ts">
  let { children } = $props();
  import "../app.css";

  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import AppSidebar from "$lib/components/app-sidebar.svelte";
  function disableContextMenu(e: MouseEvent) {
    e.preventDefault();
  }

  import { onMount } from "svelte";

  onMount(() => {
    document.documentElement.classList.add("dark");
  });
</script>

<svelte:window on:contextmenu={disableContextMenu} />
<Sidebar.Provider>
  <AppSidebar />
  <main>
    {@render children?.()}
  </main>
</Sidebar.Provider>>
