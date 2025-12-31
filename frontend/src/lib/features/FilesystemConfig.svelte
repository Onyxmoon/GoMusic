<script lang="ts">
  import { onMount } from 'svelte';
  import { SelectDirectory, GetSourceConfig } from '../../../wailsjs/go/main/App.js';
  import Button from '../components/Button.svelte';
  import Input from '../components/Input.svelte';

  export let sourceId: string = '';
  export let sourceName: string = '';
  export let isEditMode: boolean = false;
  export let onSave: (config: any) => void = () => {};
  export let onDelete: () => void = () => {};
  export let onCancel: () => void = () => {};

  let name = sourceName || 'My Music Library';
  let rootPaths: string[] = [];
  let includeSubfolders = true;
  let supportedFormats = ['.mp3', '.flac', '.m4a', '.ogg'];
  let isLoading = false;

  onMount(async () => {
    if (isEditMode && sourceId !== 'new') {
      await loadExistingConfig();
    }
  });

  async function loadExistingConfig() {
    try {
      isLoading = true;
      console.log('Loading config for source:', sourceId);
      const config = await GetSourceConfig(sourceId);
      console.log('Loaded config:', config);

      name = config.name || sourceId;

      // Extract root_paths
      if (config.config.root_paths && Array.isArray(config.config.root_paths)) {
        rootPaths = config.config.root_paths;
      } else if (config.config.root_path) {
        rootPaths = [config.config.root_path];
      }

      // Extract include_subfolders
      if (typeof config.config.include_subfolders === 'boolean') {
        includeSubfolders = config.config.include_subfolders;
      }

      // Extract supported_formats
      if (config.config.supported_formats && Array.isArray(config.config.supported_formats)) {
        supportedFormats = config.config.supported_formats;
      }
    } catch (err) {
      console.error('Failed to load source config for sourceId:', sourceId, 'Error:', err);
      alert(`Failed to load configuration for source "${sourceId}": ${err}`);
    } finally {
      isLoading = false;
    }
  }

  async function addFolder() {
    try {
      const selectedPath = await SelectDirectory();
      if (selectedPath && !rootPaths.includes(selectedPath)) {
        rootPaths = [...rootPaths, selectedPath];
      }
    } catch (err) {
      console.error('Failed to select directory:', err);
    }
  }

  function removeFolder(path: string) {
    rootPaths = rootPaths.filter(p => p !== path);
  }

  function toggleFormat(format: string) {
    if (supportedFormats.includes(format)) {
      supportedFormats = supportedFormats.filter(f => f !== format);
    } else {
      supportedFormats = [...supportedFormats, format];
    }
  }

  function handleSave() {
    onSave({
      name,
      rootPaths,
      includeSubfolders,
      supportedFormats
    });
  }
</script>

{#if isLoading}
  <div class="filesystem-config">
    <div class="loading">Loading configuration...</div>
  </div>
{:else}
  <div class="filesystem-config">
    <h2>{isEditMode ? 'Edit' : 'Add'} Filesystem Source</h2>

    <div class="form-section">
      <label>
        <span class="label">Source Name</span>
        <Input bind:value={name} placeholder="My Music Library" />
      </label>
    </div>

  <div class="form-section">
    <div class="section-header">
      <span class="label">Music Folders</span>
      <Button size="sm" on:click={addFolder}>+ Add Folder</Button>
    </div>

    {#if rootPaths.length === 0}
      <div class="empty-folders">
        <p>No folders selected. Click "Add Folder" to select music directories.</p>
      </div>
    {:else}
      <div class="folder-list">
        {#each rootPaths as path}
          <div class="folder-item">
            <span class="folder-icon">📁</span>
            <span class="folder-path">{path}</span>
            <button class="remove-btn" on:click={() => removeFolder(path)}>✕</button>
          </div>
        {/each}
      </div>
    {/if}
  </div>

  <div class="form-section">
    <label class="checkbox-label">
      <input type="checkbox" bind:checked={includeSubfolders} />
      <span>Include subfolders</span>
    </label>
  </div>

  <div class="form-section">
    <span class="label">Supported File Formats</span>
    <div class="format-grid">
      {#each ['.mp3', '.flac', '.m4a', '.ogg', '.wav', '.aac'] as format}
        <button
          class="format-btn"
          class:active={supportedFormats.includes(format)}
          on:click={() => toggleFormat(format)}
        >
          {format}
        </button>
      {/each}
    </div>
  </div>

  <div class="dialog-actions">
    {#if isEditMode}
      <Button variant="danger" on:click={onDelete}>Delete Source</Button>
      <div style="flex: 1"></div>
    {/if}
    <Button variant="secondary" on:click={onCancel}>Cancel</Button>
    <Button on:click={handleSave} disabled={rootPaths.length === 0}>
      {isEditMode ? 'Update' : 'Save'} Source
    </Button>
  </div>
</div>
{/if}

<style>
  .filesystem-config {
    display: flex;
    flex-direction: column;
    gap: 24px;
  }

  h2 {
    font-size: 24px;
    font-weight: 700;
    color: #2d2d2d;
  }

  .form-section {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .section-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .label {
    display: block;
    font-size: 14px;
    font-weight: 600;
    color: #2d2d2d;
    margin-bottom: 8px;
  }

  .empty-folders {
    padding: 32px;
    text-align: center;
    background: rgba(0, 0, 0, 0.02);
    border-radius: 8px;
    border: 2px dashed rgba(0, 0, 0, 0.1);
  }

  .empty-folders p {
    color: #6b7280;
    font-size: 14px;
  }

  .folder-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .folder-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background: rgba(255, 255, 255, 0.7);
    border-radius: 8px;
    border: 1px solid rgba(0, 0, 0, 0.06);
  }

  .folder-icon {
    font-size: 20px;
    flex-shrink: 0;
  }

  .folder-path {
    flex: 1;
    font-size: 14px;
    color: #2d2d2d;
    font-family: monospace;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .remove-btn {
    width: 24px;
    height: 24px;
    border: none;
    background: rgba(239, 68, 68, 0.1);
    color: #dc2626;
    border-radius: 4px;
    cursor: pointer;
    flex-shrink: 0;
    transition: all 0.15s ease;
  }

  .remove-btn:hover {
    background: rgba(239, 68, 68, 0.2);
  }

  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: #2d2d2d;
    cursor: pointer;
  }

  .checkbox-label input[type="checkbox"] {
    width: 18px;
    height: 18px;
    cursor: pointer;
  }

  .format-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
    gap: 8px;
  }

  .format-btn {
    padding: 8px 12px;
    background: rgba(255, 255, 255, 0.7);
    border: 2px solid rgba(0, 0, 0, 0.1);
    border-radius: 6px;
    font-size: 13px;
    font-weight: 500;
    color: #6b7280;
    cursor: pointer;
    transition: all 0.15s ease;
  }

  .format-btn:hover {
    background: rgba(255, 255, 255, 0.9);
    border-color: rgba(0, 0, 0, 0.2);
  }

  .format-btn.active {
    background: linear-gradient(135deg, #8a65ff, #ff6b9d);
    border-color: transparent;
    color: white;
  }

  .dialog-actions {
    display: flex;
    gap: 12px;
    justify-content: flex-end;
    padding-top: 8px;
  }
  .loading {
    text-align: center;
    padding: 60px;
    color: #6b7280;
    font-size: 14px;
  }
</style>