/**
 * Event types for Wails runtime events
 */

export interface ScanErrorEvent {
	sourceId: string;
	error: string;
}

export interface ScanProgressEvent {
	sourceId: string;
	current: number;
	total: number;
	currentFile?: string;
}

export interface ScanStartedEvent {
	sourceId: string;
}

export interface ScanCompleteEvent {
	sourceId: string;
}