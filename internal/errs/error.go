package errs

import "errors"

var (
	ErrEmptyVoteTitle      error = errors.New("title is empty")
	ErrEmptyChoiceTitle    error = errors.New("title is empty")
	ErrTitleNotExist       error = errors.New("the title doesn't exist")
	ErrChoiceTitleNotExist error = errors.New("the choice title doesn't exist")
	ErrTitleAlreadyExist   error = errors.New("title adready exist")
)
