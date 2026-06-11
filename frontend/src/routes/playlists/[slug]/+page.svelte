<script>
  import { page } from "$app/state";

  let slug = $derived(page.params.slug);
</script>

<h1>Hello world</h1>

<p>Current slug: {slug}</p>
