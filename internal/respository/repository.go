package respository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log/slog"
	"songs-library/internal/consts"
	"songs-library/internal/converter"
	"songs-library/internal/models"
)

var ErrSongNotFound = errors.New("song not found")

type Repository struct {
	log *slog.Logger
	db  *sqlx.DB
}

func NewRepository(conn string) (*Repository, error) {
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) CreateSong(song *models.Song) (int, error) {
	const op = "repository.CreateSong"

	q := squirrel.Insert(consts.SongsTableName).
		PlaceholderFormat(squirrel.Dollar).
		Columns(consts.SongColumn, consts.GroupColumn, consts.ReleaseDateColumn, consts.TextColumn, consts.LinkColumn).
		Values(song.Song, song.Group, song.ReleaseDate, song.Text, song.Link).
		Suffix("RETURNING id")

	var id int
	err := q.RunWith(r.db).QueryRow().Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (r *Repository) UpdateSong(song *models.UpdateSong) error {
	const op = "repository.UpdateSong"

	q := squirrel.Update(consts.SongsTableName).
		PlaceholderFormat(squirrel.Dollar).
		Set(consts.SongColumn, song.Song).
		Set(consts.GroupColumn, song.Group).
		Set(consts.ReleaseDateColumn, song.ReleaseDate).
		Set(consts.TextColumn, song.Text).
		Set(consts.LinkColumn, song.Link).
		Where(squirrel.Eq{consts.IDColumn: song.ID})

	res, err := q.RunWith(r.db).Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return ErrSongNotFound
	}

	return nil
}

func (r *Repository) DeleteSong(id int) error {
	const op = "repository.DeleteSong"

	q := squirrel.Delete(consts.SongsTableName).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{consts.IDColumn: id})

	res, err := q.RunWith(r.db).Exec()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return ErrSongNotFound
	}

	return nil
}

func (r *Repository) ListSongs(filter *models.SongsFilter) (models.Songs, error) {
	const op = "repository.ListSongs"

	q := squirrel.
		Select(consts.IDColumn, consts.SongColumn, consts.GroupColumn, consts.ReleaseDateColumn, consts.LinkColumn).
		PlaceholderFormat(squirrel.Dollar).
		From(consts.SongsTableName).
		OrderBy(consts.IDColumn + " ASC")

	q = converter.SongFilterToSqlFilters(q, filter)

	rows, err := q.RunWith(r.db).Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	songs := make([]models.Song, 0, filter.Limit)

	for rows.Next() {
		var song models.Song
		if err = rows.Scan(
			&song.ID,
			&song.Song,
			&song.Group,
			&song.ReleaseDate,
			&song.Link); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		songs = append(songs, song)
	}

	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return songs, nil
}

func (r *Repository) GetTextBySongID(songID int) (string, error) {
	const op = "repository.GetTextBySongID"

	q := squirrel.Select(consts.TextColumn).
		PlaceholderFormat(squirrel.Dollar).
		From(consts.SongsTableName).
		Where(squirrel.Eq{consts.IDColumn: songID})

	var text string
	err := q.RunWith(r.db).QueryRow().Scan(&text)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrSongNotFound
	}
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return text, nil
}
