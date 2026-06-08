package helpers

import (
	"net/mail"
	"strings"
)

func IsEmailValid(email string) (bool) {
	if len(email) < 5 {
		return false
	}

	_, err := mail.ParseAddress(email)
	if err != nil{
		return false
	}
	return true
}

type PasswordError struct{
	Errors []string
}

func (e *PasswordError) Error() string{
	return strings.Join(e.Errors, "; ")
}

func IsPasswordValid(password string) (error){
	var errors []string
	if len(password) < 8{
		errors = append(errors, "password must consist of 8 charecters") 
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch :=  range password{
		if ch >= 'A' && ch <= 'Z'{
			hasUpper = true
		}else if ch >= 'a' && ch <= 'z'{
			hasLower = true
		}else if ch >= '0' && ch <= '9'{
			hasDigit = true
		}else if ch == '!' || ch == '_' || ch == '-' || ch == '@'{
			hasSpecial = true
		}else{
			errors = append(errors, "invalid character, only A-z, 0-9, !, _, -, @ available")
		}
	}

	if !hasUpper{
		errors = append(errors, "at least 1 upper case letter required")
	}
	if !hasLower{
		errors = append(errors, "at least 1 lower case letter required")
	}
	if !hasDigit{
		errors = append(errors, "at least 1 digit required")
	}
	if !hasSpecial{
		errors = append(errors, "at least 1 special symbol required")
	}

	if len(errors) > 0{
		return &PasswordError{Errors: errors}
	}
	
	return nil
}