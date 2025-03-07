package internal

import (
	"songs-library/internal/models"
)

type Service interface {
	CreateSong(*models.CreateSong) (*models.Song, error)
	UpdateSong(song *models.UpdateSong) (*models.UpdateSong, error)
	DeleteSong(int) error
	ListSongs(*models.SongsFilter) (models.Songs, error)
	GetTextBySongID(*models.GetText) (*models.Text, error)
}
