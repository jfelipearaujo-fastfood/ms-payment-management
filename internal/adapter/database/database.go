package database

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"

	"github.com/jfelipearaujo-org/ms-payment-management/internal/environment"
	"github.com/jfelipearaujo-org/ms-payment-management/internal/shared/health"
)

type DatabaseService interface {
	GetInstance() *sql.DB
	health.HealthCheck
}

type Service struct {
	Client *sql.DB
}

func NewDatabase(config *environment.Config) DatabaseService {
	client, err := sql.Open("postgres", config.DbConfig.Url)
	if err != nil {
		panic(fmt.Errorf("error on connect to database: %v", err))
	}

	return &Service{
		Client: client,
	}
}

func (s *Service) GetInstance() *sql.DB {
	return s.Client
}

func (s *Service) Health() *health.HealthStatus {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.Client.PingContext(ctx); err != nil {
		slog.Error("could not ping the database", "error", err)
		return &health.HealthStatus{
			Status: "unhealthy",
			Err:    err.Error(),
		}
	}

	return &health.HealthStatus{
		Status: "healthy",
	}
}
