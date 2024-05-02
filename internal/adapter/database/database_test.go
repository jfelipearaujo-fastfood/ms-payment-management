package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/environment"
	"github.com/stretchr/testify/assert"
)

func TestGetInstance(t *testing.T) {
	t.Run("Should initialize the database", func(mt *testing.T) {
		// Arrange
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := NewDatabase(config)
		service.(*Service).Client = db

		// Act
		res := service.GetInstance()

		// Assert
		assert.NotNil(mt, res)
	})
}

func TestHealth(t *testing.T) {
	t.Run("Should return healthy status", func(mt *testing.T) {
		// Arrange
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := NewDatabase(config)
		service.(*Service).Client = db

		// Act
		res := service.Health()

		// Assert
		assert.NotNil(mt, res)
		assert.Equal(mt, "healthy", res.Status)
	})

	t.Run("Should return unhealthy status", func(mt *testing.T) {
		// Arrange
		db, _, err := sqlmock.New()
		if err != nil {
			t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		}
		defer db.Close()

		config := &environment.Config{
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
		}

		service := NewDatabase(config)

		// Act
		res := service.Health()

		// Assert
		assert.NotNil(mt, res)
		assert.Equal(mt, "unhealthy", res.Status)
	})
}
