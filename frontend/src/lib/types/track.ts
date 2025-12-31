// TypeScript types matching the Go DTOs

export interface TrackDTO {
  id: string;
  sourceId: string;
  sourceType: string;
  title: string;
  artist: string;
  artistId: string;
  album: string;
  albumId: string;
  albumArtist?: string;
  genre?: string;
  year?: number;
  trackNumber?: number;
  discNumber?: number;
  duration: number; // in seconds
  filePath?: string;
  streamUrl?: string;
  format?: string;
  bitRate?: number;
  sampleRate?: number;
  artworkPath?: string;
}

export interface AlbumDTO {
  id: string;
  sourceId: string;
  sourceType: string;
  title: string;
  artist: string;
  artistId: string;
  year?: number;
  genre?: string;
  artworkPath?: string;
  trackCount: number;
  totalDuration: number;
}

export interface ArtistDTO {
  id: string;
  sourceId: string;
  sourceType: string;
  name: string;
  imagePath?: string;
  albumCount: number;
  trackCount: number;
}

export interface PlaylistDTO {
  id: string;
  name: string;
  description?: string;
  trackIds: string[];
  trackCount: number;
  coverPath?: string;
}

export interface SourceDTO {
  id: string;
  name: string;
  type: string;
}

export interface ScanProgressDTO {
  isScanning: boolean;
  totalFiles: number;
  processedFiles: number;
  currentFile: string;
  errors?: string[];
}