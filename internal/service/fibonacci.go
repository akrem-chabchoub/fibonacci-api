// Package service holds Fibonacci business logic. HTTP stays in package handler.
package service

import (
	"context"
	"errors"
	"io"
	"math/big"
)

const (
	// MsgBadN is the client-facing text when n is not a plain count.
	MsgBadN = "n must be a non-negative integer"
	// MsgNTooBig is the client-facing text when n does not fit in int.
	MsgNTooBig = "n exceeds maximum"
)

var (
	// ErrBadN is returned when n cannot be parsed as a count.
	ErrBadN = errors.New(MsgBadN)
	// ErrNTooBig is returned when n does not fit in int.
	ErrNTooBig = errors.New(MsgNTooBig)
)

// WriteJSON writes the first n Fibonacci numbers as a JSON array of numbers.
func WriteJSON(ctx context.Context, w io.Writer, n int) error {
	if n < 0 {
		return ErrBadN
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := io.WriteString(w, "["); err != nil {
		return err
	}

	prev, curr := big.NewInt(0), big.NewInt(1)
	next := new(big.Int)
	digits := make([]byte, 0, 16)

	for i := 0; i < n; i++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		if i > 0 {
			if _, err := io.WriteString(w, ","); err != nil {
				return err
			}
		}

		digits = prev.Append(digits[:0], 10)
		if _, err := w.Write(digits); err != nil {
			return err
		}

		next.Add(prev, curr)
		prev, curr, next = curr, next, prev
	}

	_, err := io.WriteString(w, "]")

	return err
}
