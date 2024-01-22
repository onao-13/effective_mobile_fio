package router

import (
	"fio_service/internal/controller"
	"fio_service/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"net/http"
	"time"
)

const (
	Timeout = 10 * time.Second

	RequestLimit         = 10
	RequestLimitDuration = 5 * time.Second
)

// Router REST API сервиса
func Router(humanController controller.Human) *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		middleware.AllowContentType("application/json"),
		middleware.CleanPath,
		middleware.Logger,
		middleware.Timeout(Timeout),
		httprate.LimitByIP(RequestLimit, RequestLimitDuration),
	)

	r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		handler.NotFound(writer, "Указанный URL не существует")
	})

	// API
	r.Route("/api", func(r chi.Router) {

		// группа v1
		r.Group(func(r chi.Router) {
			r.Route("/v1", func(r chi.Router) {

				// группа API к данным человека
				r.Group(func(r chi.Router) {
					r.Route("/humans", func(r chi.Router) {
						r.Get("/", humanController.GetAll)
						r.Post("/", humanController.Create)

						r.Route("/{id}", func(r chi.Router) {
							r.Patch("/", humanController.PatchByID)
							r.Delete("/", humanController.DeleteByID)
						})
					})
				})
			})
		})
	})

	return r
}
