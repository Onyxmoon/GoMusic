<script lang="ts">
  import TopNav from './lib/components/TopNav.svelte';
  import Sidebar from './lib/components/Sidebar.svelte';
  import BottomPlayer from './lib/components/BottomPlayer.svelte';
  import AudioPlayer from './lib/components/AudioPlayer.svelte';
  import Library from './lib/features/LibraryNew.svelte';
  import Sources from './lib/features/Sources.svelte';
  import FileBrowser from './lib/features/FileBrowser.svelte';
  import { player } from './lib/stores/player.svelte';
  import { keyboardShortcuts } from './lib/stores/keyboardShortcuts.svelte';

  let activeTab = $state('library');
  let activeSection = $state('tracks');

  function handleSectionChange(event: CustomEvent) {
    activeSection = event.detail.section;
  }

  function handleSidebarAction(event: CustomEvent) {
    const { type } = event.detail;
    if (type === 'new-playlist') {
      console.log('Create new playlist');
    } else if (type === 'add-source') {
      console.log('Add new source');
    }
  }

  $effect(() => {
      const unregister = keyboardShortcuts.register({
          viewId: 'global-player',
          priority: 50, // Low priority - only when no view handles it
          shortcuts: [
              {
                  key: ' ',
                  ctrl: false,
                  handler: (e) => {
                      e.preventDefault();
                      if (player.currentTrack) {
                          player.isPlaying ? player.pause() : player.resume();
                      }
                      return true;
                  },
                  description: 'Play/Pause current track'
              },
              {
                  key: 'ArrowRight',
                  ctrl: true,
                  handler: (e) => {
                      e.preventDefault();
                      if (player.canGoNext) player.next();
                      return true;
                  },
                  description: 'Next track'
              },
              {
                  key: 'ArrowLeft',
                  ctrl: true,
                  handler: (e) => {
                      e.preventDefault();
                      if (player.canGoPrevious) player.previous();
                      return true;
                  },
                  description: 'Previous track'
              }
          ]
      });

      return unregister;
  });

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

  :global(*:hover::-webkit-scrollbar-thumb) {
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
