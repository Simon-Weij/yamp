<script>
  let { children } = $props();
  import "../app.css";
  import { onMount } from "svelte";
  import { Theme } from "../../bindings/yamp/themeservice";

  import Navbar from "./components/Navbar.svelte";

  onMount(async () => {
    const theme = await Theme();
    document.documentElement.classList.remove("dark", "light");
    document.documentElement.classList.add(theme);
  });
</script>

<main class="font-inter h-screen flex bg-bg text-text">
  <Navbar />
  <div class="flex-1">
    {@render children()}
  </div>
</main>
