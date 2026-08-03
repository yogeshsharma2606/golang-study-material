# Lab 10: Databases

Use `database/sql` with SQLite: schema setup, queries, transactions, and connection pool configuration.

## Objectives

- Open a database with `sql.Open` and verify with `Ping`.
- Run DDL and DML with context-aware APIs.
- Implement a money transfer inside a transaction with rollback on failure.
- Inspect `db.Stats()` for pool metrics.

## Setup

```bash
cd labs/10-databases
go mod tidy
go run .
```

Optional: `DB_PATH=./my.db go run .` (default file is `lab10.db`, removed on exit).

## Exercises

1. Run the program and note balances before and after the transfer.
2. Change the transfer amount to more than Alice has; observe the error and unchanged balances.
3. Log `db.Stats()` before and after a few concurrent goroutines each running a read query (extend the code).
4. Compare `BeginTx` with `Begin` and when to pass `TxOptions`.

## What to take away

- `database/sql` is a concurrency-safe pool; `sql.DB` is long-lived.
- Always use `context.Context` on queries in servers.
- Transactions group work; defer `Rollback` until `Commit` succeeds.

## Cleanup

The demo deletes `lab10.db` on exit. Remove any custom `DB_PATH` file manually if you keep it.

## Related Modules

- SQL and persistence modules in the course curriculum.
- Error handling and context cancellation.
