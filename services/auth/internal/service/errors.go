package service

import "errors"

var(
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrWrongPassword = errors.New("wrong password")
)