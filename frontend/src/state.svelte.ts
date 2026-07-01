import { PlaylistItem } from "../bindings/yamp/models"

export const currentSong = $state<PlaylistItem>({
  Artist: "Unknown",
  Album: "Unknown",
  Title: "Unknown",
  Cover: "https://placehold.co/400x400",
  Duration: "00:00"
})
