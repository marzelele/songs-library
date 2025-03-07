package converter

import (
	"github.com/Masterminds/squirrel"
	"songs-library/internal/consts"
	"songs-library/internal/models"
)

func SongFilterToSqlFilters(q squirrel.SelectBuilder, filter *models.SongsFilter) squirrel.SelectBuilder {
	if len(filter.IDs) != 0 {
		q = q.Where(squirrel.Eq{consts.IDColumn: filter.IDs})
	}

	if filter.Song != "" {
		q = q.Where(squirrel.Like{consts.SongColumn: setLike(filter.Song)})
	}

	if filter.Group != "" {
		q = q.Where(squirrel.Like{consts.GroupColumn: setLike(filter.Group)})
	}

	if filter.ReleaseDate != "" {
		q = q.Where(squirrel.Like{consts.ReleaseDateColumn: setLike(filter.ReleaseDate)})
	}

	if filter.Link != "" {
		q = q.Where(squirrel.Like{consts.LinkColumn: setLike(filter.Link)})
	}

	if filter.Page < 1 {
		filter.Page = 1

	}

	if filter.Limit < 1 {
		filter.Limit = consts.DefaultLimit
	}

	q = q.Limit(uint64(filter.Limit))
	q = q.Offset(uint64((filter.Page - 1) * filter.Limit))

	return q
}

func setLike(s string) string {
	return "%" + s + "%"
}
