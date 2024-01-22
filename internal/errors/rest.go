package errors

const (
	ErrJSONDecode           = "Ошибка декодирования JSON"
	ErrExcessPaginationSize = "Превышен лимит размера пагинации"
)

type ErrDataNotFound struct {
}

func (e *ErrDataNotFound) Error() string {
	return "Данные отсутствуют"
}
