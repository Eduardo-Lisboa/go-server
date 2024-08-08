package main

import (
	"context"
	"errors"
	"fmt"
	"go-server/internal/api"
	"go-server/internal/store/pgstore/pgstore"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("WBS_DATABASE_USER"),
		os.Getenv("WBS_DATABASE_PASSWORD"),
		os.Getenv("WBS_DATABASE_HOST"),
		os.Getenv("WBS_DATABASE_PORT"),
		os.Getenv("WBS_DATABASE_NAME"),
	))
	if err != nil {
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		panic(err)
	}

	handler := api.NewHandler(pgstore.New(pool))

	go func() {
		if err := http.ListenAndServe(":9000", handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
