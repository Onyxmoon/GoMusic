<script lang="ts">
  import { onMount } from 'svelte';
  import { GetSources } from '../../../wailsjs/go/main/App.js';
  import Button from '../components/Button.svelte';
  import Card from '../components/Card.svelte';

  interface Source {
    id: string;
    name: string;
    type: string;
  }

  let sources: Source[] = [];
  let isLoading = false;
  let selectedSource: Source | null = null;
  let showAddDialog = false;
  let isEditMode = false;

  onMount(async () => {
    await loadSources();
  });

  async function loadSources() {
    try {
      isLoading = true;
      sources = await GetSources();
    } catch (err) {
      console.error('Failed to load sources:', err);
    } finally {
      isLoading = false;
    }
  }

  function addSource() {
    showAddDialog = true;
  }

  function editSource(source: Source) {
    console.log('Editing source:', source);
    isEditMode = true;
    selectedSource = source;
  }

  function createNewSource() {
    isEditMode = false;
    showAddDialog = false;
    selectedSource = { id: 'new', name: 'New Filesystem Source', type: 'filesystem' };
  }
</script>

<div class="sources-view">
  <div class="view-header">
    <h1>Music Sources</h1>
    <Button on:click={addSource}>+ Add Source</Button>
  </div>

  <div class="sources-grid">
    {#if isLoading}
      <div class="loading">
        <div class="spinner"></div>
        <p>Loading sources...</p>
      </div>
    {:else if sources.length === 0}
      <Card>
        <div class="empty-state">
          <div class="empty-icon">📁</div>
          <h3>No sources configured</h3>
          <p>Add a music source to get started</p>
          <Button on:click={addSource}>Add Your First Source</Button>
        </div>
      </Card>
    {:else}
      {#each sources as source}
        <Card clickable on:click={() => editSource(source)}>
          <div class="source-card">
            <div class="source-icon">
              {#if source.type === 'filesystem'}
                📂
              {:else if source.type === 'api-selfhosted'}
                🌐
              {:else}
                🎵
              {/if}
            </div>
            <div class="source-info">
              <h3>{source.name || source.id}</h3>
              <p class="source-type">{source.type}</p>
            </div>
            <div class="source-actions">
              <button class="icon-btn">⚙️</button>
            </div>
          </div>
        </Card>
      {/each}
    {/if}
  </div>

  {#if showAddDialog}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="dialog-overlay" on:click={() => showAddDialog = false}>
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
      <div class="dialog" on:click|stopPropagation>
        <h2>Add Music Source</h2>
        <p>Choose a source type:</p>

        <div class="source-types">
          <button class="source-type-btn" on:click={createNewSource}>
            <span class="type-icon">📂</span>
            <span class="type-name">Local Files</span>
            <span class="type-desc">Scan music files from your computer</span>
          </button>

          <button class="source-type-btn" disabled>
            <span class="type-icon">🌐</span>
            <span class="type-name">Self-hosted API</span>
            <span class="type-desc">Connect to your music server</span>
          </button>
        </div>

        <Button variant="secondary" on:click={() => showAddDialog = false}>Cancel</Button>
      </div>
    </div>
  {/if}

  {#if selectedSource}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div class="dialog-overlay" on:click={() => selectedSource = null}>
      <!-- svelte-ignore a11y-click-events-have-key-events -->
      <!-- svelte-ignore a11y-no-static-element-interactions -->
      <div class="dialog dialog-large" on:click|stopPropagation>
        {#if selectedSource.type === 'filesystem'}
          {#await import('./FilesystemConfig.svelte')}
            <p>Loading...</p>
          {:then module}
            <svelte:component
              this={module.default}
              sourceId={selectedSource.id}
              sourceName={selectedSource.name}
              isEditMode={isEditMode}
              onSave={async (config) => {
                try {
                  if (isEditMode && selectedSource) {
                    const { UpdateFilesystemSource } = await import('../../../wailsjs/go/main/App.js');
                    await UpdateFilesystemSource(
                      selectedSource.id,
                      config.name,
                      config.rootPaths,
                      config.includeSubfolders,
                      config.supportedFormats
                    );
                  } else {
                    const { AddFilesystemSource } = await import('../../../wailsjs/go/main/App.js');
                    await AddFilesystemSource(
                      config.name,
                      config.rootPaths,
                      config.includeSubfolders,
                      config.supportedFormats
                    );
                  }
                  selectedSource = null;
                  isEditMode = false;
                  await loadSources();
                } catch (err) {
                  console.error('Failed to save source:', err);
                  alert(`Error: ${err}`);
                }
              }}
              onDelete={async () => {
                if (!selectedSource || !confirm(`Are you sure you want to remove "${selectedSource.name}"?`)) {
                  return;
                }
                try {
                  const { RemoveSource } = await import('../../../wailsjs/go/main/App.js');
                  await RemoveSource(selectedSource.id);
                  selectedSource = null;
                  isEditMode = false;
                  await loadSources();
                } catch (err) {
                  console.error('Failed to delete source:', err);
                  alert(`Error: ${err}`);
                }
              }}
              onCancel={() => {
                selectedSource = null;
                isEditMode = false;
              }}
            />
          {/await}
        {/if}
      </div>
    </div>
  {/if}
</div>

<style>
  .sources-view {
    padding: 24px;
    max-width: 1200px;
  }

  .view-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 32px;
  }

  .view-header h1 {
    font-size: 32px;
    font-weight: 700;
    color: #2d2d2d;
  }

  .sources-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 16px;
  }

  .source-card {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 8px;
  }

  .source-icon {
    font-size: 48px;
    flex-shrink: 0;
  }

  .source-info {
    flex: 1;
    min-width: 0;
  }

  .source-info h3 {
    font-size: 16px;
    font-weight: 600;
    color: #2d2d2d;
    margin-bottom: 4px;
  }

  .source-type {
    font-size: 13px;
    color: #6b7280;
    text-transform: capitalize;
  }

  .source-actions {
    display: flex;
    gap: 8px;
  }

  .icon-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .icon-btn:hover {
    background: rgba(0, 0, 0, 0.05);
  }

  .empty-state {
    text-align: center;
    padding: 60px 20px;
  }

  .empty-icon {
    font-size: 64px;
    margin-bottom: 16px;
  }

  .empty-state h3 {
    font-size: 20px;
    font-weight: 600;
    color: #2d2d2d;
    margin-bottom: 8px;
  }

  .empty-state p {
    color: #6b7280;
    margin-bottom: 24px;
  }

  .loading {
    text-align: center;
    padding: 60px;
    color: #6b7280;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(138, 101, 255, 0.2);
    border-top-color: #8a65ff;
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
    margin: 0 auto 16px;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .dialog-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .dialog {
    background: rgba(255, 255, 255, 0.95);
    backdrop-filter: blur(20px);
    border-radius: 16px;
    padding: 32px;
    max-width: 500px;
    width: 90%;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  }

  .dialog-large {
    max-width: 800px;
  }

  .dialog h2 {
    font-size: 24px;
    font-weight: 700;
    color: #2d2d2d;
    margin-bottom: 8px;
  }

  .dialog p {
    color: #6b7280;
    margin-bottom: 24px;
  }

  .source-types {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-bottom: 24px;
  }

  .source-type-btn {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding: 20px;
    background: rgba(255, 255, 255, 0.7);
    border: 2px solid rgba(0, 0, 0, 0.1);
    border-radius: 12px;
    cursor: pointer;
    transition: all 0.2s ease;
    text-align: left;
  }

  .source-type-btn:not(:disabled):hover {
    background: rgba(255, 255, 255, 0.9);
    border-color: #8a65ff;
    transform: translateY(-2px);
  }

  .source-type-btn:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  .type-icon {
    font-size: 32px;
    margin-bottom: 8px;
  }

  .type-name {
    font-size: 16px;
    font-weight: 600;
    color: #2d2d2d;
    display: block;
    margin-bottom: 4px;
  }

  .type-desc {
    font-size: 13px;
    color: #6b7280;
    display: block;
  }
</style>