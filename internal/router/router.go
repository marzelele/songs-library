package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"songs-library/internal/api/http"
	"songs-library/pkg/middlewares"
)

type Router struct {
	log     *slog.Logger
	handler *http.Handler
}

func NewRouter(log *slog.Logger, handler *http.Handler) *Router {
	return &Router{
		log:     log,
		handler: handler,
	}
}

func (r *Router) Init() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middlewares.NewMiddlewareLogger(r.log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	router.Route("/api", func(router chi.Router) {
		router.Route("/v1", func(router chi.Router) {
			router.Route("/songs", func(router chi.Router) {
				router.Post("/", r.handler.CreateSong)
				router.Delete("/{id}", r.handler.DeleteSong)
				router.Put("/", r.handler.UpdateSong)
				router.Post("/list", r.handler.ListSongs)
				router.Route("/texts", func(router chi.Router) {
					router.Get("/", r.handler.GetTextBySongID)
				})
			})
		})
	})

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	return router
}
