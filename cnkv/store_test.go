package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBasic(t *testing.T) {
	err := Put("bob", "abc")
	require.NoError(t, err)

	val, err := Get("bob")
	require.NoError(t, err)
	require.Equal(t, "abc", val)

	err = Delete("bob")
	require.NoError(t, err)

	_, err = Get("bob")
	require.Equal(t, ErrNoSuchKey, err)
}
