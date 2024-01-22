package database

import (
	"context"
	"fio_service/internal/payload"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strings"
)

type Human struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func NewHuman(ctx context.Context, pool *pgxpool.Pool) Human {
	return Human{ctx: ctx, pool: pool}
}

func (h Human) Create(human payload.Human) error {
	sql := `
	INSERT INTO
		humans(name, surname, patronymic, age, gender)
	VALUES(@name, @surname, @patronymic, @age, @gender)
	RETURNING id
	`

	args := pgx.NamedArgs{
		"name":       human.Name,
		"surname":    human.Surname,
		"patronymic": human.Patronymic,
		"age":        human.Age,
		"gender":     human.Gender,
	}

	var humanId string
	if err := h.pool.QueryRow(h.ctx, sql, args).Scan(&humanId); err != nil {
		return err
	}

	var batch = &pgx.Batch{}

	for _, nationality := range human.Nationalities {
		h.batchAddNationality(batch, humanId, nationality)
	}

	br := h.pool.SendBatch(h.ctx, batch)

	defer br.Close()

	return nil
}

func (h Human) Pagination(
	id, name, surname, patronymic, gender string,
	age, start, limit int64,
) (humans []payload.Human, err error) {
	// получение списка людей
	sql := `
	SELECT
		id,
		name,
		surname,
		patronymic,
		age,
		gender
	FROM
	    humans
	@where
	LIMIT @limit
	OFFSET @start
	`

	// создание динамического условия запроса
	where := h.sqlWhereFilter(map[string]interface{}{
		"id":         id,
		"name":       name,
		"surname":    surname,
		"patronymic": patronymic,
		"age":        age,
		"gender":     gender,
	})
	// добавление условий выборки
	sql = strings.ReplaceAll(sql, "@where", where)

	args := pgx.NamedArgs{
		"limit": limit,
		"start": start,
	}

	rows, err := h.pool.Query(h.ctx, sql, args)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var idHuman = make(map[string]*payload.Human)
	var batch = &pgx.Batch{}

	for rows.Next() {
		var human payload.Human
		err = rows.Scan(
			&human.Id, &human.Name, &human.Surname,
			&human.Patronymic, &human.Age, &human.Gender,
		)
		if err != nil {
			return nil, err
		}

		idHuman[human.Id] = &human
		h.batchGetHumanNationality(batch, human.Id)
	}

	// получение национальности людей
	br := h.pool.SendBatch(h.ctx, batch)
	defer br.Close()

	for i := 0; i < batch.Len(); i++ {
		rows, err := br.Query()
		if err != nil {
			continue
		}

		for rows.Next() {
			if err = h.setHumanNationality(rows, idHuman); err != nil {
				continue
			}
		}
	}

	humans = make([]payload.Human, 0, len(idHuman))

	for _, human := range idHuman {
		humans = append(humans, *human)
	}

	return
}

// sqlWhereFilter создает динамическое условие выборки
func (h Human) sqlWhereFilter(columnsValue map[string]interface{}) (sql string) {
	var whereColumns = make([]string, 0)

	for column, value := range columnsValue {
		switch value.(type) {
		case string:
			if len(value.(string)) != 0 {
				whereColumns = append(whereColumns, fmt.Sprintf("%s='%s'", column, value))
			}
		case int:
			if value.(int) != 0 {
				whereColumns = append(whereColumns, fmt.Sprintf("%s=%d", column, value.(int)))
			}
		}
	}

	if len(whereColumns) != 0 {
		sql += "WHERE "

		for i, column := range whereColumns {
			if i != cap(whereColumns) && cap(whereColumns)%2 == 0 {
				sql += column + " AND "
			} else {
				sql += column
			}
		}
	}

	return
}

// batchGetHumanNationality добавляет в батч SELECT запрос
// на получение национальности человека
func (h Human) batchGetHumanNationality(batch *pgx.Batch, id string) {
	sql := `
	SELECT 
	    humanid,
		countryid,
		probability
	FROM
	    humans_nationality
	WHERE 
	    humanid=@humanid
	`

	arg := pgx.NamedArgs{"humanid": id}

	batch.Queue(sql, arg)
}

// setHumanNationality присвоить человеку национальность
func (h Human) setHumanNationality(rows pgx.Rows, humans map[string]*payload.Human) error {
	var id string
	var nationality payload.Nationality

	if err := rows.Scan(&id, &nationality.CountryID, &nationality.Probability); err != nil {
		return err
	}

	human, ok := humans[id]
	if !ok {
		return fmt.Errorf("человека под этим ID не существует")
	}

	human.Nationalities = append(human.Nationalities, nationality)

	return nil
}

func (h Human) Update(id string, human payload.HumanUpdate) error {
	var batch = &pgx.Batch{}

	sql := `
	UPDATE humans
	SET
	    name=@name,
	    surname=@surname,
	    patronymic=@patronymic,
	    age=@age,
	    gender=@gender
	WHERE
	    id=@id
	`

	args := pgx.NamedArgs{
		"name":       human.Name,
		"surname":    human.Surname,
		"patronymic": human.Patronymic,
		"age":        human.Age,
		"gender":     human.Gender,
		"id":         id,
	}

	batch.Queue(sql, args)

	if len(human.AddedNationalities) != 0 {
		for _, addNationality := range human.AddedNationalities {
			h.batchAddNationality(batch, id, addNationality)
		}
	}

	if len(human.DeletedNationalities) != 0 {
		for _, deleteNationality := range human.DeletedNationalities {
			h.batchDeleteNationality(batch, id, deleteNationality)
		}
	}

	br := h.pool.SendBatch(h.ctx, batch)
	defer br.Close()

	return nil
}

// batchAddNationality добавляет в батч INSERT запрос для добавления
// человеку национальности
func (h Human) batchAddNationality(batch *pgx.Batch, humanId string, nationality payload.Nationality) {
	sql := `
	INSERT INTO 
		humans_nationality(humanid, countryid, probability) 
	VALUES(@humanid, @countryid, @probability) 
	`

	args := pgx.NamedArgs{
		"humanid":     humanId,
		"countryid":   nationality.CountryID,
		"probability": nationality.Probability,
	}

	batch.Queue(sql, args)
}

// batchDeleteNationality добавляет в батч DELETE запрос для удаления
// национальности человеку
func (h Human) batchDeleteNationality(batch *pgx.Batch, humanId string, countryId string) {
	sql := `
	DELETE FROM 
	   humans_nationality
	WHERE
	    humanid=@humanid
	AND 
	    countryid=@countryid
	`

	args := pgx.NamedArgs{
		"humanid":   humanId,
		"countryid": countryId,
	}

	batch.Queue(sql, args)
}

func (h Human) Delete(id string) error {
	sql := `
	DELETE FROM
	   humans
	WHERE 
	    id=@id
	`

	arg := pgx.NamedArgs{"id": id}

	if _, err := h.pool.Exec(h.ctx, sql, arg); err != nil {
		return err
	}

	return nil
}
