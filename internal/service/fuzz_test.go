package service

import (
	"errors"
	"testing"
)

func FuzzParseN(f *testing.F) {
	for _, seed := range []string{"", "0", "5", "-1", "1.5", "93", "94", "95", "999999999999", "007", "+5", "1e2"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, in string) {
		n, err := ParseN(in)
		if err != nil {
			if !errors.Is(err, ErrBadN) && !errors.Is(err, ErrNTooBig) {
				t.Fatalf("ParseN(%q) unexpected error: %v", in, err)
			}

			return
		}

		if n < 0 {
			t.Fatalf("ParseN(%q) = %d, want >= 0", in, n)
		}
	})
}
