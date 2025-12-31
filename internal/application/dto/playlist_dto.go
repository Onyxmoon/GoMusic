package dto

// PlaylistDTO is the data transfer object for playlists exposed to the frontend
type PlaylistDTO struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	TrackIDs    []string `json:"trackIds"`
	TrackCount  int      `json:"trackCount"`
	CoverPath   string   `json:"coverPath,omitempty"`
}