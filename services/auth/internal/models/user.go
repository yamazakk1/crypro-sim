package models

type Role string

const (
	UserRole  Role = "user"
	AdminRole Role = "admin"
)

type User struct {
	UserID   string
	Email    string
	Username string
	Passwordhash string
	Balance_usdt float64
	Role Role
}

type GetMeUser struct{
	Username string
	Email string 
	Balance_usdt float64
	Role Role
}