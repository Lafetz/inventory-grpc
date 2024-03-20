package inventory

import "github.com/google/uuid"

type Product struct {
	Id          uuid.UUID
	Title       string
	Description string
}
