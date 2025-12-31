package dto

// AlbumDTO is the data transfer object for albums exposed to the frontend
type AlbumDTO struct {
	ID            string  `json:"id"`
	SourceID      string  `json:"sourceId"`
	SourceType    string  `json:"sourceType"`
	Title         string  `json:"title"`
	Artist        string  `json:"artist"`
	ArtistID      string  `json:"artistId"`
	Year          int     `json:"year,omitempty"`
	Genre         string  `json:"genre,omitempty"`
	ArtworkPath   string  `json:"artworkPath,omitempty"`
	TrackCount    int     `json:"trackCount"`
	TotalDuration float64 `json:"totalDuration"` // Duration in seconds
}