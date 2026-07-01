<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import { ListPlaylists } from "../../../bindings/yamp/playlistrepository";
  import { onMount } from "svelte";

  let playlists: string[] = [];

  onMount(async () => {
    playlists = await ListPlaylists();
  });
</script>

<Sidebar.Root>
  <Sidebar.Content>
    <Sidebar.Group>
      <Sidebar.GroupLabel>Yamp</Sidebar.GroupLabel>
      <Sidebar.GroupContent>
        <Sidebar.Menu>
          {#each playlists as playlist (playlist)}
            <Sidebar.MenuItem>
              <Sidebar.MenuButton>
                {#snippet child({ props })}
                  <a href={"/playlists/" + playlist} {...props}>
                    <span>{playlist}</span>
                  </a>
                {/snippet}
              </Sidebar.MenuButton>
            </Sidebar.MenuItem>
          {/each}
        </Sidebar.Menu>
      </Sidebar.GroupContent>
    </Sidebar.Group>
  </Sidebar.Content>
</Sidebar.Root>
