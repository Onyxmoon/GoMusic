<script lang="ts">
  import { player } from '../stores/player.svelte';
  import { extractColorsFromImage, applyGradient, resetGradient } from '../utils/colorExtractor';

  $effect(() => {
    if (player.currentTrack?.artworkPath) {
      const artworkUrl = `/artwork/stream?file=${encodeURIComponent(player.currentTrack.artworkPath)}`;
      extractColorsFromImage(artworkUrl)
        .then(colors => applyGradient(colors))
        .catch(err => {
          console.error('Failed to extract colors:', err);
          resetGradient();
        });
    } else {
      resetGradient();
    }
  });
</script>