// Package service holds Fibonacci business logic. HTTP stays in package handler.
package service

import (
	"encoding/json"
	"errors"
	"fmt"
)

// MaxN is the largest count whose values still fit in uint64.
// The first n numbers are F(0)..F(n-1); F(93) is the last uint64 value, so n=94.
const MaxN = 94

const (
	MsgBadN    = "n must be a non-negative integer"
	MsgNTooBig = "n exceeds maximum"
)

var (
	ErrBadN    = errors.New(MsgBadN)
	ErrNTooBig = errors.New(MsgNTooBig)
)

var (
	fibonacciValues [MaxN + 1]uint64
	jsonByN         [MaxN + 1][]byte
)

func init() {
	fibonacciValues[1] = 1
	for i := 2; i <= MaxN; i++ {
		fibonacciValues[i] = fibonacciValues[i-1] + fibonacciValues[i-2]
	}

	for n := 0; n <= MaxN; n++ {
		body, err := json.Marshal(fibonacciValues[:n])
		if err != nil {
			panic(fmt.Sprintf("precompute json for n=%d: %v", n, err))
		}

		jsonByN[n] = body
	}
}

// Sequence returns a copy of the first n Fibonacci numbers starting at 0.
func Sequence(n int) ([]uint64, error) {
	if err := checkN(n); err != nil {
		return nil, err
	}

	out := make([]uint64, n)
	copy(out, fibonacciValues[:n])

	return out, nil
}

// JSON returns the precomputed JSON array for the first n numbers.
// The returned slice must not be modified.
func JSON(n int) ([]byte, error) {
	if err := checkN(n); err != nil {
		return nil, err
	}

	return jsonByN[n], nil
}

func checkN(n int) error {
	if n < 0 {
		return ErrBadN
	}

	if n > MaxN {
		return ErrNTooBig
	}

	return nil
}
