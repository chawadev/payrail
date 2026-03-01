package core

import (
	"fmt"

	"github.com/google/uuid"
)

// GenerateReference returns a unique reference starting with "PAYRAIL-"
func GenerateReference() string {
	return fmt.Sprintf("PAYRAIL-%s", uuid.New().String())
}
