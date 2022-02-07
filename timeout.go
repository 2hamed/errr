package errr

import (
	"errors"
	"os"
)

func IsTimeoutError(err error) bool {
	for {
		if os.IsTimeout(err) {
			return true
		}
		err = errors.Unwrap(err)
		if err == nil {
			return false
		}
	}
}