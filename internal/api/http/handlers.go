package http

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"songs-library/internal"
	"songs-library/internal/models"
	"songs-library/internal/respository"
	"songs-library/pkg/api/response"
	"songs-library/pkg/logger/sl"
	"strconv"
)

type Handler struct {
	log     *slog.Logger
	service internal.Service
}

func NewHandler(log *slog.Logger, service internal.Service) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

// CreateSong godoc
// @Summary      Create a song
// @Description  Добавление новой песни
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        song  body      models.CreateSong  true              "song and group"
// @Success      200   {object}  response.Response{data=models.Song}  "OK"
// @Failure      400   {object}  response.Response                    "Bad Request"
// @Failure      500   {object}  response.Response                    "Internal Server Error"
// @Router       /songs [post]
func (h *Handler) CreateSong(w http.ResponseWriter, r *http.Request) {
	const op = "handler.CreateSong"
	log := h.setLogger(r.Context(), op, h.log)

	var req models.CreateSong

	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("request body is empty"))
		return
	}
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to decode request body"))
		return
	}

	log.Info("request body decoded", slog.Any("request", req))

	err = req.Validate()
	if err != nil {
		log.Error("failed to validate request", sl.Err(err))
		w.WriteHeader(http.StatusBadRequest)

		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	song, err := h.service.CreateSong(&req)
	if err != nil {
		log.Error("failed to create song", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to create song"))
		return
	}

	render.JSON(w, r, response.OK(song))
}

// DeleteSong godoc
// @Summary      Delete a song
// @Description  Удаление песни
// @Tags         Songs
// @Param        id   path     int     true  "song_id"
// @Success      200  {object}  response.Response "OK"
// @Failure      400  {object}  response.Response "Bad Request"
// @Failure      404  {object}  response.Response "Song Not Found"
// @Failure      500  {object}  response.Response "Internal Server Error"
// @Router       /songs/{id} [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteSong"
	log := h.setLogger(r.Context(), op, h.log)

	str := chi.URLParam(r, "id")
	id, err := strconv.Atoi(str)
	if err != nil {
		log.Error("failed to decode id parameter", sl.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error(models.ErrInvalidSongID.Error()))
		return
	}

	if id <= 0 {
		log.Error("invalid song_id")

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error(models.ErrInvalidSongID.Error()))
		return
	}

	err = h.service.DeleteSong(id)
	if errors.Is(err, respository.ErrSongNotFound) {
		log.Error("song not found", sl.Err(err))

		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, response.Error("song not found"))
		return
	}
	if err != nil {
		log.Error("failed to delete song", sl.Err(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to delete song"))
		return
	}

	render.JSON(w, r, response.OK(nil))
}

// UpdateSong godoc
// @Summary      Update a song
// @Description  Изменение данных песни
// @Tags         Songs
// @Accept       json
// @Param        song  body      models.UpdateSong  true "Song Attrs"
// @Success      200   {object}  response.Response{data=models.UpdateSong}  "OK"
// @Failure      400   {object}  response.Response                          "Bad Request"
// @Failure      404   {object}  response.Response                          "Song Not Found"
// @Failure      500   {object}  response.Response                          "Internal Server Error"
// @Router       /songs [put]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	const op = "handler.UpdateSong"
	log := h.setLogger(r.Context(), op, h.log)

	var req models.UpdateSong

	err := render.DecodeJSON(r.Body, &req)
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("request body is empty"))
		return
	}
	if err != nil {
		log.Error("failed to decode request body", sl.Err(err))
		render.JSON(w, r, response.Error("failed to decode request body"))
		return
	}

	log.Info("request body decoded", slog.Any("request", req))

	err = req.Validate()
	if err != nil {
		log.Error("failed to validate request", sl.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	song, err := h.service.UpdateSong(&req)
	if errors.Is(err, respository.ErrSongNotFound) {
		log.Error("song not found", sl.Err(err))

		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, response.Error("song not found"))
		return
	}
	if err != nil {
		log.Error("failed to update song", sl.Err(err))
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to update song"))
		return
	}

	render.JSON(w, r, response.OK(song))
}

// ListSongs godoc
// @Summary      Get list of songs
// @Description  Получение данных библиотеки с фильтрацией по всем полям и пагинацией
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        song  body      models.SongsFilter  false             "songs filters"
// @Success      200   {object}  response.Response{data=models.Songs}  "OK"
// @Failure      400   {object}  response.Response                     "Bad Request"
// @Failure      500   {object}  response.Response                     "Internal Server Error"
// @Router       /songs/list [post]
func (h *Handler) ListSongs(w http.ResponseWriter, r *http.Request) {
	const op = "handler.ListSongs"
	log := h.setLogger(r.Context(), op, h.log)

	var req models.SongsFilter

	err := render.DecodeJSON(r.Body, &req)
	if err != nil && !errors.Is(err, io.EOF) {
		log.Error("failed to decode request body", sl.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("failed to decode request body"))
		return
	}

	list, err := h.service.ListSongs(&req)
	if err != nil {
		log.Error("failed to list songs", sl.Err(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to list songs"))
		return
	}

	render.JSON(w, r, response.OK(list))
}

// GetTextBySongID godoc
// @Summary      Get song's text
// @Description  Получение текста песни с пагинацией по куплетам
// @Tags         Texts
// @Accept       json
// @Produce      json
// @Param        id   query     int     true  "song_id"
// @Param        page   query     int     false  "page"
// @Param        perPage   query     int     false  "per page"
// @Success      200   {object}  response.Response{data=models.Text}  "OK"
// @Failure      400   {object}  response.Response                    "Bad Request"
// @Failure      404   {object}  response.Response                    "Song Not Found"
// @Failure      500   {object}  response.Response                    "Internal Server Error"
// @Router       /songs/texts [get]
func (h *Handler) GetTextBySongID(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetTextBySongID"
	log := h.setLogger(r.Context(), op, h.log)

	var (
		err error
		req models.GetText
	)

	req.SongID, err = strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		log.Error("failed to parse id parameter", sl.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error(models.ErrInvalidSongID.Error()))
		return
	}

	req.Page, err = strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		req.Page = 0
	}

	req.PerPage, err = strconv.Atoi(r.URL.Query().Get("perPage"))
	if err != nil {
		req.PerPage = 0
	}

	err = req.Validate()
	if err != nil {
		log.Error("failed to validate request", sl.Err(err))

		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error(err.Error()))
		return
	}

	text, err := h.service.GetTextBySongID(&req)

	if errors.Is(err, respository.ErrSongNotFound) {
		log.Error("song not found", sl.Err(err))

		w.WriteHeader(http.StatusNotFound)
		render.JSON(w, r, response.Error("song not found"))
		return
	}
	if err != nil {
		log.Error("failed to get song", sl.Err(err))

		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to get song"))
		return
	}

	render.JSON(w, r, response.OK(text))
}

func (h *Handler) setLogger(ctx context.Context, op string, log *slog.Logger) *slog.Logger {
	return log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)
}
