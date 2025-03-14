package token

import (
	"testing"
	"time"
)

func TestGenToken(t *testing.T) {
	svc := NewService()

	_, err := svc.ValidateToken(svc.GenerateToken("userId", "admin", time.Now()))
	if err != nil {
		t.Error("validate token", err)
		return
	}
}
