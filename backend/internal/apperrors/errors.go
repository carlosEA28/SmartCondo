package apperrors

import "errors"

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrUserAlreadyExists     = errors.New("user already exists")
	ErrUserInUse             = errors.New("user has related records")
	ErrApartmentRequired     = errors.New("apartment is required for residents")
	ErrApartmentNotAllowed   = errors.New("apartment is only allowed for residents")
	ErrResponsibleNotAllowed = errors.New("only residents can be responsible for an apartment")
	ErrInvalidUserData       = errors.New("invalid user data")
)
