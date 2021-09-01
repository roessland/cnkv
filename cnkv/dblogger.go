package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
)

// TODO: close everything on exit
// TODO: don't lose events on exit
// Add log compaction

type PostgresTransactionLogger struct {
	events chan<- Event // Write only channel for sending events
	errors <-chan error // Read only channel for receiving errors
	db     *pgx.Conn    // The location of the transaction log
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.events <- Event{EventType: EventDelete, Key: key}
}

func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.events <- Event{EventType: EventPut, Key: key, Value: value}
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func NewPostgresTransactionLogger(ctx context.Context) (TransactionLogger, error) {
	connString := "postgres://exampleuser:examplepass@localhost:5432/exampledb"

	var err error
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to DB: %w", err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot ping DB: %w", err)
	}

	logger := &PostgresTransactionLogger{db: conn}

	err = logger.createTableIfNotExists(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot create table: %w", err)
	}

	return logger, nil
}

func (l *PostgresTransactionLogger) Run(ctx context.Context) {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	// todo handle context.close
	go func() {
		query := `insert into transactions (event_type, key, value) values ($1, $2, $3)`
		for e := range events {
			_, err := l.db.Exec(ctx, query, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
				return
			}
		}
	}()
}

func (l *PostgresTransactionLogger) ReadEvents(ctx context.Context) (<-chan Event, <-chan error) {
	query := `select sequence, event_type, key, value from transactions order by sequence`

	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		rows, err := l.db.Query(ctx, query)
		if err != nil {
			outError <- fmt.Errorf("db query error: %w", err)
			return
		}

		defer rows.Close()

		var e Event

		for rows.Next() {
			err := rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("db scan error: %w", err)
				return
			}

			outEvent <- e
		}

		if err := rows.Err(); err != nil {
			outError <- fmt.Errorf("db error: %w", err)
			return
		}
	}()

	return outEvent, outError
}

func (l *PostgresTransactionLogger) createTableIfNotExists(ctx context.Context) error {
	query := `
		create table if not exists 
		transactions (
			sequence serial primary key,
			event_type int not null,
			key varchar(255) not null,
			value text not null default ''
		)
	`
	_, err := l.db.Exec(ctx, query)
	return err
}
