<script lang="ts">
  import { Menu, Plus, ChevronLeft, Music, Heart } from 'lucide-svelte';

  interface PlaylistItem {
    id: string;
    name: string;
    trackCount: number;
    coverUrl?: string;
  }

  export let playlists: PlaylistItem[] = [
    { id: '1', name: 'Feel Good Hits', trackCount: 121, coverUrl: 'https://via.placeholder.com/60' },
    { id: '2', name: 'Chill Vibes', trackCount: 67, coverUrl: 'https://via.placeholder.com/60' },
    { id: '3', name: 'Focus Mode', trackCount: 17, coverUrl: 'https://via.placeholder.com/60' },
    { id: '4', name: 'Weekend Tunes', trackCount: 45, coverUrl: 'https://via.placeholder.com/60' },
    { id: '5', name: 'Morning Boost', trackCount: 194, coverUrl: 'https://via.placeholder.com/60' },
    { id: '6', name: 'Relax & Unwind', trackCount: 8, coverUrl: 'https://via.placeholder.com/60' },
    { id: '7', name: 'Workout Energy', trackCount: 62, coverUrl: 'https://via.placeholder.com/60' },
    { id: '8', name: 'Party Mix', trackCount: 103, coverUrl: 'https://via.placeholder.com/60' },
    { id: '9', name: 'Rainy Day Sounds', trackCount: 29, coverUrl: 'https://via.placeholder.com/60' }
  ];

  export let selectedPlaylistId: string | null = '2';
</script>

<aside class="playlist-sidebar">
  <div class="sidebar-header">
    <span class="header-icon"><Menu size={16} /></span>
    <h3>Recent</h3>
    <button class="add-btn"><Plus size={14} /></button>
    <button class="collapse-btn"><ChevronLeft size={14} /></button>
  </div>

  <div class="playlist-list">
    {#each playlists as playlist}
      <button
        class="playlist-item"
        class:active={selectedPlaylistId === playlist.id}
        on:click={() => selectedPlaylistId = playlist.id}
      >
        <div class="playlist-cover">
          {#if playlist.coverUrl}
            <img src={playlist.coverUrl} alt={playlist.name} />
          {:else}
            <div class="cover-placeholder"><Music size={24} /></div>
          {/if}
          <div class="track-count">{playlist.trackCount}</div>
        </div>
        <div class="playlist-info">
          <div class="playlist-name">{playlist.name}</div>
          <div class="playlist-type">Playlist</div>
        </div>
      </button>
    {/each}
  </div>

  <div class="now-playing">
    <img src="https://via.placeholder.com/48" alt="Current track" />
    <div class="track-info">
      <div class="track-title">Cruel Summer</div>
      <div class="track-artist">Taylor Swift</div>
    </div>
    <button class="heart-btn"><Heart size={16} /></button>
  </div>
</aside>

<style>
  .playlist-sidebar {
    width: 220px;
    height: 100%;
    background: rgba(255, 255, 255, 0.4);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border-right: 1px solid rgba(255, 255, 255, 0.3);
    display: flex;
    flex-direction: column;
  }

  .sidebar-header {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 16px;
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  }

  .header-icon {
    font-size: 16px;
    color: #6b7280;
  }

  .sidebar-header h3 {
    flex: 1;
    font-size: 14px;
    font-weight: 600;
    color: #2d2d2d;
  }

  .add-btn,
  .collapse-btn {
    width: 24px;
    height: 24px;
    border-radius: 4px;
    border: none;
    background: rgba(0, 0, 0, 0.05);
    color: #6b7280;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
  }

  .add-btn:hover,
  .collapse-btn:hover {
    background: rgba(0, 0, 0, 0.1);
  }

  .playlist-list {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
  }

  .playlist-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    margin-bottom: 4px;
    text-align: left;
    transition: all 0.2s ease;
  }

  .playlist-item:hover {
    background: rgba(255, 255, 255, 0.6);
  }

  .playlist-item.active {
    background: rgba(255, 255, 255, 0.8);
  }

  .playlist-cover {
    width: 52px;
    height: 52px;
    border-radius: 8px;
    overflow: hidden;
    position: relative;
    flex-shrink: 0;
  }

  .playlist-cover img {
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

  .track-count {
    position: absolute;
    bottom: 4px;
    left: 4px;
    background: rgba(0, 0, 0, 0.7);
    color: white;
    padding: 2px 6px;
    border-radius: 4px;
    font-size: 10px;
    font-weight: 600;
  }

  .playlist-info {
    flex: 1;
    min-width: 0;
  }

  .playlist-name {
    font-size: 13px;
    font-weight: 500;
    color: #2d2d2d;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .playlist-type {
    font-size: 11px;
    color: #6b7280;
  }

  .now-playing {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    border-top: 1px solid rgba(0, 0, 0, 0.06);
    background: rgba(255, 255, 255, 0.5);
  }

  .now-playing img {
    width: 48px;
    height: 48px;
    border-radius: 6px;
  }

  .track-info {
    flex: 1;
    min-width: 0;
  }

  .track-title {
    font-size: 12px;
    font-weight: 500;
    color: #2d2d2d;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .track-artist {
    font-size: 11px;
    color: #6b7280;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .heart-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: #6b7280;
    font-size: 16px;
    cursor: pointer;
  }

  .heart-btn:hover {
    color: #ff6b9d;
  }
</style>
