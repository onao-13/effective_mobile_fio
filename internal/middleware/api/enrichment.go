package api

import (
	"encoding/json"
	"fio_service/internal/payload"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	AgifyApiURL       = "https://api.agify.io"
	GenderizeApiURL   = "https://api.genderize.io"
	NationalizeApiURL = "https://api.nationalize.io"
)

// Enrichment обогащение данных через стороннее API
type Enrichment struct {
	log *logrus.Logger
}

func NewEnrichment(log *logrus.Logger) Enrichment {
	return Enrichment{log: log}
}

// GetAge получение возраста с API https://agify.io/
func (e Enrichment) GetAge(humanName string) (age int64, err error) {
	e.log.Debugln("Вызов API https://agify.io/")

	resp, err := http.Get(fmt.Sprintf("%s?name=%s", AgifyApiURL, humanName))
	if err != nil {
		e.log.Errorln("Ошибка вызова API возраста: ", err)
		return
	}

	var body payload.AgifyAPI
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		e.log.Errorln("Ошибка декодирования JSON: ", err)
		return
	}

	age = body.Age

	if age == 0 {
		return age, fmt.Errorf("возраст не найден")
	}

	return
}

// GetGender получение пола с API https://genderize.io/
func (e Enrichment) GetGender(humanName string) (gender string, err error) {
	e.log.Debugln("Вызов API https://genderize.io/")

	resp, err := http.Get(fmt.Sprintf("%s?name=%s", GenderizeApiURL, humanName))
	if err != nil {
		e.log.Errorln("Ошибка вызова API пола: ", err)
		return
	}

	var body payload.GenderizeAPI
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		e.log.Errorln("Ошибка декодирова JSON: ", err)
		return
	}

	gender = body.Gender

	if len(gender) == 0 {
		return gender, fmt.Errorf("пол не найден")
	}

	return
}

// GetNationality получение национальности с API https://nationalize.io/
func (e Enrichment) GetNationality(humanName string) (nationality []payload.Nationality, err error) {
	e.log.Debugln("Вызов API https://nationalize.io/")

	resp, err := http.Get(fmt.Sprintf("%s?name=%s", NationalizeApiURL, humanName))
	if err != nil {
		e.log.Errorln("Ошибка вызова API национальности: ", err)
		return
	}

	var body payload.NationalizeAPI
	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		e.log.Errorln("Ошибка декодирования JSON: ", err)
		return
	}

	nationality = body.Country

	if len(nationality) == 0 {
		return nationality, fmt.Errorf("национальность не найдена")
	}

	return
}
