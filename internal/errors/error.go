package errors

import "errors"

var (
	ErrEmptyTitle    error = errors.New("title is empty")
	ErrTitleNotExist error = errors.New("the title doesn't exist")
)
