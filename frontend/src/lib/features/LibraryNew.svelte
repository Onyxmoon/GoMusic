<script lang="ts">
  import { onMount } from 'svelte';
  import { GetAllTracks, ScanAllLibraries } from '../../../wailsjs/go/main/App.js';
  import { tracks, isLoading, error } from '../stores/library';
  import { EventsOn } from '../../../wailsjs/runtime/runtime.js';
  import type { dto } from '../../../wailsjs/go/models';
  import type { ScanErrorEvent } from '../types/events';
  import TrackTable from '../components/TrackTable.svelte';
  import { player } from '../stores/player.svelte';
  import { RefreshCw, AlertTriangle, X, Music, Disc, Mic, List } from 'lucide-svelte';

  export let section: string = 'tracks';

  let isScanning = false;

  onMount(async () => {
    await loadTracks();
    setupScanEvents();
  });

  async function loadTracks() {
    try {
      isLoading.set(true);
      error.set(null);
      const loadedTracks = await GetAllTracks();
      tracks.set(loadedTracks || []);
    } catch (err) {
      error.set(err instanceof Error ? err.message : 'Failed to load tracks');
      console.error('Failed to load tracks:', err);
    } finally {
      isLoading.set(false);
    }
  }

  async function startScan() {
    try {
      isScanning = true;
      await ScanAllLibraries();
    } catch (err) {
      error.set(err instanceof Error ? err.message : 'Failed to start scan');
      console.error('Failed to start scan:', err);
      isScanning = false;
    }
  }

  function setupScanEvents() {
    EventsOn('scan:started', (sourceId: string) => {
      console.log('Scan started:', sourceId);
      isScanning = true;
    });

    EventsOn('scan:complete', async (sourceId: string) => {
      console.log('Scan completed:', sourceId);
      isScanning = false;
      await loadTracks();
    });

    EventsOn('scan:error', (data: ScanErrorEvent) => {
      console.error('Scan error:', data);
      error.set(data.error || 'Scan failed');
      isScanning = false;
    });
  }

  function handlePlayAll() {
    if ($tracks.length > 0) {
      player.setPlaylist($tracks, 0);
      player.play($tracks[0]);
    }
  }

  function handleTrackClick(track: dto.TrackDTO) {
    player.setPlaylist($tracks, $tracks.indexOf(track));
    player.play(track);
  }
</script>

<div class="library">
  {#if section === 'tracks'}
    <div class="section-header">
      <h2>All Tracks</h2>
      <button class="scan-btn" on:click={startScan} disabled={isScanning}>
        {#if isScanning}
          Scanning...
        {:else}
          <RefreshCw size={16} style="margin-right: 8px;" />
          Scan Library
        {/if}
      </button>
    </div>

    {#if $error}
      <div class="error-banner">
        <AlertTriangle size={18} />
        <span>{$error}</span>
        <button on:click={() => error.set(null)}><X size={16} /></button>
      </div>
    {/if}

    {#if $isLoading}
      <div class="loading-state">
        <div class="spinner"></div>
        <p>Loading your music library...</p>
      </div>
    {:else if $tracks.length === 0}
      <div class="empty-state">
        <div class="empty-icon"><Music size={64} /></div>
        <h3>No tracks yet</h3>
        <p>Scan your music library to get started</p>
        <button class="scan-btn primary" on:click={startScan} disabled={isScanning}>
          {isScanning ? 'Scanning...' : 'Scan Library'}
        </button>
      </div>
    {:else}
      <div class="track-table-container">
        <TrackTable tracks={$tracks} onTrackClick={handleTrackClick} />
      </div>
    {/if}

  {:else if section === 'albums'}
    <div class="section-header">
      <h2>Albums</h2>
    </div>
    <div class="placeholder-view">
      <div class="placeholder-icon"><Disc size={64} /></div>
      <h3>Albums View</h3>
      <p>Album grid view coming soon</p>
    </div>

  {:else if section === 'artists'}
    <div class="section-header">
      <h2>Artists</h2>
    </div>
    <div class="placeholder-view">
      <div class="placeholder-icon"><Mic size={64} /></div>
      <h3>Artists View</h3>
      <p>Artist grid view coming soon</p>
    </div>

  {:else if section === 'playlists'}
    <div class="section-header">
      <h2>Playlists</h2>
      <button class="action-btn">+ New Playlist</button>
    </div>
    <div class="placeholder-view">
      <div class="placeholder-icon"><List size={64} /></div>
      <h3>Playlists View</h3>
      <p>Playlist management coming soon</p>
    </div>
  {/if}
</div>

<style>
  .library {
    width: 100%;
    max-width: 1200px;
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  .error-banner {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background: rgba(239, 68, 68, 0.15);
    border: 1px solid rgba(239, 68, 68, 0.3);
    border-radius: 8px;
    color: #dc2626;
    margin-bottom: 16px;
  }

  .error-banner button {
    margin-left: auto;
    background: transparent;
    border: none;
    color: #dc2626;
    cursor: pointer;
    font-size: 16px;
  }

  .loading-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 80px 20px;
    color: #6b7280;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(138, 101, 255, 0.2);
    border-top-color: #8a65ff;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    margin-bottom: 16px;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 80px 20px;
    text-align: center;
  }

  .empty-icon {
    font-size: 64px;
    margin-bottom: 16px;
    opacity: 0.5;
  }

  .empty-state h3 {
    font-size: 20px;
    font-weight: 600;
    color: #2d2d2d;
    margin-bottom: 8px;
  }

  .empty-state p {
    font-size: 14px;
    color: #6b7280;
    margin-bottom: 24px;
  }

  .scan-btn {
    padding: 12px 32px;
    background: linear-gradient(135deg, #8a65ff, #ff6b9d);
    border: none;
    border-radius: 8px;
    color: white;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .scan-btn:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 8px 20px rgba(138, 101, 255, 0.3);
  }

  .scan-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .track-table-container {
    flex: 1;
    overflow-y: auto;
    overflow-x: hidden;
    min-height: 0;

    /* Firefox */
    scrollbar-width: thin;
    scrollbar-color: transparent transparent;
  }

  .track-table-container:hover {
    scrollbar-color: rgba(138, 101, 255, 0.6) rgba(0, 0, 0, 0.05);
  }

  /* WebKit browsers (Chrome, Safari, Edge) */
  .track-table-container::-webkit-scrollbar {
    width: 8px;
  }

  .track-table-container::-webkit-scrollbar-track {
    background: transparent;
    transition: background-color 0.2s ease;
  }

  .track-table-container::-webkit-scrollbar-thumb {
    background-color: transparent;
    border-radius: 4px;
    transition: background-color 0.2s ease;
  }

  .track-table-container:hover::-webkit-scrollbar-track {
    background-color: rgba(0, 0, 0, 0.05);
  }

  .track-table-container:hover::-webkit-scrollbar-thumb {
    background-color: rgba(138, 101, 255, 0.5);
  }

  .track-table-container::-webkit-scrollbar-thumb:hover {
    background-color: rgba(138, 101, 255, 0.8);
  }

  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 24px;
    flex-shrink: 0;
  }

  .section-header h2 {
    font-size: 28px;
    font-weight: 700;
    color: #2d2d2d;
  }

  .action-btn {
    padding: 10px 20px;
    background: rgba(255, 255, 255, 0.7);
    border: 1px solid rgba(0, 0, 0, 0.1);
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    color: #2d2d2d;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .action-btn:hover {
    background: rgba(255, 255, 255, 0.9);
    border-color: #8a65ff;
  }

  .placeholder-view {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 120px 20px;
    text-align: center;
    background: rgba(255, 255, 255, 0.3);
    border-radius: 16px;
    border: 2px dashed rgba(0, 0, 0, 0.1);
  }

  .placeholder-icon {
    font-size: 80px;
    margin-bottom: 24px;
    opacity: 0.3;
  }

  .placeholder-view h3 {
    font-size: 24px;
    font-weight: 600;
    color: #2d2d2d;
    margin-bottom: 8px;
  }

  .placeholder-view p {
    font-size: 14px;
    color: #6b7280;
  }
</style>
