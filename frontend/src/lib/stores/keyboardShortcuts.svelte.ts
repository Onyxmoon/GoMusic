// ============================================================================
// Types
// ============================================================================

export interface ShortcutHandler {
  key: string;              // e.g. 'ArrowDown', ' ', 'Enter'
  ctrl?: boolean;           // Optional: Ctrl/Cmd modifier
  shift?: boolean;          // Optional: Shift modifier
  alt?: boolean;            // Optional: Alt modifier
  handler: (e: KeyboardEvent) => boolean | void;  // true = handled (stop propagation)
  description?: string;     // For debugging/help screen
}

export interface ViewShortcuts {
  viewId: string;           // Unique ID (e.g. 'library-tracktable', 'player')
  priority: number;         // Higher = more important (1000 = Critical, 100 = Views, 10 = Global)
  condition?: () => boolean; // Optional: only active when true
  shortcuts: ShortcutHandler[];
}

// ============================================================================
// KeyboardShortcutManager
// ============================================================================

class KeyboardShortcutManager {
  private registrations: ViewShortcuts[] = [];

  constructor() {
    // Global listener setup (only in browser)
    if (typeof window !== 'undefined') {
      window.addEventListener('keydown', this.handleKeydown.bind(this));
    }
  }

  // ============================================================================
  // Public API
  // ============================================================================

  /**
   * Register keyboard shortcuts for a view
   * @param config View shortcuts configuration
   * @returns Unregister function for cleanup
   */
  register(config: ViewShortcuts): () => void {
    // Check for conflicts
    this.detectConflicts(config);

    this.registrations.push(config);
    // Sort by priority (highest first)
    this.registrations.sort((a, b) => b.priority - a.priority);

    // Return unregister function
    return () => this.unregister(config.viewId);
  }

  /**
   * Unregister shortcuts for a view
   * @param viewId View identifier
   */
  unregister(viewId: string): void {
    this.registrations = this.registrations.filter(r => r.viewId !== viewId);
  }

  /**
   * Get currently active shortcuts (condition === true)
   * @returns Array of active shortcut handlers
   */
  getActiveShortcuts(): Array<{ viewId: string; shortcuts: ShortcutHandler[] }> {
    return this.registrations
      .filter(r => !r.condition || r.condition())
      .map(r => ({
        viewId: r.viewId,
        shortcuts: r.shortcuts
      }));
  }

  /**
   * Get all registered views
   * @returns Array of all registrations
   */
  getAllRegistrations(): ViewShortcuts[] {
    return [...this.registrations];
  }

  // ============================================================================
  // Private Methods
  // ============================================================================

  /**
   * Handle global keydown event
   * @param e Keyboard event
   */
  private handleKeydown(e: KeyboardEvent): void {
    // Don't intercept if user is typing in input
    if (this.isInputFocused()) return;

    // Iterate through registrations (sorted by priority)
    for (const registration of this.registrations) {
      // Check condition
      if (registration.condition && !registration.condition()) {
        continue;
      }

      // Check each shortcut
      for (const shortcut of registration.shortcuts) {
        if (this.matchesShortcut(e, shortcut)) {
          const handled = shortcut.handler(e);
          if (handled) {
            // Shortcut was handled - stop processing
            return;
          }
        }
      }
    }
  }

  /**
   * Check if keyboard event matches a shortcut configuration
   * @param e Keyboard event
   * @param shortcut Shortcut handler configuration
   * @returns True if event matches shortcut
   */
  private matchesShortcut(e: KeyboardEvent, shortcut: ShortcutHandler): boolean {
    if (e.key !== shortcut.key) return false;

    // Check modifiers
    if (shortcut.ctrl && !(e.ctrlKey || e.metaKey)) return false;
    if (!shortcut.ctrl && (e.ctrlKey || e.metaKey)) return false;

    if (shortcut.shift && !e.shiftKey) return false;
    if (!shortcut.shift && e.shiftKey) return false;

    if (shortcut.alt && !e.altKey) return false;
    if (!shortcut.alt && e.altKey) return false;

    return true;
  }

  /**
   * Check if an input element currently has focus
   * @returns True if input/textarea/contenteditable is focused
   */
  private isInputFocused(): boolean {
    const active = document.activeElement;
    return active instanceof HTMLInputElement ||
           active instanceof HTMLTextAreaElement ||
           active?.getAttribute('contenteditable') === 'true';
  }

  /**
   * Detect conflicts when registering shortcuts
   * Logs warning if same shortcut is registered with same priority
   * @param config New view shortcuts configuration
   */
  private detectConflicts(config: ViewShortcuts): void {
    for (const existingReg of this.registrations) {
      // Only check if same priority
      if (existingReg.priority !== config.priority) continue;
      if (existingReg.viewId === config.viewId) continue;

      // Check for overlapping shortcuts
      for (const newShortcut of config.shortcuts) {
        for (const existingShortcut of existingReg.shortcuts) {
          if (this.shortcutsMatch(newShortcut, existingShortcut)) {
            const modifiers = this.getModifierString(newShortcut);
            const keyCombo = modifiers ? `${modifiers}+${newShortcut.key}` : newShortcut.key;

            console.warn(
              `[KeyboardShortcuts] Conflict detected: '${keyCombo}' registered by ` +
              `'${existingReg.viewId}' and '${config.viewId}' with same priority ${config.priority}`
            );
          }
        }
      }
    }
  }

  /**
   * Check if two shortcuts match (same key + modifiers)
   * @param a First shortcut
   * @param b Second shortcut
   * @returns True if shortcuts match
   */
  private shortcutsMatch(a: ShortcutHandler, b: ShortcutHandler): boolean {
    return a.key === b.key &&
           !!a.ctrl === !!b.ctrl &&
           !!a.shift === !!b.shift &&
           !!a.alt === !!b.alt;
  }

  /**
   * Get modifier string for display (e.g. 'Ctrl', 'Shift+Ctrl')
   * @param shortcut Shortcut handler
   * @returns Modifier string
   */
  private getModifierString(shortcut: ShortcutHandler): string {
    const mods: string[] = [];
    if (shortcut.ctrl) mods.push('Ctrl');
    if (shortcut.shift) mods.push('Shift');
    if (shortcut.alt) mods.push('Alt');
    return mods.join('+');
  }
}

// ============================================================================
// Export Singleton
// ============================================================================

export const keyboardShortcuts = new KeyboardShortcutManager();
