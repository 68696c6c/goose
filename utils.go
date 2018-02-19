package goose

import (
	"errors"
	"strings"
	"golang.org/x/sys/unix"
	"os"
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

func writableDir(path string) bool {
	return unix.Access(path, unix.W_OK) == nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
