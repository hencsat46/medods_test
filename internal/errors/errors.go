package customErrors

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrTokenExpired = errors.New("token is expired")
	ErrTokensMatch  = errors.New("tokens doesn't match")
)
