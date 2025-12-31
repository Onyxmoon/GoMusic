import { writable } from 'svelte/store';

export type View = 'library' | 'albums' | 'artists' | 'playlists' | 'settings';

// Current active view
export const currentView = writable<View>('library');

// Sidebar collapsed state
export const sidebarCollapsed = writable<boolean>(false);

// Navigation actions
export const navigationActions = {
  navigateTo: (view: View) => {
    currentView.set(view);
  },

  toggleSidebar: () => {
    sidebarCollapsed.update(collapsed => !collapsed);
  }
};