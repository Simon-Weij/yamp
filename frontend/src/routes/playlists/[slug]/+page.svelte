<script>
  import { page } from "$app/state";
  let playlist = $derived(page.params.slug);
</script>

<p>{playlist}</p>
