package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// User
var (
	// ErrPasswordIsNotCorrect ...
	ErrPasswordIsNotCorrect = status.Errorf(codes.InvalidArgument, "password is not correct")

	// ErrEmailNotFound ...
	ErrEmailNotFound = status.Errorf(codes.NotFound, "email not found")

	// ErrUserIDNotFound ...
	ErrUserIDNotFound = status.Errorf(codes.NotFound, "user id not found")

	// ErrEmailInvalid ...
	ErrEmailInvalid = status.Errorf(codes.InvalidArgument, "email invalid")

	// ErrEmailAlreadyExists ...
	ErrEmailAlreadyExists = status.Errorf(codes.AlreadyExists, "email already exists")

	// ErrUserDontHavePermission ...
	ErrUserDontHavePermission = status.Errorf(codes.PermissionDenied, "user dont have permission")
)

// Session
var (
	// ErrSessionNotFound ...
	ErrSessionNotFound = status.Errorf(codes.NotFound, "session not found")
)

// Token
var (
	// ErrorTokenIsNotValid ...
	ErrorTokenIsNotValid = status.Errorf(codes.InvalidArgument, "token is not valid")
	// ErrorTokenNotFound ...
	ErrorTokenNotFound = status.Errorf(codes.InvalidArgument, "token is not found")
	// ErrorGenerateTokens ...
	ErrorGenerateTokens = status.Errorf(codes.Internal, "cannot generate tokens")
)
