package service

import (
	"fmt"
	"forum/models"
	"forum/pkg/repository"
	"time"
)

type PostService struct {
	repo    repository.Post
	repocat repository.Category
}

func NewPost(repo repository.Post, repocat repository.Category) *PostService {
	return &PostService{
		repo:    repo,
		repocat: repocat,
	}
}

func (s *PostService) CreatePost(p models.Post) error {
	cats, err := s.repocat.GetCategories()
	if err != nil {
		return err
	}
	if len(p.Category) < 1 {
		return ErrPostCategory
	}
	cat := containsAll(p.Category, cats)
	if cat != "" {
		return fmt.Errorf("invalid category %s", cat)
	}
	p.CreatedAt = time.Now()
	if err = isValidPostContent(p); err != nil {
		return err
	}

	return s.repo.CreatePost(p)
}

func (s *PostService) GetCategories() ([]string, error) {
	return s.repocat.GetCategories()
}

func (s *PostService) GetPost(id int) (models.Post, error) {
	post, err := s.repo.GetPost(id)
	if err != nil {
		return models.Post{}, err
	}
	post.Created = post.CreatedAt.Format("2006-01-02 15:04:05")
	if post.UpdatedAt != nil {
		uptime := post.UpdatedAt.Format("2006-01-02 15:04:05")
		post.Updated = &uptime
	}
	return post, nil
}

func (s *PostService) GetAllPosts(filter string) ([]models.Post, error) {
	return s.repo.GetAllPosts(filter)
}

func (s *PostService) GetFilteredByUserPosts(user_id int, filter string) ([]models.Post, error) {
	return s.repo.GetFilteredByUserPosts(user_id, filter)
}

func (s *PostService) UpdatePost(user_id int, updatedPost models.Post) error {
	post, err := s.repo.GetPost(updatedPost.ID)
	if err != nil {
		return err
	}

	if post.User_ID != user_id {
		return fmt.Errorf("%v, %v, updating post with id %v", ErrPermission, user_id, post.ID)
	}

	post.Title = updatedPost.Title
	post.Content = updatedPost.Content
	now := time.Now()
	post.UpdatedAt = &now
	if err = isValidPostContent(post); err != nil {
		return err
	}
	return s.repo.UpdatePost(post)
}

func (s *PostService) DeletePost(user_id, post_id int) error {
	post, err := s.GetPost(post_id)
	if err != nil {
		return err
	}
	if user_id != post.User_ID {
		return fmt.Errorf("%v, %v, deleting post with id %v", ErrPermission, user_id, post.ID)
	}
	return s.repo.DeletePost(user_id, post_id)
}

func (s *PostService) RatePost(rate models.RatePost) error {
	return s.repo.LikeDis(rate)
}

func containsAll(list, target []string) string {
	for i := range list {
		contains := false
		for j := range target {
			if list[i] == target[j] {
				contains = true
				break
			}
		}
		if !contains {
			return list[i]
		}
	}
	return ""
}
