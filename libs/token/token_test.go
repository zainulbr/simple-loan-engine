package token

import (
	"fmt"
	"testing"
	"time"
)

func TestGenToken(t *testing.T) {
	roles := []string{
		"investor",
		"borrower",
		"field_validator",
		"field_officer",
	}

	svc := NewService()
	for _, v := range roles {
		token := svc.GenerateToken("f4fcbdce-d833-48d1-bc98-af05aa29ad1d", v, time.Now())
		fmt.Println(v, token)
		_, err := svc.ValidateToken(token)
		if err != nil {
			t.Error("validate token", err)
			return
		}

	}
}
