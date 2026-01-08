package dto

// TrackDTO is the data transfer object for tracks exposed to the frontend
type TrackDTO struct {
	ID          string  `json:"id"`
	SourceID    string  `json:"sourceId"`
	SourceType  string  `json:"sourceType"`
	Title       string  `json:"title"`
	Artist      string  `json:"artist"`
	ArtistID    string  `json:"artistId"`
	Album       string  `json:"album"`
	AlbumID     string  `json:"albumId"`
	AlbumArtist string  `json:"albumArtist,omitempty"`
	Genre       string  `json:"genre,omitempty"`
	Year        int     `json:"year,omitempty"`
	TrackNumber int     `json:"trackNumber,omitempty"`
	DiscNumber  int     `json:"discNumber,omitempty"`
	Duration    float64 `json:"duration"` // Duration in seconds
	FilePath    string  `json:"filePath,omitempty"`
	StreamURL   string  `json:"streamUrl,omitempty"`
	Format      string  `json:"format,omitempty"`
	BitRate     int     `json:"bitRate,omitempty"`
	SampleRate  int     `json:"sampleRate,omitempty"`
	HasArtwork  bool    `json:"hasArtwork"`
}