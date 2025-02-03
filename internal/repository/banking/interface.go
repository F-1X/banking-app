package banking

import (
	"banking-app/internal/model"
	"context"
)

type Repository interface {
	Send(ctx context.Context, from string, to string, amount float64) error
	GetLast(ctx context.Context, count int) ([]model.Transaction, error)
	GetBalance(ctx context.Context, address string) (float64, error)


	Close()
}
