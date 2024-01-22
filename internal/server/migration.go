package server

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

// migration совершает миграцию в БД
func migration(sourceFilePath, dbUrl string, log *logrus.Logger) {
	m, err := migrate.New(fmt.Sprintf("file://%s", sourceFilePath), dbUrl)
	if err != nil {
		panic(fmt.Sprintf("Ошибка миграции БД: %s", err))
	}
	if err := m.Up(); err != nil {
		log.Infoln("Ошибка запуска миграции: ", err)
	}
}
