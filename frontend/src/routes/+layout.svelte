<script lang="ts">
  let { children } = $props();
  import "../app.css";
  import { onMount } from "svelte";
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import AppSidebarLeft from "$lib/components/app-sidebar-left.svelte";
  import AppSidebarRight from "$lib/components/app-sidebar-right.svelte";

  function disableContextMenu(e: MouseEvent) {
    e.preventDefault();
  }

  onMount(() => {
    document.documentElement.classList.add("dark");
  });
</script>

<svelte:window on:contextmenu={disableContextMenu} />
<Sidebar.Provider>
  <AppSidebarLeft />
  <main class="w-screen">
    {@render children?.()}
  </main>
  <AppSidebarRight />
</Sidebar.Provider>>
