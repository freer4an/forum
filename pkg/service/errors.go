package service

import "errors"

var (
	// authorization errors
	ErrInvalidName     = errors.New("invalid username")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidPassword = errors.New("invalid password")
	ErrAscii           = errors.New("non-ascii character")
	ErrUserNotFound    = errors.New("user not found")
	ErrWrongPassword   = errors.New("wrong password")
	ErrUserExists      = errors.New("username or password already exists")
	// post action errors
	ErrPostContent     = errors.New("the post content must be not empty")
	ErrPostCategory    = errors.New("the post must contain at least 1 category")
	ErrPostTitleUniq   = errors.New("the post title already exists")
	ErrPostTitleSyntax = errors.New("the post title must be not empty")
	// comment action errors
	ErrPermission     = errors.New("permission denied for user with id")
	ErrCommentContent = errors.New("the comment size must be not empty")
)
