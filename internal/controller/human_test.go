package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fio_service/internal/database"
	"fio_service/internal/middleware/api"
	"fio_service/internal/payload"
	"fio_service/internal/service"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

/**
Если тесты выдают ошибку вызова API, необходимо использовать VPN
*/

type (
	funcValidate func(t *testing.T, rr *httptest.ResponseRecorder)
)

const (
	ExistTestName    = "Dmitry"
	NotValidTestName = "NotValidName"

	TestApiUrl = "http://localhost:8100/api/v1/humans"
)

var (
	humanId      string
	humanService service.Human
)

var (
	successHumanCreateBody = payload.HumanCreate{
		Name:       ExistTestName,
		Surname:    "Surname",
		Patronymic: "Patronymic",
	}
	successHumanBody = payload.Human{
		Name:       ExistTestName,
		Surname:    successHumanCreateBody.Surname,
		Patronymic: successHumanCreateBody.Patronymic,
		Age:        39,
		Gender:     "male",
		Nationalities: []payload.Nationality{
			{
				CountryID:   "RU",
				Probability: 0.403,
			},
			{
				CountryID:   "UA",
				Probability: 0.203,
			},
			{
				CountryID:   "BY",
				Probability: 0.168,
			},
			{
				CountryID:   "IL",
				Probability: 0.054,
			},
			{
				CountryID:   "LV",
				Probability: 0.033,
			},
		},
	}
)

func TestMain(m *testing.M) {
	var ctx = context.Background()

	const (
		TestDbHost = "localhost"
		TestDbPort = "5432"
		TestDbName = "test"
		TestDbUser = "test"
		TestDbPass = "test"
	)

	var testDbUrl = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		TestDbUser, TestDbPass, TestDbHost, TestDbPort, TestDbName)

	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)

	pool, err := pgxpool.New(ctx, testDbUrl)
	if err != nil {
		panic(fmt.Sprintf("Ошибка подключения к тестовой БД: %s", err))
	}

	db := database.NewHuman(ctx, pool)

	enrichment := api.NewEnrichment(log)

	humanService = service.NewHuman(log, db, enrichment)

	os.Exit(m.Run())
}

func TestHuman_Create(t *testing.T) {
	bodySuccess, err := json.Marshal(successHumanCreateBody)
	if err != nil {
		t.Fail()
	}

	var notValidHumanCreateBody = payload.HumanCreate{
		Name:       NotValidTestName,
		Surname:    "",
		Patronymic: "Patr",
	}
	bodyNotValid, err := json.Marshal(notValidHumanCreateBody)
	if err != nil {
		t.Fail()
	}

	type fields struct {
		service service.Human
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		fields fields
		name   string
		args   args
		fv     funcValidate
	}{
		{
			name:   "Успешное создание человека",
			fields: fields{service: humanService},
			args: args{
				r: httptest.NewRequest(
					http.MethodPost,
					fmt.Sprintf("%s/?name=%s", TestApiUrl, ExistTestName),
					bytes.NewBuffer(bodySuccess),
				),
			},
			fv: func(t *testing.T, rr *httptest.ResponseRecorder) {
				validateTestResponse(t, http.StatusCreated, rr)
			},
		},
		{
			name:   "Ошибка валидации при неправильно заполненных данных",
			fields: fields{service: humanService},
			args: args{
				r: httptest.NewRequest(
					http.MethodPost,
					fmt.Sprintf("%s/?name=%s", TestApiUrl, NotValidTestName),
					bytes.NewBuffer(bodyNotValid),
				),
			},
			fv: func(t *testing.T, rr *httptest.ResponseRecorder) {
				validateTestResponse(t, http.StatusBadRequest, rr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Human{
				service: tt.fields.service,
			}
			rr := httptest.NewRecorder()
			h.Create(rr, tt.args.r)
			tt.fv(t, rr)
		})
	}
}

func TestHuman_GetAll(t *testing.T) {
	type fields struct {
		service service.Human
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantHumans []payload.Human
		fv         funcValidate
	}{
		{
			name:   "Успешное получение по имени",
			fields: fields{service: humanService},
			args: args{
				r: httptest.NewRequest(
					http.MethodGet,
					fmt.Sprintf("%s/?name=%s&start=%d&size=%d", TestApiUrl, ExistTestName, 0, 1),
					nil,
				),
			},
			wantHumans: []payload.Human{
				successHumanBody,
			},
			fv: func(t *testing.T, rr *httptest.ResponseRecorder) {

			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Human{
				service:         tt.fields.service,
				paginationLimit: 100,
			}

			rr := httptest.NewRecorder()
			h.GetAll(rr, tt.args.r)

			validateTestResponse(t, http.StatusOK, rr)

			var gotHumans struct {
				Humans []payload.Human `json:"humans"`
			}
			if err := json.NewDecoder(rr.Body).Decode(&gotHumans); err != nil {
				t.Errorf("Ошибка декодирования JSON: %s", err)
			}

			if len(gotHumans.Humans) == 0 {
				t.Error("Люди не найдены")
			}

			humanId = gotHumans.Humans[0].Id
			tt.wantHumans[0].Id = humanId

			if !reflect.DeepEqual(gotHumans.Humans, tt.wantHumans) {
				t.Errorf("Pagination() gotHumans = %v, want %v", gotHumans, tt.wantHumans)
			}
		})
	}
}

func TestHuman_PatchByID(t *testing.T) {
	var successBodyUpdateHuman = payload.HumanUpdate{
		Name:       "UpdatedName",
		Surname:    "UpdatedSurname",
		Patronymic: "UpdatedPatronymic",
		Age:        43,
		Gender:     "male",
	}
	body, err := json.Marshal(successBodyUpdateHuman)
	if err != nil {
		t.Fail()
	}

	type fields struct {
		service service.Human
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		fv     funcValidate
	}{
		{
			name:   "Успешное обновление человека",
			fields: fields{service: humanService},
			args: args{
				r: httptest.NewRequest(
					http.MethodPatch,
					fmt.Sprintf("%s/{id}", TestApiUrl),
					bytes.NewBuffer(body),
				),
			},
			fv: func(t *testing.T, rr *httptest.ResponseRecorder) {
				checkTestStatusCode(http.StatusCreated, rr.Result().StatusCode, t, rr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Human{
				service: tt.fields.service,
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", humanId)
			tt.args.r = tt.args.r.WithContext(
				context.WithValue(
					tt.args.r.Context(),
					chi.RouteCtxKey,
					rctx,
				),
			)

			rr := httptest.NewRecorder()
			h.PatchByID(rr, tt.args.r)
			tt.fv(t, rr)
		})
	}
}

func TestHuman_DeleteByID(t *testing.T) {
	type fields struct {
		service service.Human
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		fv     funcValidate
	}{
		{
			name:   "Успешное удаление человека",
			fields: fields{service: humanService},
			args: args{
				r: httptest.NewRequest(
					http.MethodDelete,
					fmt.Sprintf("%s/{id}", TestApiUrl),
					nil,
				),
			},
			fv: func(t *testing.T, rr *httptest.ResponseRecorder) {
				checkTestStatusCode(http.StatusNoContent, rr.Result().StatusCode, t, rr)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Human{
				service: tt.fields.service,
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", humanId)
			tt.args.r = tt.args.r.WithContext(
				context.WithValue(
					tt.args.r.Context(),
					chi.RouteCtxKey,
					rctx,
				),
			)

			rr := httptest.NewRecorder()
			h.DeleteByID(rr, tt.args.r)
			tt.fv(t, rr)
		})
	}
}

func checkTestStatusCode(want, have int, t *testing.T, rr *httptest.ResponseRecorder) {
	if want != have {
		var body struct {
			Error string `json:"error"`
		}

		if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
			t.Errorf("Ошибка декодирования ошибки: %s", err)
		}

		t.Errorf("Ошибка выполнения запроса: %s", body.Error)
	}
}

func validateTestResponse(t *testing.T, expectedStatusCode int, rr *httptest.ResponseRecorder) {
	t.Logf("Статус код: %d", rr.Code)
	res := rr.Result()
	checkTestStatusCode(res.StatusCode, expectedStatusCode, t, rr)
}
