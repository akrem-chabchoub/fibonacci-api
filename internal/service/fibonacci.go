// Package service holds Fibonacci business logic. HTTP stays in package handler.
package service

import (
	"context"
	"errors"
	"io"
	"math/big"
)

const MaxN = 10000

const (
	MsgBadN    = "n must be a non-negative integer"
	MsgNTooBig = "n exceeds maximum"
)

var (
	ErrBadN    = errors.New(MsgBadN)
	ErrNTooBig = errors.New(MsgNTooBig)
)

var (
	jsonPrefix []byte
	endAt      [MaxN + 1]int
)

func init() {
	prev, curr := big.NewInt(0), big.NewInt(1)
	next := new(big.Int)

	buf := make([]byte, 0, MaxN*MaxN/9+MaxN+1)
	buf = append(buf, '[')
	endAt[0] = 1

	tmp := make([]byte, 0, 32)
	for i := 0; i < MaxN; i++ {
		if i > 0 {
			buf = append(buf, ',')
		}

		tmp = prev.Append(tmp[:0], 10)
		buf = append(buf, tmp...)
		endAt[i+1] = len(buf)

		next.Add(prev, curr)
		prev, curr, next = curr, next, prev
	}

	jsonPrefix = buf
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

	if _, err := w.Write(jsonPrefix[:endAt[n]]); err != nil {
		return err
	}

	_, err := io.WriteString(w, "]")

	return err
}
