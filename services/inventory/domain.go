package inventory

import "github.com/google/uuid"

type product struct {
	Id          uuid.UUID
	Title       string
	Description string
}

func NewProduct(title string, des string) *product {
	return &product{
		Id:          uuid.New(),
		Title:       title,
		Description: des,
	}
}
