package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS accounts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	balance_cents INTEGER NOT NULL DEFAULT 0
);
`

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	// Connection pool settings (meaningful for network DBs; still good practice).
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.ExecContext(ctx, schema); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func seed(ctx context.Context, db *sql.DB) error {
	var n int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM accounts`).Scan(&n); err != nil {
		return err
	}
	if n > 0 {
		return nil
	}
	_, err := db.ExecContext(ctx, `INSERT INTO accounts (name, balance_cents) VALUES (?, ?), (?, ?)`,
		"Alice", 10000, "Bob", 5000)
	return err
}

func transfer(ctx context.Context, db *sql.DB, from, to int64, amount int64) error {
	tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	var fromBal int64
	if err := tx.QueryRowContext(ctx, `SELECT balance_cents FROM accounts WHERE id = ?`, from).Scan(&fromBal); err != nil {
		return fmt.Errorf("from account: %w", err)
	}
	if fromBal < amount {
		return fmt.Errorf("insufficient funds: have %d need %d", fromBal, amount)
	}

	res, err := tx.ExecContext(ctx, `UPDATE accounts SET balance_cents = balance_cents - ? WHERE id = ?`, amount, from)
	if err != nil {
		return err
	}
	if aff, _ := res.RowsAffected(); aff != 1 {
		return fmt.Errorf("debit failed")
	}

	res, err = tx.ExecContext(ctx, `UPDATE accounts SET balance_cents = balance_cents + ? WHERE id = ?`, amount, to)
	if err != nil {
		return err
	}
	if aff, _ := res.RowsAffected(); aff != 1 {
		return fmt.Errorf("credit failed")
	}

	return tx.Commit()
}

func printBalances(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, `SELECT id, name, balance_cents FROM accounts ORDER BY id`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		var bal int64
		if err := rows.Scan(&id, &name, &bal); err != nil {
			return err
		}
		fmt.Printf("  #%d %s: $%.2f\n", id, name, float64(bal)/100)
	}
	return rows.Err()
}

func main() {
	path := "lab10.db"
	if v := os.Getenv("DB_PATH"); v != "" {
		path = v
	}
	defer os.Remove(path)

	db, err := openDB(path)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := seed(ctx, db); err != nil {
		log.Fatalf("seed: %v", err)
	}

	fmt.Println("Before transfer:")
	if err := printBalances(ctx, db); err != nil {
		log.Fatalf("balances: %v", err)
	}

	if err := transfer(ctx, db, 1, 2, 2500); err != nil {
		log.Fatalf("transfer: %v", err)
	}

	fmt.Println("After transfer ($25.00 from Alice to Bob):")
	if err := printBalances(ctx, db); err != nil {
		log.Fatalf("balances: %v", err)
	}

	stats := db.Stats()
	fmt.Printf("\nPool stats: Open=%d InUse=%d Idle=%d\n", stats.OpenConnections, stats.InUse, stats.Idle)
}
