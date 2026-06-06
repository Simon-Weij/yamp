<script>
  let { children } = $props();
  import "../app.css";
  import { onMount } from "svelte";

  import { Theme } from "../../bindings/yamp/themeservice";

  onMount(async () => {
    const theme = await Theme();
    document.documentElement.classList.remove("dark", "light");

    document.documentElement.classList.add(theme);
  });
</script>

<main
  class="font-inter light:bg-[#EDEDED] light:text-[#2B2B2B] h-screen dark:bg-black"
>
  {@render children()}
</main>
