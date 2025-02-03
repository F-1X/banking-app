package payment_system

import (
	"context"
	"fmt"
	"banking-app/internal/model"
	"banking-app/internal/repository/banking"
)

type service struct {
	repo banking.Repository
}

func New(ctx context.Context, repo banking.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Send(ctx context.Context, from string, to string, amount float64) error {
	if from == "" || to == "" || amount <= 0 {
		return fmt.Errorf("bad body parameters. write correct 'from', 'to', 'amount' parameters")
	}
	return s.repo.Send(ctx, from, to, amount)
}

func (s *service) GetLast(ctx context.Context, count int) ([]model.Transaction, error) {
	if count <= 0 {
		return []model.Transaction{}, fmt.Errorf("bad count parameter")
	}
	return s.repo.GetLast(ctx, count)
}

func (s *service) GetBalance(ctx context.Context, address string) (float64, error) {
	return s.repo.GetBalance(ctx, address)
}
