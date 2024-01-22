package api

import (
	"fio_service/internal/payload"
	"github.com/sirupsen/logrus"
	"os"
	"reflect"
	"testing"
)

var (
	log          *logrus.Logger
	humanName    = "michael"
	errHumanName = "errorname"
)

func TestMain(m *testing.M) {
	log = logrus.New()
	log.SetLevel(logrus.DebugLevel)
	os.Exit(m.Run())
}

func TestEnrichment_GetAge(t *testing.T) {
	type fields struct {
		log *logrus.Logger
	}
	type args struct {
		humanName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantAge int64
		wantErr bool
	}{
		{
			name:    "Успешное получение данных с API",
			fields:  fields{log: log},
			args:    args{humanName: humanName},
			wantAge: 63,
			wantErr: false,
		},
		{
			name:    "Имя не найдено",
			fields:  fields{log: log},
			args:    args{humanName: errHumanName},
			wantAge: 0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Enrichment{
				log: tt.fields.log,
			}
			gotAge, err := e.GetAge(tt.args.humanName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotAge, tt.wantAge) {
				t.Errorf("GetAge() gotAge = %v, want %v", gotAge, tt.wantAge)
			}
		})
	}
}

func TestEnrichment_GetGender(t *testing.T) {
	type fields struct {
		log *logrus.Logger
	}
	type args struct {
		humanName string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantGender string
		wantErr    bool
	}{
		{
			name:       "Успешное получение пола",
			fields:     fields{log: log},
			args:       args{humanName: humanName},
			wantGender: "male",
			wantErr:    false,
		},
		{
			name:       "Ошибка получения пола",
			fields:     fields{log: log},
			args:       args{humanName: errHumanName},
			wantGender: "",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Enrichment{
				log: tt.fields.log,
			}
			gotGender, err := e.GetGender(tt.args.humanName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGender() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGender, tt.wantGender) {
				t.Errorf("GetGender() gotGender = %v, want %v", gotGender, tt.wantGender)
			}
		})
	}
}

func TestEnrichment_GetNationality(t *testing.T) {
	type fields struct {
		log *logrus.Logger
	}
	type args struct {
		humanName string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantNationality []payload.Nationality
		wantErr         bool
	}{
		{
			name:   "Успешное получение национальности",
			fields: fields{log: log},
			args:   args{humanName: humanName},
			wantNationality: []payload.Nationality{
				{"AT", 0.061},
				{"DE", 0.056},
				{"DK", 0.054},
				{"IE", 0.048},
				{"GH", 0.046},
			},
			wantErr: false,
		},
		{
			name:            "Ошибка получения национальности",
			fields:          fields{log: log},
			args:            args{humanName: errHumanName},
			wantNationality: []payload.Nationality{},
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Enrichment{
				log: tt.fields.log,
			}
			gotNationality, err := e.GetNationality(tt.args.humanName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNationality() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotNationality, tt.wantNationality) {
				t.Errorf("GetNationality() gotNationality = %v, want %v", gotNationality, tt.wantNationality)
			}
		})
	}
}
