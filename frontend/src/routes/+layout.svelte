<script lang="ts">
  let { children } = $props();
  import "../app.css";
  import { onMount } from "svelte";
  import { Theme } from "../../bindings/yamp/themeservice";

  import Navbar from "./components/Navbar.svelte";

  function disableContextMenu(e: MouseEvent) {
    e.preventDefault();
  }

  onMount(async () => {
    const theme = await Theme();
    document.documentElement.classList.remove("dark", "light");
    document.documentElement.classList.add(theme);
  });
</script>

<svelte:window on:contextmenu={disableContextMenu} />
<main class="flex h-screen font-inter bg-bg text-text overflow-hidden">
  <Navbar />
  <div class="flex-1 h-full overflow-y-auto">
    {@render children()}
  </div>
</main>
