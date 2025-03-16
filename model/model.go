package model

type Product struct {
	ID          string  `json:"id" faker:"uuid_hyphenated"`
	Name        string  `json:"name" faker:"first_name"`
	Price       float64 `json:"price" faker:"amount"`
	Description string  `json:"description" faker:"paragraph len=800"`
	Score       float64 `json:"score" faker:"boundary_start=-2, boundary_end=2"`
}
