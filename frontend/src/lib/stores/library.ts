import { writable, derived } from 'svelte/store';
import type { TrackDTO, ScanProgressDTO } from '../types/track';

// Library state
export const tracks = writable<TrackDTO[]>([]);
export const isLoading = writable<boolean>(false);
export const error = writable<string | null>(null);
export const scanProgress = writable<ScanProgressDTO | null>(null);

// Search state
export const searchQuery = writable<string>('');

// Filtered tracks based on search query
export const filteredTracks = derived(
  [tracks, searchQuery],
  ([$tracks, $searchQuery]) => {
    if (!$searchQuery) return $tracks;

    const query = $searchQuery.toLowerCase();
    return $tracks.filter(track =>
      track.title.toLowerCase().includes(query) ||
      track.artist.toLowerCase().includes(query) ||
      track.album.toLowerCase().includes(query)
    );
  }
);

// Group tracks by album
export const tracksByAlbum = derived(tracks, ($tracks) => {
  const grouped = new Map<string, TrackDTO[]>();

  $tracks.forEach(track => {
    const albumKey = track.albumId || 'unknown';
    if (!grouped.has(albumKey)) {
      grouped.set(albumKey, []);
    }
    grouped.get(albumKey)!.push(track);
  });

  return grouped;
});

// Group tracks by artist
export const tracksByArtist = derived(tracks, ($tracks) => {
  const grouped = new Map<string, TrackDTO[]>();

  $tracks.forEach(track => {
    const artistKey = track.artistId || 'unknown';
    if (!grouped.has(artistKey)) {
      grouped.set(artistKey, []);
    }
    grouped.get(artistKey)!.push(track);
  });

  return grouped;
});