package user

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin          UserRole = "admin"
	RoleInvestor       UserRole = "investor"
	RoleBorrower       UserRole = "borrower"
	RoleFiledValidator UserRole = "field_validator"
	RoleFieldOfficer   UserRole = "field_officer"
)

type User struct {
	UserId    uuid.UUID `json:"user_id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      UserRole  `json:"role,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
