package apperrors

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrUserInUse              = errors.New("user has related records")
	ErrApartmentRequired      = errors.New("apartment is required for residents")
	ErrApartmentNotFound      = errors.New("apartment not found")
	ErrApartmentAlreadyExists = errors.New("apartment already registered")
	ErrInvalidUserData        = errors.New("invalid user data")
	ErrVisitorNotFound        = errors.New("visitor not found")
	ErrVisitorAlreadyExists   = errors.New("visitor already exists")
	ErrInvalidVisitorData     = errors.New("invalid visitor data")
)
