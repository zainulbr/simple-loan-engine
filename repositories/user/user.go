package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/models/user"
)

type UserService interface {
	Create(ctx context.Context, user user.User) (uuid.UUID, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*user.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
}
