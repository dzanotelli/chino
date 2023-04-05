package common

import (
	"fmt"

	"github.com/google/uuid"
)

// IsValidUUID checks that u is a valid UUID4 string
func IsValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
 }

// JoinErrors combines a list of errors in a single one
// (needed until it's not possible to upgrade to go 1.20).
// label is a string that is preponed to the errors
func JoinErrors(label string, errors []error) error {
    err := ""
    for i, e := range errors {
        if i == 0 {
            err = fmt.Sprintf("%s %v", label, e)
            continue
        }
        err = fmt.Sprintf("%s; %v", err, e)
    }
    return fmt.Errorf(err)
}