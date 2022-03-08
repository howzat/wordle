package words

import (
	"fmt"

	"github.com/pkg/errors"
)

func WrapErr(err error, message string, args ...interface{}) error {
	return errors.Wrap(err, fmt.Sprintf(message, args...))
}
