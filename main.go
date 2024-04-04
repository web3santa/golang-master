package main

import (
	"context"
	db "golang-master/sqlc"
	"golang-master/util"
	"log"

	"github.com/jackc/pgx/v5"
)

const (
	dbSource = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func run() error {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbSource)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	queries := db.New(conn)

	// list all accounts
	// accounts, err := queries.ListAccounts(ctx, db.ListAccountsParams{})
	// if err != nil {
	// 	return err
	// }
	// log.Println(accounts)

	for i := 0; i < 10; i++ {
		arg := db.CreateAccountParams{
			Owner:    util.RandomOwner(),
			Balance:  util.RandomMoney(),
			Currency: util.RandomeCurrenyCy(),
		}
		// create an accounts
		insertedAccount, err := queries.CreateAccount(ctx, arg)
		if err != nil {
			return err
		}
		log.Println(insertedAccount)

	}

	// update an accounts
	// updateAccount, err := queries.UpdateAccount(ctx, db.UpdateAccountParams{2, 250})
	// if err != nil {
	// 	return err
	// }
	// log.Println(updateAccount)

	// get a account
	// account, err := queries.GetAccount(ctx, 1)
	// if err != nil {
	// 	return err
	// }

	// delete a account
	// if err := queries.DeleteAccount(ctx, 1); err != nil {
	// 	return err
	// }

	// // prints true
	// log.Println(account)

	return nil
}

func main() {
	// if err := run(); err != nil {
	// 	log.Fatal(err)
	// }

	ctx := context.Background()

	// db 연결 생성
	conn, err := pgx.Connect(ctx, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	store := db.NewStore()

	// SendMoney 함수 호출
	result, err := store.TransferTx(ctx, db.TransferParams{
		FromAccoundId: 1,
		ToAccountId:   2,
		Amount:        100,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)

}
