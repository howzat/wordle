package wordle

import (
	"fmt"

	"github.com/pkg/errors"
)

func WrapErr(err error, format string, a ...interface{}) error {
	return errors.Wrap(err, fmt.Sprintf(format, a...))
}
