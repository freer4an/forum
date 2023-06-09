package service

import (
	"fmt"
	"forum/models"
	"forum/pkg/repository"
	"time"
)

type CommentService struct {
	repo repository.Comment
}

func NewComment(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) AddComment(com models.Comment) error {
	com.Created = time.Now().Format("02-01-2006 15:04")
	if err := isValidCommentContent(com); err != nil {
		return err
	}
	return s.repo.AddComment(com)
}

func (s *CommentService) GetComment(id int) (models.Comment, error) {
	return s.repo.GetComment(id)
}

func (s *CommentService) GetPostComments(id int) ([]models.Comment, error) {
	return s.repo.GetPostComments(id)
}

func (s *CommentService) UpdateComment(user_id int, updatedCom models.Comment) error {
	com, err := s.repo.GetComment(updatedCom.ID)
	if err != nil {
		return err
	}

	if com.User_ID != user_id {
		return fmt.Errorf("%v, %v, updating comment with id %v", ErrPermission, user_id, com.ID)
	}

	com.Content = updatedCom.Content
	current_time := time.Now().Format("02-01-2006 15:04")
	com.Updated = &current_time
	if err := isValidCommentContent(com); err != nil {
		return err
	}
	return s.repo.UpdateComment(com)
}

func (s *CommentService) DeleteComment(user_id, comment_id int) error {
	com, err := s.repo.GetComment(comment_id)
	if err != nil {
		return err
	}
	if user_id != com.User_ID {
		return fmt.Errorf("%v, %v, deleting comment with id %v", ErrPermission, user_id, com.ID)

	}
	return s.repo.DeleteComment(com.ID)
}

func (s *CommentService) RateComment(rate models.RateComment) error {
	return s.repo.RateComment(rate)
}
