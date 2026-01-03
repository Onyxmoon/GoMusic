import type { dto } from '../../../wailsjs/go/models';

/**
 * Manages audio playback state and playlist queue
 */
class PlayerStore {
  currentTrack = $state<dto.TrackDTO | null>(null);
  isPlaying = $state<boolean>(false);
  currentTime = $state<number>(0);
  volume = $state<number>(0.8);
  isMuted = $state<boolean>(false);
  playlist = $state<dto.TrackDTO[]>([]);
  currentIndex = $state<number>(-1);

  // Additional features
  repeatMode = $state<'none' | 'all' | 'one'>('none');
  shuffleEnabled = $state<boolean>(false);

  duration = $derived(this.currentTrack?.duration || 0);
  progress = $derived(
    this.duration > 0 ? (this.currentTime / this.duration) * 100 : 0
  );

  // Logarithmic volume curve
  actualVolume = $derived.by(() => {
    if (this.volume === 0) return 0;

    const minDb = -60;
    const maxDb = 0;
    const db = minDb + (maxDb - minDb) * this.volume;

    // Convert dB to linear gain: gain = 10^(dB/20)
    return Math.pow(10, db / 20);
  });

  // Check if we can go to next/previous track
  canGoNext = $derived(this.currentIndex < this.playlist.length - 1);
  canGoPrevious = $derived(this.currentIndex > 0);

  /**
   * Play a specific track
   * This sets the track as current and starts playback
   */
  play(track: dto.TrackDTO) {
    if (!track) return;

    this.currentTrack = track;
    this.isPlaying = true;

    if (this.previewMode && track.duration) {
      this.currentTime = track.duration / 15;
    } else {
      this.currentTime = 0;
    }
  }

  /**
   * Pause playback
   */
  pause() {
    this.isPlaying = false;
  }

  /**
   * Resume playback
   */
  resume() {
    if (!this.currentTrack) return;
    this.isPlaying = true;
  }

  /**
   * Stop playback and reset time
   */
  stop() {
    this.isPlaying = false;
    this.currentTime = 0;
  }

  /**
   * Seek to a specific time (in seconds)
   */
  seek(time: number) {
    if (!this.currentTrack) return;

    // Bounds checking
    this.currentTime = Math.max(0, Math.min(time, this.duration));
  }

  /**
   * Set volume (0.0 to 1.0)
   */
  setVolume(vol: number) {
    // Clamp between 0 and 1
    this.volume = Math.max(0, Math.min(1, vol));
  }

  /**
   * Toggle mute state
   */
  toggleMute() {
    this.isMuted = !this.isMuted;
  }

  /**
   * Play the next track in the playlist
   */
  next() {
    // Bounds checking
    if (this.currentIndex >= this.playlist.length - 1) {
      // Handle repeat mode
      if (this.repeatMode === 'all') {
        this.currentIndex = 0;
        this.currentTrack = this.playlist[0];
        this.currentTime = 0;
      }
      return;
    }

    // Atomically update index and track
    this.currentIndex++;
    const nextTrack = this.playlist[this.currentIndex];

    if (nextTrack) {
      this.currentTrack = nextTrack;
      this.currentTime = 0;
      // Keep playing if already playing
      // isPlaying state is preserved
    } else {
      console.error('PlayerStore.next(): track at index not found');
    }
  }

  /**
   * Play the previous track in the playlist
   * FIXES THE BROKEN IMPLEMENTATION
   */
  previous() {
    // If we're more than 3 seconds into the track, restart it
    if (this.currentTime > 3) {
      this.currentTime = 0;
      return;
    }

    // Bounds checking
    if (this.currentIndex <= 0) {
      // Handle repeat mode
      if (this.repeatMode === 'all') {
        this.currentIndex = this.playlist.length - 1;
        this.currentTrack = this.playlist[this.currentIndex];
        this.currentTime = 0;
      }
      return;
    }

    // Atomically update index and track
    this.currentIndex--;
    const prevTrack = this.playlist[this.currentIndex];

    if (prevTrack) {
      this.currentTrack = prevTrack;
      this.currentTime = 0;
      // Keep playing if already playing
    } else {
      console.error('PlayerStore.previous(): track at index not found');
    }
  }

  /**
   * Set the current playlist and optionally start playing
   */
  setPlaylist(tracks: dto.TrackDTO[], startIndex: number = 0) {
    if (!tracks || tracks.length === 0) {
      this.playlist = [];
      this.currentIndex = -1;
      return;
    }

    // Validate startIndex
    const validIndex = Math.max(0, Math.min(startIndex, tracks.length - 1));

    this.playlist = tracks;
    this.currentIndex = validIndex;
    this.currentTrack = tracks[validIndex];
  }

  /**
   * Update current time (called by audio element)
   */
  updateTime(time: number) {
    // Only update if valid number
    if (!isNaN(time) && isFinite(time)) {
      this.currentTime = time;
    }
  }

  /**
   * Toggle repeat mode: none → all → one → none
   */
  toggleRepeat() {
    if (this.repeatMode === 'none') {
      this.repeatMode = 'all';
    } else if (this.repeatMode === 'all') {
      this.repeatMode = 'one';
    } else {
      this.repeatMode = 'none';
    }
  }

  /**
   * Toggle shuffle
   */
  toggleShuffle() {
    this.shuffleEnabled = !this.shuffleEnabled;
    // TODO: Implement shuffle logic (create shuffled playlist)
  }

  /**
   * Handle track end - automatically advance or repeat
   */
  handleTrackEnd() {
    if (this.repeatMode === 'one') {
      // Repeat current track
      this.currentTime = 0;
      // Keep playing
    } else if (this.canGoNext) {
      // Go to next track
      this.next();
    } else if (this.repeatMode === 'all') {
      // Loop back to first track
      this.currentIndex = 0;
      this.currentTrack = this.playlist[0];
      this.currentTime = 0;
    } else {
      // End of playlist, stop playing
      this.isPlaying = false;
      this.currentTime = 0;
    }
  }

  /**
   * Check if a specific track is currently playing
   * @param trackId - The ID of the track to check
   * @returns Object with isCurrentTrack and isPlaying flags
   */
  isTrackPlaying(trackId: string): { isCurrentTrack: boolean; isPlaying: boolean } {
    const isCurrentTrack = this.currentTrack?.id === trackId;
    return {
      isCurrentTrack,
      isPlaying: isCurrentTrack && this.isPlaying
    };
  }
}

// Export a single instance
export const player = new PlayerStore();