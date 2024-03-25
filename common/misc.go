package common

import (
	"github.com/google/uuid"
)

// IsValidUUID checks that u is a valid UUID4 string
func IsValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
 }
