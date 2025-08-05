package utils

import (
	"fmt"

	"github.com/google/uuid"
)

type UUIDGenerator struct{}

func (g *UUIDGenerator) MustNewUUID() string {
	uuid, err := uuid.NewV7()
	if err != nil {
		panic(fmt.Errorf("failed to generate uuid: %w", err))
	}
	return uuid.String()
}
