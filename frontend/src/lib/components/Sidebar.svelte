<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { Music, Disc, Mic, List, Folder, Clock, Database, FolderOpen, Globe, Plus } from 'lucide-svelte';

  export let activeTab: string = 'library';
  export let activeSection: string = 'tracks';

  const dispatch = createEventDispatcher();

  // Library sections
  const librarySections = [
    { id: 'tracks', label: 'Titles', icon: Music },
    { id: 'albums', label: 'Albums', icon: Disc },
    { id: 'artists', label: 'Artists', icon: Mic },
    { id: 'playlists', label: 'Playlists', icon: List },
  ];

  // Files sections
  const filesSections = [
    { id: 'browser', label: 'File Browser', icon: Folder },
    { id: 'recent', label: 'Recent Folders', icon: Clock },
  ];

  // Sources sections
  const sourcesSections = [
    { id: 'all', label: 'All Sources', icon: Database },
    { id: 'filesystem', label: 'Filesystem', icon: FolderOpen },
    { id: 'api', label: 'API Sources', icon: Globe },
  ];

  $: currentSections = getSections(activeTab);

  function getSections(tab: string) {
    switch (tab) {
      case 'library':
        return librarySections;
      case 'files':
        return filesSections;
      case 'sources':
        return sourcesSections;
      default:
        return [];
    }
  }

  function handleSectionClick(sectionId: string) {
    activeSection = sectionId;
    dispatch('sectionchange', { section: sectionId });
  }

  function handleAction() {
    if (activeTab === 'library') {
      dispatch('action', { type: 'new-playlist' });
    } else if (activeTab === 'sources') {
      dispatch('action', { type: 'add-source' });
    }
  }
</script>

<aside class="sidebar">
  <div class="sidebar-header">
    <h3>
      {#if activeTab === 'library'}
        Library
      {:else if activeTab === 'files'}
        Files
      {:else if activeTab === 'sources'}
        Sources
      {:else}
        Settings
      {/if}
    </h3>
  </div>

  <nav class="sidebar-nav">
    {#each currentSections as section}
      <button
        class="nav-item"
        class:active={activeSection === section.id}
        on:click={() => handleSectionClick(section.id)}
      >
        <span class="nav-icon">
          <svelte:component this={section.icon} size={20} />
        </span>
        <span class="nav-label">{section.label}</span>
      </button>
    {/each}
  </nav>

  {#if activeTab === 'library'}
    <div class="sidebar-footer">
      <button class="footer-btn" on:click={handleAction}>
        <Plus size={16} />
        <span>New Playlist</span>
      </button>
    </div>
  {:else if activeTab === 'sources'}
    <div class="sidebar-footer">
      <button class="footer-btn" on:click={handleAction}>
        <Plus size={16} />
        <span>Add Source</span>
      </button>
    </div>
  {/if}
</aside>

<style>
  .sidebar {
    width: 240px;
    height: 100%;
    background: rgba(255, 255, 255, 0.6);
    backdrop-filter: blur(40px) saturate(180%);
    -webkit-backdrop-filter: blur(40px) saturate(180%);
    border-right: 1px solid rgba(255, 255, 255, 0.5);
    box-shadow: 4px 0 24px rgba(0, 0, 0, 0.08);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .sidebar-header {
    padding: 20px 20px 16px;
    border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  }

  .sidebar-header h3 {
    font-size: 18px;
    font-weight: 700;
    color: #2d2d2d;
  }

  .sidebar-nav {
    flex: 1;
    overflow-y: auto;
    padding: 12px;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 12px;
    background: transparent;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.15s ease;
    text-align: left;
    width: 100%;
  }

  .nav-item:hover {
    background: rgba(255, 255, 255, 0.7);
  }

  .nav-item.active {
    background: rgba(138, 101, 255, 0.15);
  }

  .nav-item.active .nav-label {
    font-weight: 600;
    color: #8a65ff;
  }

  .nav-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
  }

  .nav-label {
    font-size: 14px;
    color: #2d2d2d;
    flex: 1;
  }

  .sidebar-footer {
    padding: 12px;
    border-top: 1px solid rgba(0, 0, 0, 0.06);
  }

  .footer-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 10px 16px;
    width: 100%;
    background: rgba(255, 255, 255, 0.7);
    border: 1px solid rgba(0, 0, 0, 0.06);
    border-radius: 8px;
    cursor: pointer;
    font-size: 14px;
    font-weight: 500;
    color: #2d2d2d;
    transition: all 0.15s ease;
  }

  .footer-btn:hover {
    background: rgba(255, 255, 255, 0.9);
    border-color: #8a65ff;
  }
</style>