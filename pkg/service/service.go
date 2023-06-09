package service

import (
	"forum/models"
	"forum/pkg/repository"
)

type Authorization interface {
	CreateUser(user models.User) (int, error)
	GetUser(name, password string) (models.User, error)
	GetUserByID(id int) (models.User, error)
}

type Session interface {
	CreateSession(models.Session) error
	GetSession(token string) (models.Session, error)
	DeleteSession(user_id int) error
}

type Post interface {
	CreatePost(p models.Post) error
	GetPost(id int) (models.Post, error)
	GetAllPosts(filter string) ([]models.Post, error)
	GetFilteredByUserPosts(user_id int, filter string) ([]models.Post, error)
	UpdatePost(user_id int, p models.Post) error
	DeletePost(user_id, post_id int) error
	GetCategories() ([]string, error)
	RatePost(models.RatePost) error
}

type Comment interface {
	AddComment(models.Comment) error
	GetComment(id int) (models.Comment, error)
	GetPostComments(id int) ([]models.Comment, error)
	UpdateComment(user_id int, com models.Comment) error
	DeleteComment(user_id, comment_id int) error
	RateComment(models.RateComment) error
}

type Service struct {
	Authorization
	Session
	Post
	Comment
}

func New(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuth(repo.Authorization),
		Session:       NewSession(repo.Session),
		Post:          NewPost(repo.Post, repo.Category),
		Comment:       NewComment(repo.Comment),
	}
}
