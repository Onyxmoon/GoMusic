<script lang="ts">
  import type { TrackDTO } from '../types/track';
  import { player } from '../stores/player.svelte';
  import { Play } from 'lucide-svelte';

  export let tracks: TrackDTO[];
  export let showAlbum: boolean = true;
  export let showArtist: boolean = true;

  function formatDuration(seconds: number): string {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  }

  function playTrack(track: TrackDTO) {
    player.play(track);
  }

  function playAll() {
    if (tracks.length > 0) {
      player.setPlaylist(tracks, 0);
      player.play(tracks[0]);
    }
  }
</script>

<div class="track-list">
  {#if tracks.length > 0}
    <div class="track-list-header">
      <button class="play-all-btn" on:click={playAll}>
        <Play size={16} style="margin-right: 6px;" />
        Play All
      </button>
      <span class="track-count">{tracks.length} tracks</span>
    </div>

    <div class="track-list-table">
      <div class="track-list-thead">
        <div class="track-col track-col-number">#</div>
        <div class="track-col track-col-title">Title</div>
        {#if showArtist}
          <div class="track-col track-col-artist">Artist</div>
        {/if}
        {#if showAlbum}
          <div class="track-col track-col-album">Album</div>
        {/if}
        <div class="track-col track-col-duration">Duration</div>
      </div>

      {#each tracks as track, index}
        <div class="track-row" on:click={() => playTrack(track)}>
          <div class="track-col track-col-number">{index + 1}</div>
          <div class="track-col track-col-title">
            <div class="track-title">{track.title}</div>
          </div>
          {#if showArtist}
            <div class="track-col track-col-artist">{track.artist}</div>
          {/if}
          {#if showAlbum}
            <div class="track-col track-col-album">{track.album}</div>
          {/if}
          <div class="track-col track-col-duration">
            {formatDuration(track.duration)}
          </div>
        </div>
      {/each}
    </div>
  {:else}
    <div class="empty-state">
      <p>No tracks found</p>
    </div>
  {/if}
</div>

<style>
  .track-list {
    width: 100%;
  }

  .track-list-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 0;
    margin-bottom: 8px;
  }

  .play-all-btn {
    padding: 10px 20px;
    background: #5b8cff;
    border: none;
    border-radius: 6px;
    color: white;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .play-all-btn:hover {
    background: #4a7de8;
  }

  .track-count {
    color: #9ca3af;
    font-size: 14px;
  }

  .track-list-table {
    width: 100%;
  }

  .track-list-thead {
    display: grid;
    grid-template-columns: 50px 1fr 1fr 1fr 100px;
    padding: 12px 16px;
    border-bottom: 1px solid #3a3f4b;
    color: #9ca3af;
    font-size: 13px;
    font-weight: 500;
    text-transform: uppercase;
  }

  .track-row {
    display: grid;
    grid-template-columns: 50px 1fr 1fr 1fr 100px;
    padding: 12px 16px;
    border-radius: 4px;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .track-row:hover {
    background: #2a2e39;
  }

  .track-col {
    display: flex;
    align-items: center;
    color: #e4e6eb;
    font-size: 14px;
  }

  .track-col-number {
    color: #9ca3af;
  }

  .track-title {
    font-weight: 500;
  }

  .track-col-duration {
    justify-content: flex-end;
    color: #9ca3af;
  }

  .empty-state {
    padding: 60px 20px;
    text-align: center;
    color: #6b7280;
  }
</style>