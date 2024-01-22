package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Create обертка для ответа клиенту с кодом 201
func Create(w http.ResponseWriter, msg string) {
	if len(msg) == 0 {
		w.WriteHeader(http.StatusCreated)
		return
	}
	data := map[string]interface{}{
		"message": msg,
	}
	handle(w, http.StatusCreated, data)
}

// InternalServerError обертка для ответа клиенту с кодом 500
func InternalServerError(w http.ResponseWriter, msg string) {
	data := map[string]interface{}{
		"error": msg,
	}
	handle(w, http.StatusInternalServerError, data)
}

// BadRequest обертка для ответа клиенту с кодом 400
func BadRequest(w http.ResponseWriter, msg string) {
	data := map[string]interface{}{
		"error": msg,
	}
	handle(w, http.StatusBadRequest, data)
}

// OkData обертка для отправки данных клиенту с кодом 200
func OkData(w http.ResponseWriter, data interface{}) {
	handle(w, http.StatusOK, data)
}

// NotFound обертка для ответа клиенту с кодом 404
func NotFound(w http.ResponseWriter, msg string) {
	data := map[string]interface{}{
		"error": msg,
	}
	handle(w, http.StatusNotFound, data)
}

// NoContent обертка для ответа клиенту с кодом 204
func NoContent(w http.ResponseWriter) {
	handle(w, http.StatusNoContent, nil)
}

// handle ответ клиенту
func handle(w http.ResponseWriter, status int, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Errorln("Ошибка создания JSON ответа: ", err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(b)
}
