<script lang="ts">
  import { player } from '../stores/player.svelte';
  import { Pause } from 'lucide-svelte';

  interface Props {
    trackId: string;
    size?: number;
    showWhenPaused?: boolean;
    className?: string;
  }

  let { trackId, size = 16, showWhenPaused = true, className = '' }: Props = $props();

  // Reactive state - updates only when player state changes
  let trackState = $derived(player.isTrackPlaying(trackId));
  let isCurrentTrack = $derived(trackState.isCurrentTrack);
  let isPlaying = $derived(trackState.isPlaying);
  let isPaused = $derived(isCurrentTrack && !isPlaying);

  let shouldShowIndicator = $derived(isPlaying || (isPaused && showWhenPaused));
</script>

{#if shouldShowIndicator}
  <div
    class="now-playing-indicator {className}"
    class:playing={isPlaying}
    class:paused={isPaused}
    role="status"
    aria-label={isPlaying ? 'Wird gerade abgespielt' : 'Pausiert'}
  >
    {#if isPlaying}
      <!-- Animated equalizer bars -->
      <div class="equalizer">
        <span class="bar"></span>
        <span class="bar"></span>
        <span class="bar"></span>
      </div>
    {:else}
      <Pause size={size} />
    {/if}
  </div>
{/if}

<style>
  .now-playing-indicator {
    position: absolute;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    transition: opacity 0.2s ease;
  }

  /* Equalizer Bars */
  .equalizer {
    display: flex;
    align-items: flex-end;
    gap: 2px;
    height: 14px;
  }

  .bar {
    width: 3px;
    background: #8a65ff; /* Brand purple */
    border-radius: 2px;
    animation: equalize 1.2s ease-in-out infinite;
  }

  .bar:nth-child(1) {
    animation-delay: 0s;
  }

  .bar:nth-child(2) {
    animation-delay: 0.2s;
  }

  .bar:nth-child(3) {
    animation-delay: 0.4s;
  }

  @keyframes equalize {
    0%,
    100% {
      height: 4px;
    }
    50% {
      height: 14px;
    }
  }

  /* Paused state styling */
  .paused {
    color: #6b7280; /* Muted gray */
  }

  .playing {
    color: #8a65ff;
  }

  /* Accessibility: Reduced motion */
  @media (prefers-reduced-motion: reduce) {
    .bar {
      animation: none;
      height: 8px; /* Static mid-height */
    }
  }
</style>