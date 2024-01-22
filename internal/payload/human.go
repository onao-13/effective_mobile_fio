package payload

import "fmt"

type (
	// Human струрктура человека со всеми его полями
	Human struct {
		Id            string        `json:"id"`
		Name          string        `json:"name"`
		Surname       string        `json:"surname"`
		Patronymic    string        `json:"patronymic"`
		Age           int64         `json:"age"`
		Gender        string        `json:"gender"`
		Nationalities []Nationality `json:"nationalities"`
	}
	// HumanCreate структура создания человека
	HumanCreate struct {
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic"`
	}
	// HumanUpdate структура обновления данных человеку
	HumanUpdate struct {
		Id                   string        `json:"id"`
		Name                 string        `json:"name"`
		Surname              string        `json:"surname"`
		Patronymic           string        `json:"patronymic"`
		Age                  int64         `json:"age"`
		Gender               string        `json:"gender"`
		AddedNationalities   []Nationality `json:"added_nationalities"`
		DeletedNationalities []string      `json:"deleted_nationalities"`
	}
)

func (h Human) Validate() error {
	switch {
	case len(h.Id) == 0:
		return fmt.Errorf("ID не указан")
	case h.Age <= 0:
		return fmt.Errorf("возраст не указан")
	case len(h.Name) == 0:
		return fmt.Errorf("имя не указано")
	case len(h.Surname) == 0:
		return fmt.Errorf("фамилия не указана")
	case len(h.Patronymic) == 0:
		return fmt.Errorf("отчество не указано")
	case len(h.Gender) == 0:
		return fmt.Errorf("пол не указан")
	}

	return nil
}

func (h HumanCreate) Validate() error {
	switch {
	case len(h.Name) == 0:
		return fmt.Errorf("имя не указано")
	case len(h.Surname) == 0:
		return fmt.Errorf("фамилия не указана")
	case len(h.Patronymic) == 0:
		return fmt.Errorf("отчество не указано")
	}

	return nil
}

func (h HumanUpdate) Validate() error {
	switch {
	case h.Age <= 0:
		return fmt.Errorf("возраст не указан")
	case len(h.Name) == 0:
		return fmt.Errorf("имя не указано")
	case len(h.Surname) == 0:
		return fmt.Errorf("фамилия не указана")
	case len(h.Patronymic) == 0:
		return fmt.Errorf("отчество не указано")
	case len(h.Gender) == 0:
		return fmt.Errorf("пол не указан")
	}

	return nil
}
