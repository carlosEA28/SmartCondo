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
	ErrVisitNotFound          = errors.New("visit not found")
	ErrPorteiroNotFound       = errors.New("porteiro not found")
	ErrInvalidPorteiroData    = errors.New("invalid porteiro data")
	ErrFilterRequired          = errors.New("at least one search filter is required")
	ErrComunicadoNotFound      = errors.New("comunicado not found")
	ErrInvalidComunicadoData   = errors.New("invalid comunicado data")
	ErrMissingAuthHeader       = errors.New("missing authentication header")
	ErrUnauthorizedSindico     = errors.New("user is not authorized as sindico")
	ErrComunicadoNotOwner      = errors.New("you can only delete your own comunicados")
)
