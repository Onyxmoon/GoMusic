<script lang="ts">
  import type { TrackDTO } from '../types/track';
  import { Play, MoreVertical } from 'lucide-svelte';

  export let tracks: TrackDTO[] = [];
  export let onTrackClick: (track: TrackDTO) => void = () => {};

  function formatDuration(seconds: number): string {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  }

  function formatDate(date: string): string {
    // Simple date formatting
    return date || '22.6.2024';
  }
</script>

<div class="track-table">
  <div class="table-header">
    <div class="col col-number">#</div>
    <div class="col col-title">Title</div>
    <div class="col col-album">Album</div>
    <div class="col col-date">Added date</div>
    <div class="col col-time">Time</div>
  </div>

  <div class="table-body">
    {#each tracks as track, index}
      <div
        class="track-row"
        on:click={() => onTrackClick(track)}
        role="button"
        tabindex="0"
        on:keydown={(e) => e.key === 'Enter' && onTrackClick(track)}
      >
        <div class="col col-number">
          <span class="play-icon"><Play size={14} /></span>
          <span class="track-number">{index + 1}</span>
        </div>

        <div class="col col-title">
          <div class="track-cover">
            {#if track.artworkPath}
              <img src={`/artwork/stream?file=${encodeURIComponent(track.artworkPath)}`} alt={track.title} />
            {:else}
              <div class="cover-placeholder"></div>
            {/if}
          </div>
          <div class="title-info">
            <div class="title">{track.title}</div>
            <div class="artist">{track.artist}</div>
          </div>
        </div>

        <div class="col col-album">{track.album}</div>
        <div class="col col-date">{formatDate(track.addedAt?.toString())}</div>
        <div class="col col-time">
          {formatDuration(track.duration)}
          <button class="more-btn" on:click|stopPropagation><MoreVertical size={16} /></button>
        </div>
      </div>
    {/each}
  </div>
</div>

<style>
  .track-table {
    width: 100%;
  }

  .table-header {
    display: grid;
    grid-template-columns: 50px minmax(200px, 2fr) minmax(120px, 1fr) minmax(100px, 0.8fr) 100px;
    gap: 16px;
    padding: 12px 16px;
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
    font-size: 12px;
    color: #6b7280;
    font-weight: 500;
  }

  .table-body {
    display: flex;
    flex-direction: column;
  }

  .track-row {
    display: grid;
    grid-template-columns: 50px minmax(200px, 2fr) minmax(120px, 1fr) minmax(100px, 0.8fr) 100px;
    gap: 16px;
    padding: 10px 16px;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    text-align: left;
    transition: all 0.15s ease;
  }

  .track-row:hover {
    background: rgba(255, 255, 255, 0.5);
  }

  .track-row:hover .play-icon {
    opacity: 1;
  }

  .track-row:hover .track-number {
    opacity: 0;
  }

  .col {
    display: flex;
    align-items: center;
    color: #2d2d2d;
    font-size: 14px;
  }

  .col-number {
    position: relative;
    justify-content: center;
    color: #6b7280;
    font-weight: 500;
  }

  .play-icon {
    position: absolute;
    opacity: 0;
    transition: opacity 0.15s ease;
    font-size: 12px;
    color: #2d2d2d;
  }

  .track-number {
    transition: opacity 0.15s ease;
  }

  .col-title {
    gap: 12px;
  }

  .track-cover {
    width: 40px;
    height: 40px;
    border-radius: 6px;
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
    font-size: 18px;
  }

  .title-info {
    flex: 1;
    min-width: 0;
  }

  .title {
    font-weight: 500;
    color: #2d2d2d;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .artist {
    font-size: 12px;
    color: #6b7280;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .col-album,
  .col-date {
    color: #6b7280;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    min-width: 0;
  }

  .col-time {
    justify-content: space-between;
    color: #6b7280;
  }

  .more-btn {
    width: 28px;
    height: 28px;
    border: none;
    background: transparent;
    border-radius: 4px;
    color: #6b7280;
    cursor: pointer;
    opacity: 0;
    transition: all 0.15s ease;
  }

  .track-row:hover .more-btn {
    opacity: 1;
  }

  .more-btn:hover {
    background: rgba(0, 0, 0, 0.05);
    color: #2d2d2d;
  }
</style>
