package payload

type (
	// AgifyAPI структура для получения возраста с API
	// https://agify.io/
	AgifyAPI struct {
		Age int64 `json:"age"`
	}

	// GenderizeAPI структура для получения пола с API
	// https://genderize.io/
	GenderizeAPI struct {
		Gender string `json:"gender"`
	}

	// NationalizeAPI структура для получения национальности с API
	// https://nationalize.io/
	NationalizeAPI struct {
		Country []Nationality `json:"country"`
	}
	// Nationality структура национальности человека
	Nationality struct {
		CountryID   string  `json:"country_id"`
		Probability float32 `json:"probability"`
	}
)
