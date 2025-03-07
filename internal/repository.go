package internal

import (
	"songs-library/internal/models"
)

type Repository interface {
	CreateSong(*models.Song) (int, error)
	UpdateSong(song *models.UpdateSong) error
	DeleteSong(int) error
	ListSongs(*models.SongsFilter) (models.Songs, error)
	GetTextBySongID(int) (string, error)
}
