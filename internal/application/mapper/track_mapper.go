package mapper

import (
	"GoMusic/internal/application/dto"
	"GoMusic/internal/domain/model"
)

// TrackMapper converts domain models to DTOs
type TrackMapper struct{}

// NewTrackMapper creates a new track mapper
func NewTrackMapper() *TrackMapper {
	return &TrackMapper{}
}

// ToDTO converts a Track to TrackDTO
func (m *TrackMapper) ToDTO(track *model.Track) *dto.TrackDTO {
	if track == nil {
		return nil
	}

	return &dto.TrackDTO{
		ID:          track.ID,
		SourceID:    track.SourceID,
		SourceType:  string(track.SourceType),
		Title:       track.Title,
		Artist:      track.Artist,
		ArtistID:    track.ArtistID,
		Album:       track.Album,
		AlbumID:     track.AlbumID,
		AlbumArtist: track.AlbumArtist,
		Genre:       track.Genre,
		Year:        track.Year,
		TrackNumber: track.TrackNumber,
		DiscNumber:  track.DiscNumber,
		Duration:    track.Duration.Seconds(),
		FilePath:    track.FilePath,
		StreamURL:   track.StreamURL,
		Format:      track.Format,
		BitRate:     track.BitRate,
		SampleRate:  track.SampleRate,
		HasArtwork:  track.ArtworkPath != "",
	}
}

// ToDTOList converts a slice of Tracks to TrackDTOs
func (m *TrackMapper) ToDTOList(tracks []*model.Track) []*dto.TrackDTO {
	if tracks == nil {
		return nil
	}

	dtos := make([]*dto.TrackDTO, len(tracks))
	for i, track := range tracks {
		dtos[i] = m.ToDTO(track)
	}

	return dtos
}