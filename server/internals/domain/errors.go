package domain

import "errors"

var (
	ErrUserNotFound          = errors.New("User Not found")
	ErrNotFound              = errors.New("Not found")
	ErrConflict              = errors.New("This email is already in use")
	ErrTokeNotFoundOrExpired = errors.New("Token Not found or expired")
	ErrInvalidCredentials    = errors.New("Invalid Credentials")
	ErrInvalidRefreshToken   = errors.New("Invalid Refresh Token")
	ErrAlreadyVerified       = errors.New("Email already verified")
	ErrNotVerified           = errors.New("account not verified")
	ErrWeakPassword          = errors.New("Password must be more than 8 characters")
	ErrAccountDeactivated    = errors.New("account has been deactivated")
	ErrNameConflict          = errors.New("This name is alredy in use")
	ErrInvalidFileExtension  = errors.New("This is not a valid image extension %s")
)
