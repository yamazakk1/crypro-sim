package service

import "errors"

var ErrEmptyFields = errors.New("symbol and fullname are required")