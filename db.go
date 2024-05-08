package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	db  *pgxpool.Pool
	ctx = context.Background()
)

func initDB() (err error) {
	db, err = pgxpool.New(ctx, "postgres://postgres:password3@localhost:5432/task?")
	if err != nil {
		fmt.Println("Error creating connection!", err.Error())
		return
	}

	err = db.Ping(ctx)
	if err != nil {
		log.Println("Error in ping", err.Error())
		return
	}

	return
}
