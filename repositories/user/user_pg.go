package user

import (
	"context"

	"github.com/zainulbr/simple-loan-engine/models/user"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
)

type userRepository struct {
	db *pg.DB
}

func NewuserRepository(db *pg.DB) UserRepository {
	return &userRepository{db: db}
}

// Create User
func (r *userRepository) Create(ctx context.Context, user user.User) (uuid.UUID, error) {
	var userId uuid.UUID
	_, err := r.db.QueryOneContext(ctx,
		pg.Scan(&userId),
		`INSERT INTO "user".users (email, role) VALUES (?, ?) RETURNING user_id`,
		user.Email, user.Role)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

// Get User by ID
func (r *userRepository) GetByID(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	user := new(user.User)
	_, err := r.db.QueryOneContext(ctx,
		user,
		`SELECT user_id, email, role, created_at, updated_at FROM "user".users WHERE user_id = ?`,
		userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Delete User
func (r *userRepository) Delete(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM "user".users WHERE user_id = ?`, userID)
	return err
}
