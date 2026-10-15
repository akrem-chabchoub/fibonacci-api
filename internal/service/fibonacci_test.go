package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"testing"

	"go.uber.org/goleak"
)

func jsonBytes(n int) ([]byte, error) {
	var buf bytes.Buffer
	err := WriteJSON(context.Background(), &buf, n)
	return buf.Bytes(), err
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestParseN(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      string
		want    int
		wantErr error
	}{
		{name: "zero", in: "0", want: 0},
		{name: "one", in: "1", want: 1},
		{name: "five", in: "5", want: 5},
		{name: "past uint64 count", in: "95", want: 95},
		{name: "max", in: "1000", want: 1000},
		{name: "empty", in: "", wantErr: ErrBadN},
		{name: "negative", in: "-1", wantErr: ErrBadN},
		{name: "plus", in: "+5", wantErr: ErrBadN},
		{name: "float", in: "1.5", wantErr: ErrBadN},
		{name: "scientific", in: "1e2", wantErr: ErrBadN},
		{name: "letters", in: "abc", wantErr: ErrBadN},
		{name: "leading zeros", in: "007", wantErr: ErrBadN},
		{name: "above max", in: "1001", wantErr: ErrNTooBig},
		{name: "overflow digits", in: "18446744073709551616", wantErr: ErrNTooBig},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseN(tc.in)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("ParseN(%q) error = %v, want %v", tc.in, err, tc.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("ParseN(%q) unexpected error: %v", tc.in, err)
			}

			if got != tc.want {
				t.Fatalf("ParseN(%q) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

func TestJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		n       int
		want    string
		wantErr error
	}{
		{name: "empty", n: 0, want: "[]"},
		{name: "one", n: 1, want: "[0]"},
		{name: "two", n: 2, want: "[0,1]"},
		{name: "five", n: 5, want: "[0,1,1,2,3]"},
		{name: "ten", n: 10, want: "[0,1,1,2,3,5,8,13,21,34]"},
		{name: "negative", n: -1, wantErr: ErrBadN},
		{name: "too large", n: MaxN + 1, wantErr: ErrNTooBig},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := jsonBytes(tc.n)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("JSON(%d) error = %v, want %v", tc.n, err, tc.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("JSON(%d) unexpected error: %v", tc.n, err)
			}

			if string(got) != tc.want {
				t.Fatalf("JSON(%d) = %s, want %s", tc.n, got, tc.want)
			}
		})
	}
}

func TestJSONPastUint64(t *testing.T) {
	t.Parallel()

	got, err := jsonBytes(95)
	if err != nil {
		t.Fatal(err)
	}

	var seq []json.Number
	if err := json.Unmarshal(got, &seq); err != nil {
		t.Fatal(err)
	}

	if len(seq) != 95 {
		t.Fatalf("len = %d, want 95", len(seq))
	}

	if seq[93].String() != "12200160415121876738" {
		t.Fatalf("F(93) = %s", seq[93])
	}

	if seq[94].String() != "19740274219868223167" {
		t.Fatalf("F(94) = %s", seq[94])
	}
}

func TestJSONMatchesRecurrence(t *testing.T) {
	t.Parallel()

	for n := 0; n <= MaxN; n++ {
		got, err := jsonBytes(n)
		if err != nil {
			t.Fatalf("JSON(%d) unexpected error: %v", n, err)
		}

		var seq []json.Number
		if err := json.Unmarshal(got, &seq); err != nil {
			t.Fatalf("JSON(%d) unmarshal: %v", n, err)
		}

		if len(seq) != n {
			t.Fatalf("JSON(%d) len = %d", n, len(seq))
		}

		prev, curr := big.NewInt(0), big.NewInt(1)
		for i, num := range seq {
			want := prev.String()
			if num.String() != want {
				t.Fatalf("JSON(%d)[%d] = %s, want %s", n, i, num, want)
			}

			next := new(big.Int).Add(prev, curr)
			prev, curr = curr, next
		}
	}
}

func TestWriteJSONCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := WriteJSON(ctx, io.Discard, 50)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("WriteJSON() error = %v, want context.Canceled", err)
	}
}

func TestConcurrentJSON(t *testing.T) {
	t.Parallel()

	const workers = 32

	errCh := make(chan error, workers)
	for range workers {
		go func() {
			errCh <- WriteJSON(context.Background(), io.Discard, 80)
		}()
	}

	for range workers {
		if err := <-errCh; err != nil {
			t.Fatal(err)
		}
	}
}
