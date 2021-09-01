package main

import "context"

type TransactionLogger interface {
	WriteDelete(key string)
	WritePut(key, value string)
	Err() <-chan error
	ReadEvents(context.Context) (<-chan Event, <-chan error)
	Run(ctx context.Context)
}

type Event struct {
	Sequence  uint64 // A unique record ID
	EventType EventType // The action taken
	Key       string // The key affected by this transaction
	Value     string // The value of a PUT in the transaction
}

type EventType byte

const (
	_ EventType = iota // skip 0
	EventDelete
	EventPut
)