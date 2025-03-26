package token

import (
	"fmt"
	"testing"
	"time"
)

func TestGenToken(t *testing.T) {
	roles := map[string]string{
		"": "",
		// "e957a5ad-d6be-4ed6-a328-d2d9cad45955": "investor",
		// "99385f18-5655-432c-a383-163e9e634b0d": "investor",
		// "f4fcbdce-d833-48d1-bc98-af05aa29ad1d": "borrower",
		// "e6df5aef-090b-4ec0-974f-03f58abc00b2": "field_validator",
		// "c76a1d43-5758-4e05-9f08-a16e991bc656": "field_officer",
	}

	svc := NewService()
	for k, v := range roles {
		token := svc.GenerateToken(k, v, time.Now().Add(time.Hour*24)) // 1 day
		fmt.Println(v, token)
		_, err := svc.ValidateToken(token)
		if err != nil {
			t.Error("validate token", err)
			return
		}

	}
}
