<script lang="ts">
  import { Heart } from 'lucide-svelte';
  import type { dto } from '../../../../wailsjs/go/models';

  interface Props {
    track: dto.TrackDTO;
  }

  let { track }: Props = $props();
</script>

<div class="track-info-container">
  <div class="track-cover">
    {#if track.hasArtwork}
      <img src={`/artwork/stream?id=${encodeURIComponent(track.id)}`} alt={track.title} />
    {:else}
      <div class="cover-placeholder"></div>
    {/if}
  </div>
  <div class="track-info">
    <div class="track-title">{track.title}</div>
    <div class="track-artist">{track.artist}</div>
  </div>
  <button class="heart-btn" aria-label="Add to favorites">
    <Heart size={18} />
  </button>
</div>

<style>
  .track-info-container {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;
    min-width: 0;
    position: relative;
    z-index: 1;
  }

  .track-cover {
    width: 56px;
    height: 56px;
    border-radius: 8px;
    overflow: hidden;
    flex-shrink: 0;
  }

  .track-cover img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .cover-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #8a65ff, #ff6b9d);
    font-size: 24px;
  }

  .track-info {
    flex: 1;
    min-width: 0;
  }

  .track-title {
    font-size: 14px;
    font-weight: 500;
    color: #2d2d2d;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .track-artist {
    font-size: 12px;
    color: #31343a;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .heart-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: #31343a;
    font-size: 18px;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .heart-btn:hover {
    color: #ff5686;
  }
</style>