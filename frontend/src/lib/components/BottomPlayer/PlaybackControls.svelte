<script lang="ts">
  import { SkipBack, Play, Pause, SkipForward } from 'lucide-svelte';

  interface Props {
    isPlaying: boolean;
    onPlayPause: () => void;
    onPrevious: () => void;
    onNext: () => void;
  }

  let { isPlaying, onPlayPause, onPrevious, onNext }: Props = $props();
</script>

<div class="controls">
  <button class="control-btn" onclick={onPrevious} aria-label="Previous track">
    <SkipBack size={18} />
  </button>
  <button class="control-btn play-btn" onclick={onPlayPause} aria-label={isPlaying ? 'Pause' : 'Play'}>
    {#if isPlaying}
      <Pause size={18} />
    {:else}
      <Play size={18} />
    {/if}
  </button>
  <button class="control-btn" onclick={onNext} aria-label="Next track">
    <SkipForward size={18} />
  </button>
</div>

<style>
  .controls {
    display: flex;
    gap: 12px;
    align-items: center;
  }

  .control-btn {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border: none;
    background: transparent;
    color: #2d2d2d;
    font-size: 16px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.15s ease;
  }

  .control-btn:hover {
    background: rgba(0, 0, 0, 0.05);
  }

  .play-btn {
    width: 40px;
    height: 40px;
    background: #2d2d2d;
    color: white;
  }

  .play-btn:hover {
    background: #1a1a1a;
  }
</style>