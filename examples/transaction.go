package examples

import (
	"context"
	"fmt"

	"github.com/pakkasys/fluidapi/database"
)

func main() {
	cfg := database.ConnectConfig{
		Driver:   database.SQLite3,
		Database: "example.db",
	}

	dsn, _ := database.DSN(cfg)
	db, err := database.Connect(cfg, database.NewSQLDB, *dsn)
	if err != nil {
		panic(err)
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	result, err := database.Transaction(context.Background(), tx, func(ctx context.Context, tx database.Tx) (int64, error) {
		res, err := tx.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, name TEXT)")
		if err != nil {
			return 0, err
		}
		res, err = tx.Exec("INSERT INTO users (name) VALUES (?)", "Bob")
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Inserted user with ID:", result)
}
