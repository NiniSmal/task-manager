package entity

import "errors"

var ErrNotFound = errors.New("not found")
var ErrNotVerification = errors.New("not verification")
var ErrNotAuthenticated = errors.New("not authenticated")
var ErrEmailExists = errors.New("email exists")
var ErrIncorrectEmail = errors.New("incorrect email")
var ErrIncorrectName = errors.New("name is too long or too short")
var ErrForbidden = errors.New("access denied")
