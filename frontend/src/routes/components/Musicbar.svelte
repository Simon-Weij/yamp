<script>
  import { currentSong } from "../../state.svelte";
  import { ChevronFirst, ChevronLast, Pause } from "@lucide/svelte";
  let progress = $state(45);
  let currentTime = $state("11:11");
</script>

<main class="bg-musicbar flex items-center px-3 py-2">
  <img
    src={currentSong.Cover}
    alt={currentSong.Album}
    class="w-16 h-16 rounded object-cover shrink-0 m-2"
  />

  <div class="ml-3">
    <div class="font-bold text-sm">{currentSong.Title}</div>
    <div class="text-xs text-gray-600">{currentSong.Artist}</div>
  </div>

  <div class="flex-1 flex flex-col items-center">
    <div class="flex gap-4 items-center mb-2">
      <ChevronFirst
        size={36}
        class="cursor-pointer text-gray-300 hover:text-white"
      />
      <Pause size={32} class="cursor-pointer text-gray-300 hover:text-white" />
      <ChevronLast
        size={36}
        class="cursor-pointer text-gray-300 hover:text-white"
      />
    </div>

    <div class="flex items-center gap-2 w-96 max-w-full">
      <span class="text-xs text-gray-400 w-10 text-right">
        {currentTime}
      </span>

      <div class="flex-1 bg-gray-700 rounded-full h-1">
        <div
          class="bg-white h-1 rounded-full transition-all"
          style={`width: ${progress}%`}
        ></div>
      </div>

      <span class="text-xs text-gray-400 w-10">
        {currentSong.Duration}
      </span>
    </div>
  </div>
</main>
