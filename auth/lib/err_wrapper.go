package lib

import "fmt"

func ErrWrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}
