<script lang="ts">
  import { MoreVertical, User } from 'lucide-svelte';

  interface Artist {
    id: string;
    name: string;
    listeners: string;
    albums: number;
    avatarUrl?: string;
  }

  export let topArtists: Artist[] = [
    { id: '1', name: 'Taylor Swift', listeners: '221M listeners', albums: 14, avatarUrl: 'https://via.placeholder.com/48' },
    { id: '2', name: 'The Weeknd', listeners: '190M listeners', albums: 5, avatarUrl: 'https://via.placeholder.com/48' },
    { id: '3', name: 'Taylor Swift', listeners: '156M listeners', albums: 7, avatarUrl: 'https://via.placeholder.com/48' },
    { id: '4', name: 'Lana Del Rey', listeners: '142M listeners', albums: 9, avatarUrl: 'https://via.placeholder.com/48' },
    { id: '5', name: 'Bruno Mars', listeners: '119M listeners', albums: 4, avatarUrl: 'https://via.placeholder.com/48' },
    { id: '6', name: 'Lady Gaga', listeners: '104M listeners', albums: 7, avatarUrl: 'https://via.placeholder.com/48' },
    { id: '7', name: 'Eminem', listeners: '98M listeners', albums: 14, avatarUrl: 'https://via.placeholder.com/48' },
    { id: '8', name: 'Drake', listeners: '85M listeners', albums: 8, avatarUrl: 'https://via.placeholder.com/48' }
  ];
</script>

<aside class="right-sidebar">
  <div class="top-artists">
    <div class="section-header">
      <h3>Top Artists</h3>
      <button class="more-btn"><MoreVertical size={16} /></button>
    </div>

    <div class="artist-list">
      {#each topArtists as artist}
        <button class="artist-item">
          <div class="artist-avatar">
            {#if artist.avatarUrl}
              <img src={artist.avatarUrl} alt={artist.name} />
            {:else}
              <div class="avatar-placeholder"><User size={24} /></div>
            {/if}
          </div>
          <div class="artist-info">
            <div class="artist-name">{artist.name}</div>
            <div class="artist-stats">
              {artist.listeners} • {artist.albums} albums
            </div>
          </div>
        </button>
      {/each}
    </div>
  </div>

  <div class="premium-banner">
    <h3>Uninterrupted Music Awaits</h3>
    <p>Enjoy ad-free music with our premium plan.</p>
    <button class="upgrade-btn">Upgrade now</button>
  </div>
</aside>

<style>
  .right-sidebar {
    width: 280px;
    height: 100%;
    background: rgba(255, 255, 255, 0.4);
    backdrop-filter: blur(20px);
    -webkit-backdrop-filter: blur(20px);
    border-left: 1px solid rgba(255, 255, 255, 0.3);
    display: flex;
    flex-direction: column;
    overflow-y: auto;
  }

  .top-artists {
    padding: 20px;
    flex: 1;
  }

  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
  }

  .section-header h3 {
    font-size: 16px;
    font-weight: 600;
    color: #2d2d2d;
  }

  .more-btn {
    width: 28px;
    height: 28px;
    border: none;
    background: transparent;
    border-radius: 4px;
    color: #6b7280;
    font-size: 16px;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .more-btn:hover {
    background: rgba(0, 0, 0, 0.05);
  }

  .artist-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .artist-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 8px;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    text-align: left;
    transition: all 0.15s ease;
  }

  .artist-item:hover {
    background: rgba(255, 255, 255, 0.6);
  }

  .artist-avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    overflow: hidden;
    flex-shrink: 0;
  }

  .artist-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .avatar-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #8a65ff, #ff6b9d);
    font-size: 24px;
  }

  .artist-info {
    flex: 1;
    min-width: 0;
  }

  .artist-name {
    font-size: 14px;
    font-weight: 500;
    color: #2d2d2d;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .artist-stats {
    font-size: 12px;
    color: #6b7280;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .premium-banner {
    margin: 20px;
    padding: 20px;
    background: linear-gradient(135deg, rgba(138, 101, 255, 0.15), rgba(255, 107, 157, 0.15));
    backdrop-filter: blur(10px);
    -webkit-backdrop-filter: blur(10px);
    border-radius: 12px;
    border: 1px solid rgba(255, 255, 255, 0.4);
    text-align: center;
  }

  .premium-banner h3 {
    font-size: 16px;
    font-weight: 600;
    color: #2d2d2d;
    margin-bottom: 8px;
  }

  .premium-banner p {
    font-size: 13px;
    color: #6b7280;
    margin-bottom: 16px;
  }

  .upgrade-btn {
    width: 100%;
    padding: 10px 20px;
    background: linear-gradient(135deg, #8a65ff, #ff6b9d);
    border: none;
    border-radius: 8px;
    color: white;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .upgrade-btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 8px 20px rgba(138, 101, 255, 0.3);
  }

  .upgrade-btn:active {
    transform: translateY(0);
  }
</style>
