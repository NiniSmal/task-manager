package entity

import "errors"

var ErrNotFound = errors.New("not found")
var ErrNotVerification = errors.New("not verification")
var ErrNotAuthenticated = errors.New("not authenticated")
var ErrEmailExists = errors.New("email exists")
var ErrValidate = errors.New("validate")
var ErrForbidden = errors.New("access denied")
