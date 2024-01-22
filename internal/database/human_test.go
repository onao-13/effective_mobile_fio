package database

import (
	"context"
	"fio_service/internal/payload"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
	"reflect"
	"testing"
)

var pool *pgxpool.Pool

const (
	postgresTestUrl = "postgres://test:test@localhost:5432/test"

	nameTest = "testname"
)

var humanTestId string

var humanTest = payload.Human{
	Name:       nameTest,
	Surname:    "testsurname",
	Patronymic: "testpatronymic",
	Age:        10,
	Gender:     "male",
	Nationalities: []payload.Nationality{
		{
			CountryID:   "RU",
			Probability: 0.66,
		},
		{
			CountryID:   "US",
			Probability: 0.131,
		},
	},
}

func TestMain(m *testing.M) {
	var err error

	pool, err = pgxpool.New(context.TODO(), postgresTestUrl)
	if err != nil {
		panic(fmt.Errorf("ошибка подключения к тестовой БД: %s", err))
	}

	os.Exit(m.Run())
}

func TestHuman_Create(t *testing.T) {
	type fields struct {
		ctx  context.Context
		pool *pgxpool.Pool
	}
	type args struct {
		human payload.Human
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Успешное создание человека",
			fields: fields{
				ctx:  context.TODO(),
				pool: pool,
			},
			args: args{
				human: humanTest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Human{
				ctx:  tt.fields.ctx,
				pool: tt.fields.pool,
			}
			if err := h.Create(tt.args.human); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHuman_GetAll(t *testing.T) {
	type fields struct {
		ctx  context.Context
		pool *pgxpool.Pool
	}
	type args struct {
		id         string
		name       string
		surname    string
		patronymic string
		gender     string
		age        int64
		start      int64
		limit      int64
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantHumans []payload.Human
		wantErr    bool
	}{
		{
			name: "Успешный поиск людей по имени",
			fields: fields{
				ctx:  context.TODO(),
				pool: pool,
			},
			args: args{
				id:         "",
				name:       nameTest,
				surname:    "",
				patronymic: "",
				age:        0,
				gender:     "",
				start:      0,
				limit:      5,
			},
			wantHumans: []payload.Human{
				humanTest,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Human{
				ctx:  tt.fields.ctx,
				pool: tt.fields.pool,
			}
			gotHumans, err := h.Pagination(tt.args.id, tt.args.name, tt.args.surname, tt.args.patronymic, tt.args.gender, tt.args.age, tt.args.start, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pagination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			humanTestId = gotHumans[0].Id
			tt.wantHumans[0].Id = humanTestId

			if !reflect.DeepEqual(gotHumans, tt.wantHumans) {
				t.Errorf("Pagination() gotHumans = %v, want %v", gotHumans, tt.wantHumans)
			}
		})
	}
}

func TestHuman_Patch(t *testing.T) {
	type fields struct {
		ctx  context.Context
		pool *pgxpool.Pool
	}
	type args struct {
		id    string
		human payload.HumanUpdate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Успешное обновление данных пользователя",
			fields: fields{ctx: context.TODO(), pool: pool},
			args: args{
				id: humanTestId,
				human: payload.HumanUpdate{
					Name:       "testname2",
					Surname:    humanTest.Surname,
					Patronymic: humanTest.Patronymic,
					Age:        12,
					Gender:     "female",
					AddedNationalities: []payload.Nationality{
						{
							CountryID:   "TID",
							Probability: 0.234,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Human{
				ctx:  tt.fields.ctx,
				pool: tt.fields.pool,
			}
			if err := h.Update(tt.args.id, tt.args.human); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHuman_Delete(t *testing.T) {
	type fields struct {
		ctx  context.Context
		pool *pgxpool.Pool
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Успешное удаление человека",
			fields:  fields{ctx: context.TODO(), pool: pool},
			args:    args{id: humanTestId},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Human{
				ctx:  tt.fields.ctx,
				pool: tt.fields.pool,
			}
			if err := h.Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

//func dropTestTableContent(t *testing.T) {
//	var dropScrips = []string{
//		`DELETE FROM humans_nationality`,
//		`DELETE FROM humans`,
//	}
//
//	for _, sql := range dropScrips {
//		if _, err := pool.Exec(context.TODO(), sql); err != nil {
//			t.Errorf("ошибка чистки БД: %s", err)
//		}
//	}
//}
