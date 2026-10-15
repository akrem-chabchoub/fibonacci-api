package service

import "strconv"

// ParseN parses n from a query string. Only plain digits are accepted: 0 or [1-9][0-9]*.
func ParseN(in string) (int, error) {
	if in == "" || !isPlainUint(in) {
		return 0, ErrBadN
	}

	n, err := strconv.Atoi(in)
	if err != nil {
		return 0, ErrNTooBig
	}

	return n, nil
}

func isPlainUint(in string) bool {
	if in == "0" {
		return true
	}

	if in[0] < '1' || in[0] > '9' {
		return false
	}

	for i := 1; i < len(in); i++ {
		if in[i] < '0' || in[i] > '9' {
			return false
		}
	}

	return true
}
