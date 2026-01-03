<script lang="ts">
  import TopNav from './lib/components/TopNav.svelte';
  import Sidebar from './lib/components/Sidebar.svelte';
  import BottomPlayer from './lib/components/BottomPlayer.svelte';
  import AudioPlayer from './lib/components/AudioPlayer.svelte';
  import Library from './lib/features/LibraryNew.svelte';
  import Sources from './lib/features/Sources.svelte';
  import FileBrowser from './lib/features/FileBrowser.svelte';

  let activeTab = 'library';
  let activeSection = 'tracks';

  $: currentView = getView(activeTab, activeSection);

  // Reset section when tab changes
  $: if (activeTab) {
    activeSection = getDefaultSection(activeTab);
  }

  function getDefaultSection(tab: string): string {
    switch (tab) {
      case 'library':
        return 'tracks';
      case 'files':
        return 'browser';
      case 'sources':
        return 'all';
      default:
        return 'tracks';
    }
  }

  function getView(tab: string, section: string) {
    switch (tab) {
      case 'library':
        return Library;
      case 'sources':
        return Sources;
      case 'files':
        return FileBrowser;
      default:
        return Library;
    }
  }

  function handleSectionChange(event: CustomEvent) {
    activeSection = event.detail.section;
  }

  function handleSidebarAction(event: CustomEvent) {
    const { type } = event.detail;
    if (type === 'new-playlist') {
      console.log('Create new playlist');
      // TODO: Implement playlist creation
    } else if (type === 'add-source') {
      console.log('Add new source');
      // TODO: Trigger add source dialog
    }
  }
</script>

<div class="app">
  <div class="app-layout">
    <TopNav bind:activeTab />

    <div class="app-body">
      <Sidebar
        {activeTab}
        bind:activeSection
        on:sectionchange={handleSectionChange}
        on:action={handleSidebarAction}
      />

      <main class="main-content">
        {#if activeTab === 'library'}
          <Library section={activeSection} />
        {:else if activeTab === 'sources'}
          <Sources />
        {:else if activeTab === 'files'}
          <FileBrowser />
        {/if}
      </main>
    </div>
  </div>

  <BottomPlayer />

  <!-- Audio Player (always running in background) -->
  <AudioPlayer />
</div>

<style>
  /* Global Scrollbar Styles */
  :global(::-webkit-scrollbar) {
    width: 8px;
    height: 8px;
  }

  :global(::-webkit-scrollbar-track) {
    background: transparent;
  }

  :global(::-webkit-scrollbar-thumb) {
    background-color: transparent;
    border-radius: 4px;
    transition: background-color 0.2s ease;
  }

  :global(div:hover::-webkit-scrollbar-thumb) {
    background-color: rgba(138, 101, 255, 0.5);
  }

  :global(::-webkit-scrollbar-thumb:hover) {
    background-color: rgba(138, 101, 255, 0.8);
  }

  .app {
    width: 100vw;
    height: 100vh;
    overflow: hidden;
  }

  .app-layout {
    height: 100vh;
    display: flex;
    flex-direction: column;
  }

  .app-body {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .main-content {
    flex: 1;
    overflow: hidden;
    padding: 24px;
    display: flex;
    flex-direction: column;
  }
</style>
