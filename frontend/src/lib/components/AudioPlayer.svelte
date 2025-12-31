<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { GetTrackFilePath } from '../../../wailsjs/go/main/App.js';
  import { player } from '../stores/player.svelte';

  let audioElement: HTMLAudioElement;
  let isLoading = $state(false);
  let isSeeking = $state(false);
  let lastSeekTime = $state(0);
  let lastLoadedTrackId = $state<string | null>(null);

  // Effect: Load track when currentTrack changes
  $effect(() => {
    if (player.currentTrack && audioElement && !isLoading && player.currentTrack.id !== lastLoadedTrackId) {
      loadTrack(player.currentTrack.id);
    }
  });

  // Effect: Play/Pause based on store state
  $effect(() => {
    if (audioElement && !isLoading) {
      if (player.isPlaying) {
        audioElement.play().catch(err => {
          console.error('Playback failed:', err);
          player.pause();
        });
      } else {
        audioElement.pause();
      }
    }
  });

  // Effect: Volume control
  $effect(() => {
    if (audioElement) {
      audioElement.volume = player.isMuted ? 0 : player.volume;
    }
  });

  // Effect: Seek control (when user drags seekbar)
  $effect(() => {
    if (audioElement && !isSeeking && Math.abs(player.currentTime - audioElement.currentTime) > 1) {
      audioElement.currentTime = player.currentTime;
    }
  });

  async function loadTrack(trackId: string) {
    if (!trackId) return;

    // Prevent multiple simultaneous loads
    if (isLoading) return;

    // Set both guards IMMEDIATELY before any async operations
    lastLoadedTrackId = trackId;
    isLoading = true;

    try {
      const filePath = await GetTrackFilePath(trackId);

      // Use HTTP endpoint to stream the file (avoids file:// CORS issues)
      const audioUrl = `/audio/stream?path=${encodeURIComponent(filePath)}`;
      audioElement.src = audioUrl;

      // Auto-play if isPlaying is true
      if (player.isPlaying) {
        await audioElement.play();
      }
    } catch (err) {
      console.error('Failed to load track:', err);
      player.pause();
      // Reset lastLoadedTrackId on error so user can retry
      lastLoadedTrackId = null;
    } finally {
      isLoading = false;
    }
  }

  function handleTimeUpdate() {
    if (audioElement && !isNaN(audioElement.currentTime)) {
      player.updateTime(audioElement.currentTime);
    }
  }

  function handleEnded() {
    player.handleTrackEnd();
  }

  function handleError(e: Event) {
    const audio = e.target as HTMLAudioElement;
    const error = audio.error;

    console.error('Audio playback error:', {
      code: error?.code,
      message: error?.message,
      src: audio.src,
      networkState: audio.networkState,
      readyState: audio.readyState
    });

    // Error codes:
    // 1 = MEDIA_ERR_ABORTED
    // 2 = MEDIA_ERR_NETWORK
    // 3 = MEDIA_ERR_DECODE
    // 4 = MEDIA_ERR_SRC_NOT_SUPPORTED

    if (error?.code === 4) {
      console.error('File format not supported or file not accessible');
    } else if (error?.code === 2) {
      console.error('Network error - file might not exist or is not accessible');
    }

    player.pause();
  }

  function handleLoadedMetadata() {
    updateMediaSessionMetadata();
  }

  function updateMediaSessionMetadata() {
    if ('mediaSession' in navigator && player.currentTrack) {
      const metadata: MediaMetadataInit = {
        title: player.currentTrack.title || 'Unknown Title',
        artist: player.currentTrack.artist || 'Unknown Artist',
        album: player.currentTrack.album || 'Unknown Album',
      };

      // Add artwork if available
      if (player.currentTrack.artworkPath) {
        const artworkUrl = `/artwork/stream?file=${encodeURIComponent(player.currentTrack.artworkPath)}`;
        metadata.artwork = [
          { src: artworkUrl, sizes: '512x512', type: 'image/jpeg' },
        ];
      }

      navigator.mediaSession.metadata = new MediaMetadata(metadata);
    }
  }

  function handlePlay() {
    // Audio started playing (could be triggered by OS/browser)
    // Only update store if it's out of sync to avoid infinite loops
    if (!player.isPlaying) {
      player.resume();
    }
  }

  function handlePause() {
    // Audio was paused (could be triggered by OS/browser)
    // Only update store if it's out of sync to avoid infinite loops
    if (player.isPlaying) {
      player.pause();
    }
  }

  onMount(() => {
    // Initialize audio element properties
    if (audioElement) {
      audioElement.volume = player.volume;
      audioElement.preload = 'metadata';
    }

    if ('mediaSession' in navigator) {
      navigator.mediaSession.setActionHandler('play', () => {
        player.resume();
      });

      navigator.mediaSession.setActionHandler('pause', () => {
        player.pause();
      });

      navigator.mediaSession.setActionHandler('previoustrack', () => {
        player.previous();
      });

      navigator.mediaSession.setActionHandler('nexttrack', () => {
        player.next();
      });

      navigator.mediaSession.setActionHandler('seekbackward', (details) => {
        const skipTime = details.seekOffset || 10;
        if (audioElement) {
          audioElement.currentTime = Math.max(0, audioElement.currentTime - skipTime);
        }
      });

      navigator.mediaSession.setActionHandler('seekforward', (details) => {
        const skipTime = details.seekOffset || 10;
        if (audioElement) {
          audioElement.currentTime = Math.min(audioElement.duration, audioElement.currentTime + skipTime);
        }
      });
    }
  });

  onDestroy(() => {
    // Cleanup: pause and clear source
    if (audioElement) {
      audioElement.pause();
      audioElement.src = '';
    }
  });
</script>

<!-- Hidden audio element - managed programmatically -->
<audio
  bind:this={audioElement}
  ontimeupdate={handleTimeUpdate}
  onended={handleEnded}
  onerror={handleError}
  onloadedmetadata={handleLoadedMetadata}
  onplay={handlePlay}
  onpause={handlePause}
  style="display: none;"
></audio>