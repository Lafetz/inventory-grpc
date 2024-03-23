package inventory

import "errors"

var (
	ErrProductNotFound = errors.New("product not found")
	ErrInvalidId       = errors.New("the id is invalid")
)
