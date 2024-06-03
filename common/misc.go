package common

import (
	"github.com/google/uuid"
)

// IsValidUUID checks that u is a valid UUID4 string
func IsValidUUID(u string) bool {
    _, err := uuid.Parse(u)
    return err == nil
 }



func GetFakeAuth() *ClientAuth {
	fakeAuth := NewClientAuth(map[string]interface{}{
		"customerId": "00000000-0000-0000-0000-000000000000",
		"customerKey": "00000000-0000-0000-0000-000000000000",
	})
	fakeAuth.SwitchAuthTo(CustomerAuth)

	return fakeAuth
}
