package goose

import (
	"errors"
	"strings"
)

func prependErrors(errs []error, err error) []error {
	return append([]error{err}, errs...)
}

func errorsToError(errs []error) error {
	msg := make([]string, len(errs))
	for _, err := range errs {
		msg = append(msg, err.Error())
	}
	return errors.New(strings.Join(msg, ", "))
}