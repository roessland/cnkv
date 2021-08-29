package main

import "errors"

var ErrNoSuchKey = errors.New("no such key")

var store = make(map[string]string)

func Put(key string, value string) error {
	store[key] = value

	return nil
}

func Get(key string) (string, error) {
	value, ok := store[key]

	if !ok {
		return "", ErrNoSuchKey
	}

	return value, nil
}

func Delete(key string) error {
	delete(store, key)

	return nil
}