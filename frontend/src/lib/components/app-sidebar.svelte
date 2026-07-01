<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { ListPlaylists } from "../../../bindings/yamp/playlistrepository";
  import { Button } from "$lib/components/ui/button/index.js";
  import { CreatePlaylist } from "../../../bindings/yamp/playlistrepository";

  let playlistName = $state("");
  let open = $state(false);

  import { onMount } from "svelte";
  import { Plus } from "@lucide/svelte";
  import Input from "./ui/input/input.svelte";

  let playlists: string[] = $state([]);

  async function createPlaylist() {
    console.log("making playlist " + playlistName);
    if (!playlistName.trim()) return;
    await CreatePlaylist(playlistName);
    // Sort alphabetically
    playlists = [...playlists, playlistName.trim()].sort((a, b) =>
      a.localeCompare(b),
    );
    playlistName = "";

    open = false;
  }

  onMount(async () => {
    playlists = await ListPlaylists();
  });
</script>

<Dialog.Root bind:open>
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

            <Sidebar.MenuItem>
              <Dialog.Trigger>
                {#snippet child({ props })}
                  <Sidebar.MenuButton class="cursor-pointer" {...props}>
                    <Plus /> Create playlist
                  </Sidebar.MenuButton>
                {/snippet}
              </Dialog.Trigger>
            </Sidebar.MenuItem>
          </Sidebar.Menu>
        </Sidebar.GroupContent>
      </Sidebar.Group>
    </Sidebar.Content>

    <Dialog.Content>
      <Dialog.Header>
        <Dialog.Title>Create playlist</Dialog.Title>
      </Dialog.Header>

      <Input placeholder="playlist" bind:value={playlistName} />

      <Button
        variant="outline"
        class="cursor-pointer"
        onclick={() => createPlaylist()}>Create</Button
      >
    </Dialog.Content>
  </Sidebar.Root>
</Dialog.Root>
