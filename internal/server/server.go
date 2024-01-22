package server

import (
	"context"
	"fio_service/internal/config"
	"fio_service/internal/controller"
	"fio_service/internal/database"
	"fio_service/internal/middleware/api"
	"fio_service/internal/router"
	"fio_service/internal/service"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	log  *logrus.Logger
	srv  *http.Server
	ctx  context.Context
	cfg  config.Config
	pool *pgxpool.Pool
}

// NewServer создает новый сервер
func NewServer(log *logrus.Logger, cfg config.Config) Server {
	return Server{
		log: log,
		cfg: cfg,
	}
}

// Serve запускает сервер
func (s Server) Serve() {

	s.log.Infoln("Сервер запускается")

	s.ctx = context.Background()
	var err error

	migration("migrations", s.cfg.DbURL(), s.log)
	s.log.Infoln("Миграция завершена")

	s.pool, err = pgxpool.New(s.ctx, s.cfg.DbURL())
	if err != nil {
		panic(fmt.Sprintf("Ошибка подключение к БД: %s", err))
	}

	// DATABASE
	humanDb := database.NewHuman(s.ctx, s.pool)

	// MIDDLEWARE
	enrichmentMiddleware := api.NewEnrichment(s.log)

	// SERVICE
	humanService := service.NewHuman(s.log, humanDb, enrichmentMiddleware)

	// CONTROLLER
	humanController := controller.NewHuman(humanService, s.cfg)

	mux := router.Router(humanController)

	s.log.Infoln("Все системы запущены")

	s.srv = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.Port),
		Handler: mux,
	}

	s.log.Infoln("Сервер запущен на порту ", s.cfg.Port)

	if err := s.srv.ListenAndServe(); err != nil {
		panic(fmt.Sprintf("Ошибка запуска сервера: %s", err))
	}
}

// Down останавливает сервер и производит отключение к БД
func (s Server) Down() {
	s.pool.Close()
	s.srv.Shutdown(s.ctx)
}
