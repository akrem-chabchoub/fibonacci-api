package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"go.uber.org/goleak"
)

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
		{name: "max", in: "94", want: 94},
		{name: "empty", in: "", wantErr: ErrBadN},
		{name: "negative", in: "-1", wantErr: ErrBadN},
		{name: "plus", in: "+5", wantErr: ErrBadN},
		{name: "float", in: "1.5", wantErr: ErrBadN},
		{name: "scientific", in: "1e2", wantErr: ErrBadN},
		{name: "letters", in: "abc", wantErr: ErrBadN},
		{name: "leading zeros", in: "007", wantErr: ErrBadN},
		{name: "too large", in: "95", wantErr: ErrNTooBig},
		{name: "overflow digits", in: "999999999999", wantErr: ErrNTooBig},
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

func TestSequence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		n       int
		want    []uint64
		wantErr error
	}{
		{name: "empty", n: 0, want: []uint64{}},
		{name: "one", n: 1, want: []uint64{0}},
		{name: "two", n: 2, want: []uint64{0, 1}},
		{name: "five", n: 5, want: []uint64{0, 1, 1, 2, 3}},
		{name: "ten", n: 10, want: []uint64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}},
		{name: "negative", n: -1, wantErr: ErrBadN},
		{name: "too large", n: 95, wantErr: ErrNTooBig},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := Sequence(tc.n)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("Sequence(%d) error = %v, want %v", tc.n, err, tc.wantErr)
				}

				return
			}

			if err != nil {
				t.Fatalf("Sequence(%d) unexpected error: %v", tc.n, err)
			}

			if len(got) != len(tc.want) {
				t.Fatalf("Sequence(%d) len = %d, want %d", tc.n, len(got), len(tc.want))
			}

			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Fatalf("Sequence(%d)[%d] = %d, want %d", tc.n, i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestSequenceMaxFitsUint64(t *testing.T) {
	t.Parallel()

	got, err := Sequence(MaxN)
	if err != nil {
		t.Fatalf("Sequence(%d) unexpected error: %v", MaxN, err)
	}

	const f93 = uint64(12200160415121876738)
	if got[MaxN-1] != f93 {
		t.Fatalf("F(93) = %d, want %d", got[MaxN-1], f93)
	}
}

func TestJSONMatchesMarshal(t *testing.T) {
	t.Parallel()

	for n := 0; n <= MaxN; n++ {
		got, err := JSON(n)
		if err != nil {
			t.Fatalf("JSON(%d) unexpected error: %v", n, err)
		}

		seq, err := Sequence(n)
		if err != nil {
			t.Fatalf("Sequence(%d) unexpected error: %v", n, err)
		}

		want, err := json.Marshal(seq)
		if err != nil {
			t.Fatalf("json.Marshal(%d): %v", n, err)
		}

		if !bytes.Equal(got, want) {
			t.Fatalf("JSON(%d) = %s, want %s", n, got, want)
		}
	}
}

func TestJSONErrors(t *testing.T) {
	t.Parallel()

	if _, err := JSON(-1); !errors.Is(err, ErrBadN) {
		t.Fatalf("JSON(-1) error = %v, want %v", err, ErrBadN)
	}

	if _, err := JSON(MaxN + 1); !errors.Is(err, ErrNTooBig) {
		t.Fatalf("JSON(%d) error = %v, want %v", MaxN+1, err, ErrNTooBig)
	}
}

func TestJSONFive(t *testing.T) {
	t.Parallel()

	got, err := JSON(5)
	if err != nil {
		t.Fatal(err)
	}

	if string(got) != "[0,1,1,2,3]" {
		t.Fatalf("JSON(5) = %s, want [0,1,1,2,3]", got)
	}
}

func TestConcurrentJSON(t *testing.T) {
	t.Parallel()

	const workers = 32

	errCh := make(chan error, workers)
	for range workers {
		go func() {
			_, err := JSON(93)
			errCh <- err
		}()
	}

	for range workers {
		if err := <-errCh; err != nil {
			t.Fatal(err)
		}
	}
}
