package domain

import "errors"

var (
    ErrUserExists    = errors.New("user already exists")
    ErrNotFound      = errors.New("not found")
    ErrInvalidRole   = errors.New("invalid user role")
    ErrInvalidAuth   = errors.New("invalid authentication credentials")
    ErrUserNotFound  = errors.New("user not found")
    ErrInvalidID     = errors.New("invalid uuid format")
)