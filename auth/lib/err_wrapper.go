package lib

import "fmt"

func ErrWrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
