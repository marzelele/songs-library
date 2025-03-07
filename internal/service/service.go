package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"songs-library/internal"
	"songs-library/internal/models"
	"strings"
)

type Service struct {
	log             *slog.Logger
	repo            internal.Repository
	songsInfoAPIURL string
}

func NewService(log *slog.Logger, repo internal.Repository, songsInfoAPIURL string) internal.Service {
	return &Service{
		log:             log,
		repo:            repo,
		songsInfoAPIURL: songsInfoAPIURL,
	}
}

func (s *Service) CreateSong(in *models.CreateSong) (*models.Song, error) {
	const op = "service.CreateSong"

	log := s.log.With(
		slog.String("op", op),
	)

	details, err := s.getSongDetail(in.Song, in.Group)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	song := models.Song{
		Song:        in.Song,
		Group:       in.Group,
		ReleaseDate: details.ReleaseDate,
		Text:        details.Text,
		Link:        details.Link,
	}

	id, err := s.repo.CreateSong(&song)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	song.ID = id

	log.Info("created song", slog.Any("song", song))
	return &song, err
}

func (s *Service) DeleteSong(id int) error {
	const op = "service.DeleteSong"

	log := s.log.With(
		slog.String("op", op),
	)

	err := s.repo.DeleteSong(id)
	if err != nil {
		return err
	}

	log.Info("deleted song", slog.Int("songID", id))

	return nil
}

func (s *Service) UpdateSong(song *models.UpdateSong) (*models.UpdateSong, error) {
	const op = "service.CreateSong"

	log := s.log.With(
		slog.String("op", op),
	)

	err := s.repo.UpdateSong(song)
	if err != nil {
		return nil, err
	}

	log.Debug("updated song", slog.Any("song", song))

	return song, nil
}

func (s *Service) ListSongs(filter *models.SongsFilter) (models.Songs, error) {
	return s.repo.ListSongs(filter)
}

func (s *Service) GetTextBySongID(in *models.GetText) (*models.Text, error) {
	text, err := s.repo.GetTextBySongID(in.SongID)
	if err != nil {
		return nil, err
	}

	return &models.Text{
		SongID: in.SongID,
		Text:   s.paginateText(text, in.Page, in.PerPage),
	}, nil
}

func (s *Service) paginateText(text string, page, perPage int) string {
	verses := strings.Split(text, "\n\n")

	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = len(verses)
	}

	totalVerses := len(verses)
	if totalVerses == 0 {
		return ""
	}
	start := (page - 1) * perPage
	if start >= totalVerses {
		return ""
	}

	end := start + perPage
	if end > totalVerses {
		end = totalVerses
	}

	return strings.Join(verses[start:end], "\n\n")
}

func (s *Service) getSongDetail(song, group string) (models.SongDetail, error) {
	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)

	resp, err := http.Get(s.songsInfoAPIURL + "/info" + "?" + params.Encode())
	if err != nil {
		return models.SongDetail{}, fmt.Errorf("cannot get song detail: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return models.SongDetail{}, fmt.Errorf("request failed: %v", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.SongDetail{}, fmt.Errorf("cannot get response body: %w", err)
	}

	var songDetail models.SongDetail
	if err = json.Unmarshal(body, &songDetail); err != nil {
		return models.SongDetail{}, fmt.Errorf("cannot unmarshal song detail: %w", err)
	}

	err = resp.Body.Close()
	if err != nil {
		return models.SongDetail{}, fmt.Errorf("cannot close response body: %w", err)
	}

	return songDetail, nil
}
