package db

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

const (
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

type Store struct {
	*Queries
	conn *pgx.Conn
}

func NewStore() *Store {
	ctx := context.Background()

	conn, _ := pgx.Connect(ctx, dbSource)

	// defer conn.Close(ctx)
	return &Store{
		conn:    conn,
		Queries: New(conn),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {

	tx, err := store.conn.Begin(ctx)

	if err != nil {
		return err
	}

	err = fn(store.Queries)
	if err != nil {
		defer tx.Rollback(ctx)
	}

	return tx.Commit(ctx)
}

type TransferParams struct {
	FromAccoundId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Transfer `json:"from_account"`
	ToAccount   Transfer `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccoundId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccoundId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// TODO update accounts 'balance

		// 검사: 계정의 잔액이 음수가 되지 않도록 함
		fromAccount, err := q.GetAccount(ctx, arg.FromAccoundId)
		if err != nil {
			return err
		}
		if fromAccount.Balance-arg.Amount <= 0 {
			return errors.New("insufficient funds")
		}

		toAccount, err := q.GetAccount(ctx, arg.ToAccountId)
		if err != nil {
			return err
		}

		_, err = q.UpdateAccount(
			ctx, UpdateAccountParams{
				ID:      arg.FromAccoundId,
				Balance: fromAccount.Balance - arg.Amount,
			})
		if err != nil {
			return err
		}

		_, err = q.UpdateAccount(
			ctx, UpdateAccountParams{
				ID:      arg.ToAccountId,
				Balance: toAccount.Balance + arg.Amount,
			})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return TransferTxResult{}, err
	}

	return result, nil

}
