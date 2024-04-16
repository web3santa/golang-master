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

	return nil
}

func main() {
	// generate randome accouts and balance

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
		FromAccoundId: 2,
		ToAccountId:   1,
		Amount:        29,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)

}
