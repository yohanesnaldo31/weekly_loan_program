package app

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"weekly_loan_program/repo/cache"
	"weekly_loan_program/repo/db"
	"weekly_loan_program/server/httphandler"
	loanHandler "weekly_loan_program/server/httphandler/loan"
	"weekly_loan_program/service/loan"
	loanUC "weekly_loan_program/usecase/loan"
)

type Application struct {
	HTTPHandler httphandler.HTTPHandler
}

func NewApplication(ctx context.Context) (*Application, error) {
	dbPool, err := NewDatabase(ctx)
	if err != nil {
		return nil, fmt.Errorf("app: init database: %w", err)
	}

	redisClient, err := NewRedis(ctx)
	if err != nil {
		return nil, fmt.Errorf("app: init redis: %w", err)
	}

	// initialize repositories
	dbRepo := db.NewRepository(dbPool)
	cacheRepo := cache.NewRepository(redisClient)

	// initialize service/domain layers
	loanService := loan.NewService(dbRepo, cacheRepo)

	// initialize usecase/business logic layers
	loanUsecase := loanUC.NewUsecase(loanService)

	// initialize handler layers -> validating user inpnut & building response
	loanHandler := loanHandler.NewHandler(loanUsecase)

	return &Application{
		HTTPHandler: *httphandler.InitHTTPHandler(loanHandler),
	}, nil
}

// NewDatabase opens a connection pool to the Postgres instance defined in docker-compose.yml
func NewDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	dsn := "postgres://postgres:postgres@localhost:5432/loan?sslmode=disable"

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("db: create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("db: ping: %w", err)
	}

	return pool, nil
}

// NewRedis opens a client connection to the Redis instance defined in docker-compose.yml
func NewRedis(ctx context.Context) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	if err := client.Ping(ctx).Err(); err != nil {
		client.Close()
		return nil, fmt.Errorf("redis: ping: %w", err)
	}

	return client, nil
}
