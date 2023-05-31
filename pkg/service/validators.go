package service

import (
	"forum/models"
	"regexp"
	"strings"
	"unicode"
)

const (
	nameRegExp  = `^[a-zA-Z0-9_.'-]{3,15}$`
	emailRegExp = `^[\w-.]+@([\w-]+.)+[\w-]{2,4}$`
)

func validateNewUserData(user models.User) error {
	if !isAscii(user.Name) || !isAscii(user.Password) || !isAscii(user.Password) {
		return ErrAscii
	}
	if !isValidName(user.Name) {
		return ErrInvalidName
	}
	if !isValidEmail(user.Email) {
		return ErrInvalidEmail
	}
	return isValidPassword(user.Password)
}

func isAscii(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func isValidName(name string) bool {
	re := regexp.MustCompile(nameRegExp)
	return re.MatchString(name)
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(emailRegExp)
	return re.MatchString(email)
}

func isValidPassword(password string) error {
	var err error
	match1, _ := regexp.MatchString("[A-Z]", password)
	match2, _ := regexp.MatchString("[a-z]", password)
	match3, _ := regexp.MatchString("[~!@#$%^&*_()-+={[}]|\\:;\"'<,>.?/]", password)

	if !match1 || !match2 || !match3 {
		return ErrInvalidPassword
	}

	return err
}

func isValidPostContent(p models.Post) error {
	if len([]rune(strings.TrimSpace(p.Content))) < 1 {
		return ErrPostContent
	} else if len([]rune(p.Title)) < 1 {
		return ErrPostTitleSyntax
	}
	return nil
}

func isValidCommentContent(p models.Comment) error {
	if len([]rune(strings.TrimSpace(p.Content))) < 1 {
		return ErrCommentContent
	}
	return nil
}
