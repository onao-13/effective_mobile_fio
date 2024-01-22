package controller

import (
	"encoding/json"
	"fio_service/internal/config"
	"fio_service/internal/errors"
	"fio_service/internal/handler"
	"fio_service/internal/payload"
	"fio_service/internal/service"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

// Human контроллер для взаимодействия с API людей
type Human struct {
	service         service.Human
	paginationLimit int64
}

func NewHuman(service service.Human, cfg config.Config) (h Human) {
	var err error

	h.service = service
	h.paginationLimit, err = strconv.ParseInt(cfg.PaginationLimit, 10, 32)
	if err != nil {
		panic(fmt.Sprintf("Ошибка установки параметра пагинации"))
	}

	return
}

// GetAll получение всех людей
func (h Human) GetAll(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()
	var (
		// параметры запроса
		id         = urlQuery.Get("id")
		name       = urlQuery.Get("name")
		surname    = urlQuery.Get("surname")
		patronymic = urlQuery.Get("patronymic")
		ageStr     = urlQuery.Get("age")
		gender     = urlQuery.Get("gender")

		// пагинация
		startStr = urlQuery.Get("start")
		sizeStr  = urlQuery.Get("size")

		age int64

		err error
	)

	if len(ageStr) != 0 {
		age, err = strconv.ParseInt(ageStr, 10, 64)
		if err != nil {
			handler.BadRequest(w, "Недопустимый тип age. Age должен быть числом")
			return
		}
	}

	start, err := strconv.ParseInt(startStr, 10, 32)
	if err != nil {
		handler.BadRequest(w, "Ошибка параметра start. Start должен быть числом")
		return
	}

	size, err := strconv.ParseInt(sizeStr, 10, 32)
	if err != nil {
		handler.BadRequest(w, "Ошибка параметра size. Size должен быть числом")
		return
	}

	if size > h.paginationLimit {
		handler.BadRequest(w, errors.ErrExcessPaginationSize)
		return
	}

	humans, err := h.service.GetAll(id, name, surname, patronymic, gender, age, start, size)
	if err != nil {
		switch err.(type) {
		case *errors.ErrDataNotFound:
			handler.NotFound(w, err.Error())
			return
		default:
			handler.InternalServerError(w, errors.ErrServer)
			return
		}
	}

	body := map[string]interface{}{
		"humans": humans,
	}

	handler.OkData(w, body)
}

// Create создание нового человека
func (h Human) Create(w http.ResponseWriter, r *http.Request) {
	var (
		human payload.HumanCreate
		err   error
	)

	if err = json.NewDecoder(r.Body).Decode(&human); err != nil {
		handler.BadRequest(w, errors.ErrJSONDecode)
		return
	}

	if err = human.Validate(); err != nil {
		handler.BadRequest(w, err.Error())
		return
	}

	if err = h.service.Create(human); err != nil {
		handler.InternalServerError(w, errors.ErrServer)
		return
	}

	handler.Create(w, "")
}

// PatchByID обновление человека по ID
func (h Human) PatchByID(w http.ResponseWriter, r *http.Request) {
	var (
		id    string
		err   error
		human payload.HumanUpdate
	)

	id = chi.URLParam(r, "id")
	if len(id) == 0 {
		handler.BadRequest(w, "ID не указан")
		return
	}

	if err = json.NewDecoder(r.Body).Decode(&human); err != nil {
		handler.BadRequest(w, errors.ErrJSONDecode)
		return
	}

	if err = human.Validate(); err != nil {
		handler.BadRequest(w, err.Error())
		return
	}

	if err = h.service.Update(id, human); err != nil {
		handler.InternalServerError(w, errors.ErrServer)
		return
	}

	handler.Create(w, "")
}

// DeleteByID удаление человека по ID
func (h Human) DeleteByID(w http.ResponseWriter, r *http.Request) {
	var id string
	var err error

	id = chi.URLParam(r, "id")
	if err != nil {
		handler.BadRequest(w, "ID не указан")
		return
	}

	if err = h.service.Delete(id); err != nil {
		handler.InternalServerError(w, errors.ErrServer)
		return
	}

	handler.NoContent(w)
}
