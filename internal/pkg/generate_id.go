package pkg

import "github.com/gofrs/uuid"

func GenerateID() uuid.UUID {
	return uuid.Must(uuid.NewV7())
}
