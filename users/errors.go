package users

import "errors"

var (
	errEmptyID   = errors.New("User ID must be present")
	errEmptyName = errors.New("User name must be present")
	errNoUsers   = errors.New("No categories present")
	errNoUserId  = errors.New("User is not present")
)
