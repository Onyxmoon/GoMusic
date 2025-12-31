<script lang="ts">
  import { onMount } from 'svelte';
  import { GetAllTracks, ScanLibrary, GetSources } from '../../../wailsjs/go/main/App.js';
  import { tracks, isLoading, error, searchQuery, filteredTracks, scanProgress } from '../stores/library';
  import { EventsOn } from '../../../wailsjs/runtime/runtime.js';
  import Input from '../components/Input.svelte';
  import Button from '../components/Button.svelte';
  import TrackList from '../components/TrackList.svelte';

  let isScanning = false;
  let firstSourceId: string | null = null;

  onMount(async () => {
    await loadSources();
    await loadTracks();
    setupScanEvents();
  });

  async function loadSources() {
    try {
      const sources = await GetSources();
      if (sources && sources.length > 0) {
        firstSourceId = sources[0].ID;
      }
    } catch (err) {
      console.error('Failed to load sources:', err);
    }
  }

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
    if (!firstSourceId) {
      error.set('No sources configured. Please add a source first.');
      return;
    }

    try {
      isScanning = true;
      await ScanLibrary(firstSourceId);
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

    EventsOn('scan:error', (data: any) => {
      console.error('Scan error:', data);
      error.set(data.error || 'Scan failed');
      isScanning = false;
    });
  }
</script>

<div class="library">
  <div class="library-header">
    <h1 class="library-title">Library</h1>
    <div class="library-actions">
      <Button on:click={startScan} disabled={isScanning}>
        {isScanning ? 'Scanning...' : 'Scan Library'}
      </Button>
      <Button variant="secondary" on:click={loadTracks} disabled={$isLoading}>
        Refresh
      </Button>
    </div>
  </div>

  <div class="search-bar">
    <Input
      bind:value={$searchQuery}
      placeholder="Search tracks, artists, albums..."
    />
  </div>

  {#if $error}
    <div class="error-message">
      {$error}
    </div>
  {/if}

  {#if $isLoading}
    <div class="loading">
      <div class="spinner"></div>
      <p>Loading library...</p>
    </div>
  {:else}
    <div class="library-content">
      <TrackList tracks={$filteredTracks} />
    </div>
  {/if}
</div>

<style>
  .library {
    padding: 24px;
    height: 100%;
    overflow-y: auto;
  }

  .library-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 24px;
  }

  .library-title {
    font-size: 32px;
    font-weight: 700;
    color: #e4e6eb;
    margin: 0;
  }

  .library-actions {
    display: flex;
    gap: 12px;
  }

  .search-bar {
    margin-bottom: 24px;
    max-width: 500px;
  }

  .error-message {
    padding: 16px;
    background: rgba(239, 68, 68, 0.1);
    border: 1px solid rgba(239, 68, 68, 0.3);
    border-radius: 8px;
    color: #fca5a5;
    margin-bottom: 24px;
  }

  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 80px 20px;
    color: #9ca3af;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid #3a3f4b;
    border-top-color: #5b8cff;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    margin-bottom: 16px;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .library-content {
    background: #1e2128;
    border-radius: 8px;
    padding: 16px;
  }
</style>