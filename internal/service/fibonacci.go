// Package service holds Fibonacci business logic. HTTP stays in package handler.
package service

import (
	"context"
	"errors"
	"io"
	"math/big"
)

const MaxN = 1000

const (
	MsgBadN    = "n must be a non-negative integer"
	MsgNTooBig = "n exceeds maximum"
)

var (
	ErrBadN    = errors.New(MsgBadN)
	ErrNTooBig = errors.New(MsgNTooBig)
)

var decimals [MaxN][]byte

func init() {
	prev, curr := big.NewInt(0), big.NewInt(1)
	next := new(big.Int)
	buf := make([]byte, 0, 16)

	for i := 0; i < MaxN; i++ {
		buf = prev.Append(buf[:0], 10)
		decimals[i] = append([]byte(nil), buf...)
		next.Add(prev, curr)
		prev, curr, next = curr, next, prev
	}
}

func WriteJSON(ctx context.Context, w io.Writer, n int) error {
	if n < 0 {
		return ErrBadN
	}

	if n > MaxN {
		return ErrNTooBig
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if _, err := io.WriteString(w, "["); err != nil {
		return err
	}

	for i := 0; i < n; i++ {
		if i > 0 {
			if _, err := io.WriteString(w, ","); err != nil {
				return err
			}
		}

		if _, err := w.Write(decimals[i]); err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, "]")

	return err
}
