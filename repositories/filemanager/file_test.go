package filemanager

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zainulbr/simple-loan-engine/libs/db/pgsql"
	"github.com/zainulbr/simple-loan-engine/models/filemanager"
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

func TestFile(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	data := &filemanager.File{
		Label:        "testlabel",
		Location:     "testlocation",
		LocationType: "testlocationtype",
		FileType:     "testfiletype",
	}

	svc := NewFileRepository(db)

	err := svc.Create(context.Background(), data)
	assert.NoError(t, err)
	assert.NotEmpty(t, data.FileID)
	t.Logf("Created  ID: %s", data.FileID)

	resp, err := svc.GetByID(context.Background(), data.FileID)
	assert.NoError(t, err)
	assert.Equal(t, resp.Label, data.Label)
	assert.Equal(t, resp.Location, data.Location)
	assert.Equal(t, resp.LocationType, data.LocationType)
	assert.Equal(t, resp.FileType, data.FileType)
	t.Logf("Retrieved : %+v", resp)

	err = svc.Delete(context.Background(), data.FileID)
	assert.NoError(t, err)

	_, err = svc.GetByID(context.Background(), data.FileID)
	assert.Error(t, err) // Should be error

	t.Log("Data deleted successfully")
}
