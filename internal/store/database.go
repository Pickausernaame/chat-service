// Code generated by ent, DO NOT EDIT.

package store

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
)

// Database is the client that holds all ent builders.
type Database struct {
	client *Client
}

// NewDatabase creates a new database based on Client.
func NewDatabase(client *Client) *Database {
	return &Database{client: client}
}

// RunInTx runs the given function f within a transaction.
// Inspired by https://entgo.io/docs/transactions/#best-practices.
// If there is already a transaction in the context, then the method uses it.
func (db *Database) RunInTx(ctx context.Context, f func(context.Context) error) (err error) {
	tx := TxFromContext(ctx)
	if tx == nil {
		tx, err = db.client.Tx(ctx)
		if err != nil {
			return fmt.Errorf("creating transaction: %v", err)
		}
	}

	defer func() {
		if er := recover(); er != nil {
			err = tx.Rollback()
			panic(er)
		}
	}()

	if err := f(NewTxContext(ctx, tx)); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			err = fmt.Errorf("%w: rolling back transaction: %v", err, rerr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %v", err)
	}

	return nil
}

func (db *Database) loadClient(ctx context.Context) *Client {
	tx := TxFromContext(ctx)
	if tx != nil {
		return tx.Client()
	}
	return db.client
}

// Exec executes a query that doesn't return rows. For example, in SQL, INSERT or UPDATE.
func (db *Database) Exec(ctx context.Context, query string, args ...interface{}) (*sql.Result, error) {
	var res sql.Result
	err := db.loadClient(ctx).driver.Exec(ctx, query, args, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Query executes a query that returns rows, typically a SELECT in SQL.
func (db *Database) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var rows sql.Rows
	err := db.loadClient(ctx).driver.Query(ctx, query, args, &rows)
	if err != nil {
		return nil, err
	}
	return &rows, nil
}

// Chat is the client for interacting with the Chat builders.
func (db *Database) Chat(ctx context.Context) *ChatClient {
	return db.loadClient(ctx).Chat
}

// Message is the client for interacting with the Message builders.
func (db *Database) Message(ctx context.Context) *MessageClient {
	return db.loadClient(ctx).Message
}

// Problem is the client for interacting with the Problem builders.
func (db *Database) Problem(ctx context.Context) *ProblemClient {
	return db.loadClient(ctx).Problem
}

// Problem is the client for interacting with the Problem builders.
func (db *Database) Job(ctx context.Context) *JobClient {
	return db.loadClient(ctx).Job
}

// Problem is the client for interacting with the Problem builders.
func (db *Database) FailedJob(ctx context.Context) *FailedJobClient {
	return db.loadClient(ctx).FailedJob
}
