package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zainulbr/simple-loan-engine/libs/db/pgsql"
	"github.com/zainulbr/simple-loan-engine/models/user"
	"github.com/zainulbr/simple-loan-engine/settings"

	"github.com/go-pg/pg/v10"
)

func setupTestDB(t *testing.T) *pg.DB {
	config := settings.Settings{
		Conn: settings.ConnectionSettings{
			Postgres: settings.PostgresOption{
				URI:     "postgres://admin:admin@localhost:5432/loan-db?sslmode=disable",
				Enabled: true,
			},
		},
	}

	if err := pgsql.Open(&config); err != nil {
		t.Error(err)
	}

	return pgsql.DB()
}

func TestUser(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	data := user.User{
		Email: "testuser@example.com",
		Role:  user.RoleBorrower,
	}
	svc := NewuserRepository(db)

	userID, err := svc.Create(context.Background(), data)
	assert.NoError(t, err)
	assert.NotEmpty(t, userID)
	t.Logf("Created User ID: %s", userID)

	respUser, err := svc.GetByID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Equal(t, respUser.Email, data.Email)
	assert.Equal(t, respUser.Role, data.Role)

	t.Logf("Retrieved User: %+v", respUser)

	err = svc.Delete(context.Background(), userID)
	assert.NoError(t, err)

	_, err = svc.GetByID(context.Background(), userID)
	assert.Error(t, err) // Should be error

	t.Log("User deleted successfully")
}
