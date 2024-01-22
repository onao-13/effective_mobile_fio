package service

import (
	"fio_service/internal/database"
	"fio_service/internal/errors"
	"fio_service/internal/middleware/api"
	"fio_service/internal/payload"
	"github.com/sirupsen/logrus"
	"sync"
)

type Human struct {
	log        *logrus.Logger
	db         database.Human
	enrichment api.Enrichment
}

func NewHuman(log *logrus.Logger, db database.Human, enrichment api.Enrichment) Human {
	return Human{
		log:        log,
		db:         db,
		enrichment: enrichment,
	}
}

func (h Human) Create(humanCreate payload.HumanCreate) error {
	var err error
	var age int64
	var gender string
	var nationalities []payload.Nationality

	var wg *sync.WaitGroup = &sync.WaitGroup{}
	wg.Add(3)

	// age
	go func() {
		defer wg.Done()
		age, err = h.enrichment.GetAge(humanCreate.Name)
	}()

	// gender
	go func() {
		defer wg.Done()
		gender, err = h.enrichment.GetGender(humanCreate.Name)
	}()

	// nationalities
	go func() {
		defer wg.Done()
		nationalities, err = h.enrichment.GetNationality(humanCreate.Name)
	}()

	wg.Wait()

	if err != nil {
		return err
	}

	human := payload.Human{
		Name:          humanCreate.Name,
		Surname:       humanCreate.Surname,
		Patronymic:    humanCreate.Patronymic,
		Age:           age,
		Gender:        gender,
		Nationalities: nationalities,
	}

	if err := h.db.Create(human); err != nil {
		h.log.Errorln("Ошибка создания человека: ", err)
		return err
	}
	h.log.Debugln("Создан человек")

	return nil
}

func (h Human) GetAll(id, name, surname, patronymic, gender string, age, start, limit int64) (humans []payload.Human, err error) {
	humans, err = h.db.Pagination(id, name, surname, patronymic, gender, age, start, limit)
	if err != nil {
		h.log.Errorln("Ошибка пагинации: ", err)
		return
	}

	if len(humans) == 0 {
		return nil, &errors.ErrDataNotFound{}
	}

	return
}

func (h Human) Update(id string, human payload.HumanUpdate) error {
	if err := h.db.Update(id, human); err != nil {
		h.log.Errorln("Ошибка обновления данных человеку: ", err)
		return err
	}
	h.log.Debugln("Обновлен человек с ID ", id)
	return nil
}

func (h Human) Delete(id string) error {
	if err := h.db.Delete(id); err != nil {
		h.log.Errorln("Ошибка удаления человека: ", err)
		return err
	}
	h.log.Debugln("Удален человек с ID ", id)
	return nil
}
