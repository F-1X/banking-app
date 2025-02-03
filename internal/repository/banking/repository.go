package banking

import (
	"banking-app/internal/model"
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

const (
	getBalanceForUpdateQuery = "SELECT balance FROM payment_system.wallets WHERE address = '%s' FOR UPDATE"
	getBalanceQuery          = "SELECT balance FROM payment_system.wallets WHERE address = '%s'"
	updateSubQuery           = "UPDATE payment_system.wallets SET balance = balance - %f WHERE address = '%s'"
	updateAddQuery           = "UPDATE payment_system.wallets SET balance = balance + %f WHERE address = '%s'"
	transationQuery          = "INSERT INTO payment_system.transactions (from_address, to_address, amount) VALUES ('%s', '%s', '%f')"
	getLastNTrasactions      = "SELECT from_address, to_address, amount FROM payment_system.transactions ORDER BY id DESC LIMIT %d"
)

type repository struct {
	db *pgx.Conn
}

func New(ctx context.Context, connString string) Repository {
	db, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("cant connect to db: %+v", err)
	}

	if db.Ping(ctx); err != nil {
		log.Fatalf("cant ping db: %+v", err)
	}

	return &repository{db: db}
}

func (r *repository) Send(ctx context.Context, from string, to string, amount float64) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	queryBalance := fmt.Sprintf(getBalanceQuery, from)
	querySub := fmt.Sprintf(updateSubQuery, amount, from)
	queryAdd := fmt.Sprintf(updateAddQuery, amount, to)

	var balance float64
	if err := tx.QueryRow(ctx, queryBalance).Scan(&balance); err != nil {
		return err
	}
	if balance < amount {
		return NotEnough
	}
	_, err = tx.Exec(ctx, querySub)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, queryAdd)
	if err != nil {
		return err
	}
	queryTransaction := fmt.Sprintf(transationQuery, from, to, amount)
	_, err = tx.Exec(ctx, queryTransaction)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *repository) GetLast(ctx context.Context, count int) ([]model.Transaction, error) {
	query := fmt.Sprintf(getLastNTrasactions, count)

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return []model.Transaction{}, err
	}
	defer rows.Close()

	transactions := make([]model.Transaction, 0, count)
	for rows.Next() {
		var t model.Transaction
		if err := rows.Scan(&t.From, &t.To, &t.Amount); err != nil {
			log.Printf("failed to scan transaction row: %+v", err)
			continue
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
func (r *repository) GetBalance(ctx context.Context, address string) (balance float64, err error) {
	query := fmt.Sprintf(getBalanceQuery, address)
	err = r.db.QueryRow(ctx, query).Scan(&balance)
	return balance, err
}

func (r *repository) Close() {
	r.db.Close(context.Background())
}
