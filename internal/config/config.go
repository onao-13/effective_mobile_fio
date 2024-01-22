package config

import (
	"fmt"
	"os"
	"reflect"
)

// Config струткура настроек сервиса
type Config struct {
	Port            string
	DbUser          string
	DbPass          string
	DbHost          string
	DbPort          string
	DbName          string
	PaginationLimit string
	DebugMode       string
}

// Load загрузка конфига, считывая системные переменные
func Load() (c Config) {
	c = Config{
		Port:            os.Getenv("PORT"),
		DbUser:          os.Getenv("DB_USER"),
		DbPass:          os.Getenv("DB_PASS"),
		DbHost:          os.Getenv("DB_HOST"),
		DbPort:          os.Getenv("DB_PORT"),
		DbName:          os.Getenv("DB_NAME"),
		PaginationLimit: os.Getenv("PAGINATION_LIMIT"),
		DebugMode:       os.Getenv("DEBUG_MODE"),
	}

	if err := c.validate(); err != nil {
		panic(fmt.Sprintf("Ошибка инициализации конфигурации сервера: %s", err))
	}

	return
}

// validate проверяет, заполнены ли поля.
// Возвращает ошибку, если поле пустое
func (c Config) validate() error {
	fields := reflect.ValueOf(&c).Elem()

	for i := 0; i < fields.NumField(); i++ {
		if fields.Field(i).Len() == 0 {
			return fmt.Errorf("поле %v пустое", fields.Type().Field(i).Name)
		}
	}

	return nil
}

// DbURL возвращает путь подключения к базе PostgreSQL
func (c Config) DbURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DbUser, c.DbPass, c.DbHost, c.DbPort, c.DbName)
}
