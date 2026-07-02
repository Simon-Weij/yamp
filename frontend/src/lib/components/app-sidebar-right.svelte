<script lang="ts">
  import * as Sidebar from "$lib/components/ui/sidebar/index.js";
  import { onMount } from "svelte";
  import { CurrentlyPlaying } from "../../../bindings/yamp/songrepository";
  import { PlaylistItem } from "../../../bindings/yamp/models";
  let currentSong = $state<PlaylistItem>({
    Artist: "",
    Album: "",
    Title: "",
    Cover: "",
    Duration: 0,
  });
  let currentTime = $state(10000);
  let totalTime = $state(200000);

  let progressPercent = $derived(
    totalTime > 0 ? (currentTime / totalTime) * 100 : 0,
  );

  function msToMinutesSeconds(ms: number): string {
    const totalSeconds = Math.floor(ms / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    const paddedSeconds = seconds.toString().padStart(2, "0");
    return `${minutes}:${paddedSeconds}`;
  }
  onMount(async () => {
    currentSong = await CurrentlyPlaying();
  });
</script>

<Sidebar.Root side="right">
  <Sidebar.Content>
    <div class="m-5">
      <img
        src={currentSong.Cover}
        class="rounded-2xl"
        alt={currentSong.Album}
      />
    </div>
    <h1 class="text-center font-bold text-xl">{currentSong.Title}</h1>
    <p class="text-center text-gray-300">
      {currentSong.Artist} - {currentSong.Album}
    </p>

    <div class="mt-3 px-2 mx-7">
      <div class="relative h-1 w-full rounded-full bg-neutral-700">
        <div
          class="absolute left-0 top-0 h-full rounded-full bg-white"
          style="width: {progressPercent}%"
        ></div>
      </div>
    </div>

    <div class="mt-1 flex justify-between text-xs text-neutral-400 mx-7">
      <span>{msToMinutesSeconds(currentTime)}</span>
      <span>{msToMinutesSeconds(totalTime)}</span>
    </div>
  </Sidebar.Content>
</Sidebar.Root>
