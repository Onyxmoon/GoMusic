<script lang="ts">
  import { onMount } from 'svelte';
  import { BrowseDirectory, GetSourceRootPath, GetSources } from '../../../wailsjs/go/main/App.js';
  import { Folder, Music } from 'lucide-svelte';

  export let sourceId: string = '';

  interface FileNode {
    name: string;
    path: string;
    isDirectory: boolean;
    size?: number;
    extension?: string;
  }

  interface DirectoryContents {
    currentPath: string;
    files: FileNode[];
    directories: FileNode[];
  }

  let currentPath = '';
  let files: FileNode[] = [];
  let directories: FileNode[] = [];
  let breadcrumbs: string[] = [];
  let rootPath: string = '';
  let isLoading = false;
  let error: string | null = null;
  let availableSources: Array<{id: string, name: string, type: string}> = [];

  onMount(async () => {
    try {
      // Load available sources
      const sources = await GetSources();
      console.log('Available sources:', sources);
      availableSources = sources;

      // If no valid sourceId, use first available
      if (!sourceId && sources && sources.length > 0) {
        sourceId = sources[0].id;
        console.log('Using first source:', sourceId);
      }

      if (!sourceId) {
        error = 'No source ID available';
        return;
      }

      // Get root path for the source
      rootPath = await GetSourceRootPath(sourceId);
      await loadDirectory('');
    } catch (err) {
      error = `Failed to initialize: ${err}`;
      console.error(error);
    }
  });

  async function loadDirectory(path: string) {
    try {
      isLoading = true;
      error = null;

      const result: DirectoryContents = await BrowseDirectory(sourceId, path);

      currentPath = result.currentPath;
      files = result.files || [];
      directories = result.directories || [];

      // Update breadcrumbs
      if (!path || path === '' || path === '/') {
        breadcrumbs = ['Root'];
      } else {
        const parts = path.split('/').filter(Boolean);
        breadcrumbs = ['Root', ...parts];
      }
    } catch (err) {
      error = `Failed to load directory: ${err}`;
      console.error(error);
    } finally {
      isLoading = false;
    }
  }

  function navigateToFolder(node: FileNode) {
    if (node.isDirectory) {
      // Calculate relative path from root
      let relativePath = node.path;
      if (rootPath && node.path.startsWith(rootPath)) {
        relativePath = node.path.substring(rootPath.length);
      }
      loadDirectory(relativePath);
    }
  }

  function navigateToBreadcrumb(index: number) {
    if (index === 0) {
      // Navigate to root
      loadDirectory('');
    } else {
      // Build path from breadcrumbs
      const pathParts = breadcrumbs.slice(1, index + 1);
      const newPath = '/' + pathParts.join('/');
      loadDirectory(newPath);
    }
  }

  async function changeSource(newSourceId: string) {
    sourceId = newSourceId;
    try {
      rootPath = await GetSourceRootPath(sourceId);
      await loadDirectory('');
    } catch (err) {
      error = `Failed to change source: ${err}`;
      console.error(error);
    }
  }
</script>

<div class="file-browser">
  <div class="browser-header">
    <h2>File Browser</h2>
    <div class="breadcrumbs">
      {#each breadcrumbs as crumb, index}
        <button
          class="breadcrumb"
          class:active={index === breadcrumbs.length - 1}
          on:click={() => navigateToBreadcrumb(index)}
        >
          {crumb}
        </button>
        {#if index < breadcrumbs.length - 1}
          <span class="separator">/</span>
        {/if}
      {/each}
    </div>
  </div>

  <div class="file-list">
    <div class="file-list-header">
      <div class="col col-name">Name</div>
      <div class="col col-type">Type</div>
      <div class="col col-path">Path</div>
    </div>

    {#if isLoading}
      <div class="empty-state">
        <div class="spinner"></div>
        <p>Loading directory...</p>
      </div>
    {:else if error}
      <div class="empty-state error">
        <p>Error: {error}</p>
      </div>
    {:else if directories.length === 0 && files.length === 0}
      <div class="empty-state">
        <p>No files or folders found</p>
        <p class="hint">Configure a source in the Sources view to scan music files</p>
      </div>
    {:else}
      {#each directories as dir}
        <button
          class="file-row is-directory"
          on:click={() => navigateToFolder(dir)}
        >
          <div class="col col-name">
            <span class="file-icon"><Folder size={20} /></span>
            <span class="file-name">{dir.name}</span>
          </div>
          <div class="col col-type">Folder</div>
          <div class="col col-path">{dir.path}</div>
        </button>
      {/each}
      {#each files as file}
        <button
          class="file-row"
          on:click={() => navigateToFolder(file)}
        >
          <div class="col col-name">
            <span class="file-icon"><Music size={20} /></span>
            <span class="file-name">{file.name}</span>
          </div>
          <div class="col col-type">
            {file.extension || 'Audio File'}
          </div>
          <div class="col col-path">{file.path}</div>
        </button>
      {/each}
    {/if}
  </div>
</div>

<style>
  .file-browser {
    padding: 24px;
    max-width: 1400px;
  }

  .browser-header {
    margin-bottom: 24px;
  }

  .browser-header h2 {
    font-size: 32px;
    font-weight: 700;
    color: #2d2d2d;
    margin-bottom: 16px;
  }

  .breadcrumbs {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 12px 16px;
    background: rgba(255, 255, 255, 0.7);
    border-radius: 8px;
    overflow-x: auto;
  }

  .breadcrumb {
    padding: 4px 12px;
    background: transparent;
    border: none;
    border-radius: 4px;
    color: #6b7280;
    font-size: 14px;
    cursor: pointer;
    white-space: nowrap;
    transition: all 0.15s ease;
  }

  .breadcrumb:hover {
    background: rgba(0, 0, 0, 0.05);
    color: #2d2d2d;
  }

  .breadcrumb.active {
    color: #2d2d2d;
    font-weight: 600;
  }

  .separator {
    color: #9ca3af;
    font-size: 14px;
  }

  .file-list {
    background: rgba(255, 255, 255, 0.5);
    border-radius: 12px;
    padding: 16px;
  }

  .file-list-header {
    display: grid;
    grid-template-columns: 2fr 150px 1fr;
    gap: 16px;
    padding: 12px 16px;
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
    font-size: 12px;
    font-weight: 600;
    color: #6b7280;
    text-transform: uppercase;
  }

  .file-row {
    display: grid;
    grid-template-columns: 2fr 150px 1fr;
    gap: 16px;
    padding: 12px 16px;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    text-align: left;
    transition: all 0.15s ease;
    width: 100%;
  }

  .file-row:hover {
    background: rgba(255, 255, 255, 0.7);
  }

  .file-row.is-directory:hover {
    background: rgba(138, 101, 255, 0.1);
  }

  .col {
    display: flex;
    align-items: center;
    font-size: 14px;
    color: #2d2d2d;
  }

  .col-name {
    gap: 12px;
  }

  .file-icon {
    font-size: 20px;
    flex-shrink: 0;
  }

  .file-name {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .col-type,
  .col-path {
    color: #6b7280;
  }

  .col-path {
    font-family: monospace;
    font-size: 12px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .empty-state {
    padding: 80px 20px;
    text-align: center;
  }

  .empty-state p {
    color: #6b7280;
    margin-bottom: 8px;
  }

  .empty-state .hint {
    font-size: 13px;
    color: #9ca3af;
  }

  .empty-state.error p {
    color: #dc2626;
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
</style>