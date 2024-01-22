package main

import (
	"fio_service/internal/config"
	"fio_service/internal/server"
	"github.com/sirupsen/logrus"
	"strconv"
)

func main() {
	cfg := config.Load()
	log := logrus.New()

	debugMode(cfg, log)

	srv := server.NewServer(log, cfg)

	srv.Serve()

	defer srv.Down()
}

// debugMode проверяет, включен ли флаг отладки.
// Если да, то включает логи с режимом "DEBUG"
func debugMode(cfg config.Config, log *logrus.Logger) {
	isDebug, err := strconv.ParseBool(cfg.DebugMode)
	if err != nil {
		return
	}

	if isDebug {
		log.SetLevel(logrus.DebugLevel)
		log.Debugln("Включен режим отладки")
	}
}
