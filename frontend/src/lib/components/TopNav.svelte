<script lang="ts">
  import { Search, Ellipsis, Bell } from 'lucide-svelte';

  export let activeTab: string = 'playlists';

  const tabs = [
    { id: 'library', label: 'Library' },
    { id: 'files', label: 'Files' },
    { id: 'sources', label: 'Sources' },
    { id: 'settings', label: 'Settings' }
  ];

  let searchQuery = '';
</script>

<nav class="top-nav">
  <div class="nav-left">
    <div class="app-icon">

    </div>
    <div class="nav-tabs">
      {#each tabs as tab}
        <button
          class="nav-tab"
          class:active={activeTab === tab.id}
          on:click={() => activeTab = tab.id}
        >
          {tab.label}
        </button>
      {/each}
    </div>
  </div>

  <div class="nav-center">
    <div class="search-bar">
      <span class="search-icon"><Search size={16} /></span>
      <input
        type="text"
        placeholder="What do you want to listen?"
        bind:value={searchQuery}
      />
    </div>
  </div>

  <div class="nav-right">
    <button class="icon-btn"><Ellipsis size={18} /></button>
    <button class="icon-btn"><Bell size={18} /></button>
    <div class="user-avatar">
      <img src="https://placehold.co/32" alt="User" />
    </div>
  </div>
</nav>

<style>
  .top-nav {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 20px;
    background: rgba(255, 255, 255, 0.6);
    backdrop-filter: blur(40px) saturate(180%);
    -webkit-backdrop-filter: blur(40px) saturate(180%);
    border-bottom: 1px solid rgba(255, 255, 255, 0.3);
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
    z-index: 1;
  }

  .nav-left {
    display: flex;
    align-items: center;
    gap: 16px;
    flex: 1;
  }

  .app-icon {
    width: 40px;
    height: 40px;
    background: linear-gradient(135deg, #8a65ff, #ff6b9d);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .nav-tabs {
    display: flex;
    gap: 4px;
  }

  .nav-tab {
    padding: 6px 16px;
    background: transparent;
    border: none;
    border-radius: 6px;
    color: #6b7280;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .nav-tab:hover {
    background: rgba(0, 0, 0, 0.05);
    color: #2d2d2d;
  }

  .nav-tab.active {
    background: rgba(255, 255, 255, 0.9);
    color: #2d2d2d;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }

  .nav-center {
    flex: 1;
    max-width: 500px;
  }

  .search-bar {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    background: rgba(255, 255, 255, 0.7);
    border-radius: 8px;
    border: 1px solid rgba(0, 0, 0, 0.06);
  }

  .search-icon {
    color: #9ca3af;
    font-size: 14px;
  }

  .search-bar input {
    flex: 1;
    border: none;
    background: transparent;
    outline: none;
    color: #2d2d2d;
    font-size: 14px;
  }

  .search-bar input::placeholder {
    color: #9ca3af;
  }

  .nav-right {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
    justify-content: flex-end;
  }

  .icon-btn {
    width: 36px;
    height: 36px;
    border-radius: 8px;
    border: none;
    background: rgba(255, 255, 255, 0.7);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 16px;
    transition: all 0.2s ease;
  }

  .icon-btn:hover {
    background: rgba(255, 255, 255, 0.9);
  }

  .user-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    overflow: hidden;
  }

  .user-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
</style>
