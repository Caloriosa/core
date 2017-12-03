package validation

import (
	"core/types"
	"errors"
)

func ValidateNewUser (user *types.User) error {
	if user.Email == "" {
		return errors.New("Missing e-mail")
	}
	if user.Password == "" {
		return errors.New("Missing password")
	}
	if user.Login == "" {
		return errors.New("Missing login")
	}
	return ValidateEmail(user.Email)
}