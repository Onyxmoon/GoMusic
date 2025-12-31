<script lang="ts">
  import { player } from '../stores/player.svelte';
  import { SkipBack, Play, Pause, SkipForward, Volume2 } from 'lucide-svelte';
  import { formatTime } from '../utils/timeFormat';

  function handleProgressClick(e: MouseEvent) {
    const target = e.currentTarget as HTMLElement;
    const rect = target.getBoundingClientRect();
    const percent = (e.clientX - rect.left) / rect.width;
    const newTime = percent * player.duration;
    player.seek(newTime);
  }

  function handleVolumeChange(e: Event) {
    const target = e.target as HTMLInputElement;
    player.setVolume(parseFloat(target.value));
  }

  function togglePlayPause() {
    if (player.isPlaying) {
      player.pause();
    } else {
      player.resume();
    }
  }
</script>

<div class="player" class:player-active={player.currentTrack}>
  {#if player.currentTrack}
    <div class="player-track-info">
      <div class="track-artwork">
        {#if player.currentTrack.hasArtwork}
          <img src={`/artwork/stream?id=${encodeURIComponent(player.currentTrack.id)}`} alt={player.currentTrack.title} />
        {:else}
          <div class="artwork-placeholder"></div>
        {/if}
      </div>
      <div class="track-details">
        <div class="track-title">{player.currentTrack.title}</div>
        <div class="track-artist">{player.currentTrack.artist}</div>
      </div>
    </div>

    <div class="player-controls">
      <div class="control-buttons">
        <button class="control-btn" on:click={player.previous}>
          <SkipBack size={18} />
        </button>
        <button class="control-btn control-btn-play" on:click={togglePlayPause}>
          {#if player.isPlaying}
            <Pause size={20} />
          {:else}
            <Play size={20} />
          {/if}
        </button>
        <button class="control-btn" on:click={player.next}>
          <SkipForward size={18} />
        </button>
      </div>

      <div class="progress-container">
        <span class="time-display">{formatTime(player.currentTime)}</span>
        <div class="progress-bar" on:click={handleProgressClick}>
          <div class="progress-fill" style="width: {player.progress}%"></div>
        </div>
        <span class="time-display">{formatTime(player.duration)}</span>
      </div>
    </div>

    <div class="player-volume">
      <span class="volume-icon"><Volume2 size={18} /></span>
      <input
        type="range"
        min="0"
        max="1"
        step="0.01"
        value={player.volume}
        on:input={handleVolumeChange}
        class="volume-slider"
      />
    </div>
  {:else}
    <div class="player-empty">
      <p>No track selected</p>
    </div>
  {/if}
</div>

<style>
  .player {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 90px;
    background: #1a1d24;
    border-top: 1px solid #2a2e39;
    display: flex;
    align-items: center;
    padding: 0 24px;
    gap: 24px;
    z-index: 100;
  }

  .player-track-info {
    display: flex;
    align-items: center;
    gap: 12px;
    min-width: 250px;
  }

  .track-artwork {
    width: 56px;
    height: 56px;
    border-radius: 4px;
    overflow: hidden;
    flex-shrink: 0;
  }

  .track-artwork img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .artwork-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #2a2e39;
    font-size: 24px;
  }

  .track-details {
    min-width: 0;
  }

  .track-title {
    font-size: 14px;
    font-weight: 500;
    color: #e4e6eb;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .track-artist {
    font-size: 12px;
    color: #9ca3af;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .player-controls {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
    align-items: center;
  }

  .control-buttons {
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
    color: #e4e6eb;
    font-size: 16px;
    cursor: pointer;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .control-btn:hover {
    background: #2a2e39;
  }

  .control-btn-play {
    width: 44px;
    height: 44px;
    background: #5b8cff;
    font-size: 18px;
  }

  .control-btn-play:hover {
    background: #4a7de8;
  }

  .progress-container {
    display: flex;
    align-items: center;
    gap: 12px;
    width: 100%;
    max-width: 600px;
  }

  .time-display {
    font-size: 12px;
    color: #9ca3af;
    min-width: 40px;
    text-align: center;
  }

  .progress-bar {
    flex: 1;
    height: 4px;
    background: #3a3f4b;
    border-radius: 2px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: #5b8cff;
    transition: width 0.1s linear;
  }

  .player-volume {
    display: flex;
    align-items: center;
    gap: 8px;
    min-width: 150px;
  }

  .volume-icon {
    font-size: 18px;
  }

  .volume-slider {
    flex: 1;
    height: 4px;
    -webkit-appearance: none;
    appearance: none;
    background: #3a3f4b;
    border-radius: 2px;
    outline: none;
  }

  .volume-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background: #5b8cff;
    cursor: pointer;
  }

  .volume-slider::-moz-range-thumb {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background: #5b8cff;
    cursor: pointer;
    border: none;
  }

  .player-empty {
    width: 100%;
    text-align: center;
    color: #6b7280;
  }
</style>