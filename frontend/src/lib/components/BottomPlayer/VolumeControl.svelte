<script lang="ts">
  import { Volume2, Volume1, VolumeX } from 'lucide-svelte';

  interface Props {
    volume: number;
    isMuted: boolean;
    onVolumeChange: (volume: number) => void;
    onToggleMute: () => void;
  }

  let { volume, isMuted, onVolumeChange, onToggleMute }: Props = $props();

  let isExpanded = $state(false);

  function handleVolumeChange(e: Event) {
    const target = e.target as HTMLInputElement;
    const value = parseInt(target.value) / 100;
    onVolumeChange(value);

    // Update CSS custom property for visual feedback
    target.style.setProperty('--volume-percent', `${target.value}%`);
  }

  function handleMouseEnter() {
    isExpanded = true;
  }

  function handleMouseLeave() {
    isExpanded = false;
  }
</script>

<div class="volume-container"
     role="group"
     aria-label="Volume control"
     class:expanded={isExpanded}
     onmouseenter={handleMouseEnter}
     onmouseleave={handleMouseLeave}>
  <span class="volume-value">{Math.round(volume * 100)}%</span>
  <input
    type="range"
    class="volume-slider"
    min="0"
    max="100"
    value={volume * 100}
    oninput={handleVolumeChange}
    style="--volume-percent: {volume * 100}%" />
  <button class="icon-btn volume-btn" onclick={onToggleMute}>
    {#if isMuted || volume === 0}
      <VolumeX size={16} />
    {:else if volume < 0.5}
      <Volume1 size={16} />
    {:else}
      <Volume2 size={16} />
    {/if}
  </button>
</div>

<style>
  .volume-container {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    height: 32px;
    padding: 0;
    border-radius: 6px;
    background: transparent;
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    max-width: 32px;
  }

  .volume-container.expanded {
    max-width: 160px;
    background: rgba(0, 0, 0, 0.03);
    padding-left: 8px;
  }

  .volume-value {
    font-size: 10px;
    color: #40434c;
    font-weight: 500;
    min-width: 0;
    width: 0;
    opacity: 0;
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    overflow: hidden;
    white-space: nowrap;
  }

  .expanded .volume-value {
    width: 32px;
    margin-right: 6px;
    opacity: 1;
  }

  .volume-slider {
    -webkit-appearance: none;
    appearance: none;
    width: 0;
    height: 4px;
    background: transparent;
    border-radius: 2px;
    outline: none;
    opacity: 0;
    transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
    cursor: pointer;
    margin: 0;
    padding: 0;
  }

  .expanded .volume-slider {
    width: 80px;
    margin-right: 6px;
    opacity: 1;
  }

  .volume-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 12px;
    height: 12px;
    background: #2d2d2d;
    border-radius: 50%;
    cursor: pointer;
    transition: all 0.15s ease;
    margin-top: -4px;
  }

  .volume-slider::-webkit-slider-thumb:hover {
    transform: scale(1.2);
  }

  .volume-slider::-webkit-slider-runnable-track {
    background: linear-gradient(
      to right,
      #2d2d2d 0%,
      #2d2d2d var(--volume-percent, 80%),
      rgba(0, 0, 0, 0.1) var(--volume-percent, 80%),
      rgba(0, 0, 0, 0.1) 100%
    );
    border-radius: 2px;
    height: 4px;
  }

  .icon-btn {
    width: 32px;
    height: 32px;
    border-radius: 6px;
    border: none;
    background: transparent;
    color: #40434c;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .icon-btn:hover {
    background: rgba(0, 0, 0, 0.05);
    color: #2d2d2d;
  }

  .volume-btn {
    flex-shrink: 0;
  }
</style>