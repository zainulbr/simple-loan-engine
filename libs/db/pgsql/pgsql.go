package pgsql

import (
	"fmt"
	"log"
	"sync"

	"github.com/go-pg/pg/v10"
	"github.com/zainulbr/simple-loan-engine/settings"
)

const (
	defaultKey = "default"
)

var store sync.Map

// Open opens a postgres connection with the specified URI,
// then saved to the specified key. If key is not specified,
// use default key. Access with DB(key) function.
func open(option *settings.PostgresOption) (*pg.DB, error) {
	pgOption, err := pg.ParseURL(option.URI)
	if err != nil {
		log.Println("pg parseurl", err)
		return nil, err
	}

	v := pg.Connect(pgOption)

	version := ""
	if _, err := v.QueryOne(pg.Scan(&version), "SELECT version()"); err != nil {
		log.Println("pg version", err)
		return nil, err
	}
	log.Println(version)

	return v, nil
}

// Close closes connections of specified keys,
func Close() error {
	store.Range(func(key, value interface{}) bool {
		log.Printf("Closing SQL DB: %s", key)

		if err := (value.(*pg.DB)).Close(); err != nil {
			log.Println(err)
		}

		store.Delete(key)

		return true
	})

	return nil
}

// Create connection
func Create(option *settings.PostgresOption, key string) error {
	if _, ok := store.Load(key); ok {
		return fmt.Errorf("SQL DB '%s' already exists", key)
	}

	log.Printf("Initializing SQL DB: %s", key)

	db, err := open(option)
	if err != nil {
		return err
	}

	store.Store(key, db)

	return nil
}

// Open connection
func Open(settings *settings.Settings) error {
	return Create(&settings.Conn.Postgres, defaultKey)
}

// DB returns the database connection of the specified key,
// if none specified, return default connection
func DB(dbKey ...string) *pg.DB {
	key := defaultKey

	if len(dbKey) > 0 && len(dbKey[0]) > 0 {
		key = dbKey[0]
	}

	instance, ok := store.Load(key)
	if !ok {
		log.Fatalf("SQL DB '%s' not found, please call Create() or Open() first.", key)
	}

	return instance.(*pg.DB)
}

func init() {
	store = sync.Map{}
}
