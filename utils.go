package payrail

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// GenerateReference returns a unique reference starting with "PAYRAIL-"
func GenerateReference() string {
	return fmt.Sprintf("PAYRAIL-%s", uuid.New().String())
}

// ValidatePhone returns network name ("MTN", "Airtel", "Zamtel")
// or an error if the number is not recognized as a valid Zambian mobile.
func ValidatePhone(phone string) (string, error) {
	// Normalize: remove spaces, dashes, plus sign
	num := strings.ReplaceAll(phone, " ", "")
	num = strings.ReplaceAll(num, "-", "")
	num = strings.TrimPrefix(num, "+")

	// Remove Zambia country code if present
	if strings.HasPrefix(num, "260") {
		num = num[3:]
	}

	// Remove leading zero if present
	if strings.HasPrefix(num, "0") {
		num = num[1:]
	}

	// Must have at least prefix + subscriber (min ~9 digits)
	if len(num) < 9 {
		return "", errors.New("invalid phone number length")
	}

	// Extract operator prefix (first 2 digits)
	prefix := num[:2]

	switch prefix {
	case "76", "96":
		return "mtn", nil
	case "77", "97", "55":
		return "airtel", nil
	case "75", "95":
		return "zamtel", nil
	default:
		return "", errors.New("unknown or unsupported network prefix")
	}
}
