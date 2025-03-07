package models

import "errors"

var (
	ErrInvalidSongID   = errors.New("invalid song_id parameter")
	ErrSongIsRequired  = errors.New("song is required")
	ErrGroupIsRequired = errors.New("group is required")
)

type Song struct {
	ID          int    `json:"id"`
	Song        string `json:"song"`
	Group       string `json:"group"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"-"`
	Link        string `json:"link"`
}

type Songs []Song

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type UpdateSong struct {
	ID          int    `json:"id"`
	Song        string `json:"song"`
	Group       string `json:"group"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func (s *UpdateSong) Validate() error {
	if s.ID <= 0 {
		return ErrInvalidSongID
	}

	if s.Song == "" {
		return ErrSongIsRequired
	}

	if s.Group == "" {
		return ErrGroupIsRequired
	}

	if s.ReleaseDate == "" {
		return errors.New("release_date is required")
	}

	return nil
}

type CreateSong struct {
	Song  string `json:"song"`
	Group string `json:"group"`
}

func (c *CreateSong) Validate() error {
	if c.Song == "" {
		return ErrSongIsRequired
	}

	if c.Group == "" {
		return ErrGroupIsRequired
	}

	return nil
}

type SongsFilter struct {
	IDs         []int  `json:"ids"`
	Song        string `json:"song"`
	Group       string `json:"group"`
	ReleaseDate string `json:"release_date"`
	Link        string `json:"link"`
	Page        int    `json:"page"`
	Limit       int    `json:"limit"`
}

type Text struct {
	SongID int    `json:"song_id"`
	Text   string `json:"text"`
}

type GetText struct {
	SongID  int `json:"song_id"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}

func (s *GetText) Validate() error {
	if s.SongID <= 0 {
		return ErrInvalidSongID
	}

	return nil
}
