package entity

import "errors"

var ErrNotFound = errors.New("not found")
var ErrNotVerification = errors.New("not verification")
var ErrNotAuthenticated = errors.New("not authenticated")
var ErrLoginExists = errors.New("login exists")
