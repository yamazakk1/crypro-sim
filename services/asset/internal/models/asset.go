package models

import "time"

type Asset struct {
	ID        string
	Symbol    string
	Fullname  string
	IsActive  bool
	InitialPrice float64
	CreatedAt time.Time
	UpdatedAt time.Time
}
