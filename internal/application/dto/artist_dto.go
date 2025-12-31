package dto

// ArtistDTO is the data transfer object for artists exposed to the frontend
type ArtistDTO struct {
	ID         string `json:"id"`
	SourceID   string `json:"sourceId"`
	SourceType string `json:"sourceType"`
	Name       string `json:"name"`
	ImagePath  string `json:"imagePath,omitempty"`
	AlbumCount int    `json:"albumCount"`
	TrackCount int    `json:"trackCount"`
}