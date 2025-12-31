export namespace dto {
	
	export class FileNodeDTO {
	    name: string;
	    path: string;
	    isDirectory: boolean;
	    size?: number;
	    extension?: string;
	
	    static createFrom(source: any = {}) {
	        return new FileNodeDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.isDirectory = source["isDirectory"];
	        this.size = source["size"];
	        this.extension = source["extension"];
	    }
	}
	export class DirectoryContentsDTO {
	    currentPath: string;
	    files: FileNodeDTO[];
	    directories: FileNodeDTO[];
	
	    static createFrom(source: any = {}) {
	        return new DirectoryContentsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentPath = source["currentPath"];
	        this.files = this.convertValues(source["files"], FileNodeDTO);
	        this.directories = this.convertValues(source["directories"], FileNodeDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ScanProgressDTO {
	    isScanning: boolean;
	    totalFiles: number;
	    processedFiles: number;
	    currentFile: string;
	    errors?: string[];
	
	    static createFrom(source: any = {}) {
	        return new ScanProgressDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.isScanning = source["isScanning"];
	        this.totalFiles = source["totalFiles"];
	        this.processedFiles = source["processedFiles"];
	        this.currentFile = source["currentFile"];
	        this.errors = source["errors"];
	    }
	}
	export class SourceDTO {
	    id: string;
	    name: string;
	    type: string;
	
	    static createFrom(source: any = {}) {
	        return new SourceDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	    }
	}
	export class TrackDTO {
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
	    duration: number;
	    filePath?: string;
	    streamUrl?: string;
	    format?: string;
	    bitRate?: number;
	    sampleRate?: number;
	    artworkPath?: string;
	
	    static createFrom(source: any = {}) {
	        return new TrackDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.sourceId = source["sourceId"];
	        this.sourceType = source["sourceType"];
	        this.title = source["title"];
	        this.artist = source["artist"];
	        this.artistId = source["artistId"];
	        this.album = source["album"];
	        this.albumId = source["albumId"];
	        this.albumArtist = source["albumArtist"];
	        this.genre = source["genre"];
	        this.year = source["year"];
	        this.trackNumber = source["trackNumber"];
	        this.discNumber = source["discNumber"];
	        this.duration = source["duration"];
	        this.filePath = source["filePath"];
	        this.streamUrl = source["streamUrl"];
	        this.format = source["format"];
	        this.bitRate = source["bitRate"];
	        this.sampleRate = source["sampleRate"];
	        this.artworkPath = source["artworkPath"];
	    }
	}

}

export namespace model {
	
	export class SourceConfiguration {
	    id: string;
	    name: string;
	    type: string;
	    enabled: boolean;
	    config: Record<string, any>;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new SourceConfiguration(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.enabled = source["enabled"];
	        this.config = source["config"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

