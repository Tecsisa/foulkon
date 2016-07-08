package postgresql

import (
	"fmt"
	"os"
	"testing"

	"github.com/tecsisa/authorizr/database"
	"time"
)

var repoDB PostgresRepo

func TestMain(m *testing.M) {
	// Wait for DB
	time.Sleep(5 * time.Second)
	// Retrieve db connector to run test
	dbmap, err := InitDb("postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		fmt.Fprintln(os.Stderr, "There was an error starting connector", err)
		os.Exit(1)
	}
	repoDB = PostgresRepo{
		Dbmap: dbmap,
	}

	result := m.Run()

	os.Exit(result)
}

// User Table aux methods
func insertUser(id string, externalID string, path string, createAt int64, urn string) error {
	err := repoDB.Dbmap.Exec("INSERT INTO public.users (id, external_id, path, create_at, urn) VALUES (?, ?, ?, ?, ?)",
		id, externalID, path, createAt, urn).Error

	// Error handling
	if err != nil {
		return &database.Error{
			Code:    database.INTERNAL_ERROR,
			Message: err.Error(),
		}
	}
	return nil
}

func getUsersCountFiltered(id string, externalID string, path string, createAt int64, urn string) (int, error) {
	query := repoDB.Dbmap.Table(User{}.TableName())
	if id != "" {
		query = query.Where("id = ?", id)
	}
	if path != "" {
		query = query.Where("path = ?", path)
	}
	if createAt != 0 {
		query = query.Where("create_at = ?", createAt)
	}
	if urn != "" {
		query = query.Where("urn = ?", urn)
	}
	var number int
	if err := query.Count(&number).Error; err != nil {
		return 0, err
	}

	return number, nil
}

func cleanUserTable() error {
	if err := repoDB.Dbmap.Delete(&User{}).Error; err != nil {
		return err
	}
	return nil
}
