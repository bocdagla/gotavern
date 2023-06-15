package command

import "github.com/google/uuid"

type Create struct {
	Holder      uuid.UUID `json:holder`
	Description string    `json:description`
}
