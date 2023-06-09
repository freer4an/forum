package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/models"
	"strings"
)

const (
	TopPosts       = "popular"
	UserLikedPosts = "my-liked-posts"
)

type PostSQL struct {
	db *sql.DB
}

func NewPostSQL(db *sql.DB) *PostSQL {
	return &PostSQL{
		db: db,
	}
}

func (r *PostSQL) CreatePost(p models.Post) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO post (user_id, author, title, content, created, updated) values (?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(p.User_ID, p.Author, p.Title, p.Content, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	postID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, categoryName := range p.Category {
		_, err = tx.Exec("INSERT INTO post_category(post_id, category_name) VALUES (?, ?)", postID, categoryName)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *PostSQL) GetPost(id int) (models.Post, error) {
	var p models.Post
	query := `
	SELECT post.*,
	COUNT(DISTINCT CASE WHEN pr.islike = 1 THEN pr.user_id || '-' || pr.post_id END) AS like_count,
	COUNT(DISTINCT CASE WHEN pr.islike = -1 THEN pr.user_id || '-' || pr.post_id END) AS dislike_count,
		GROUP_CONCAT(DISTINCT category.name) AS categories
	FROM post 
	LEFT JOIN post_rating AS pr ON post.id = pr.post_id
	LEFT JOIN post_category AS pc ON post.id = pc.post_id
	LEFT JOIN category ON pc.category_name = category.name
	WHERE
		post.id = ?
	GROUP BY post.id;
	`
	row := r.db.QueryRow(query, id)
	var categories sql.NullString
	if err := row.Scan(&p.ID, &p.User_ID, &p.Author, &p.Title, &p.Content, &p.CreatedAt, &p.UpdatedAt, &p.Likes, &p.Dislikes, &categories); err != nil {
		return models.Post{}, fmt.Errorf("error scanning Post Details: %v", err)
	}
	if categories.Valid {
		p.Category = strings.Split(categories.String, ",")
	} else {
		p.Category = []string{}
	}
	return p, nil
}

// Get all posts with their categories and number of likes, dislikes and comments.
func (r *PostSQL) GetAllPosts(filter string) ([]models.Post, error) {
	var query string
	if filter == TopPosts {
		query = `
		SELECT post.*,
    		COUNT(DISTINCT comment.id) AS comment_count,
		    COUNT(DISTINCT CASE WHEN pr.islike = 1 THEN pr.user_id || '-' || pr.post_id END) AS like_count,
		    COUNT(DISTINCT CASE WHEN pr.islike = -1 THEN pr.user_id || '-' || pr.post_id END) AS dislike_count,
		    GROUP_CONCAT(DISTINCT category.name) AS categories
		FROM post
		LEFT JOIN comment ON post.id = comment.post_id
		LEFT JOIN post_rating AS pr ON post.id = pr.post_id
		LEFT JOIN post_category AS pc ON post.id = pc.post_id
		LEFT JOIN category ON pc.category_name = category.name
		GROUP BY post.id
		ORDER by like_count DESC;
		`
	} else {
		query = `
		SELECT post.*,
		    COUNT(DISTINCT comment.id) AS comment_count,
		    COUNT(DISTINCT CASE WHEN pr.islike = 1 THEN pr.user_id || '-' || pr.post_id END) AS like_count,
		    COUNT(DISTINCT CASE WHEN pr.islike = -1 THEN pr.user_id || '-' || pr.post_id END) AS dislike_count,
		    GROUP_CONCAT(DISTINCT category.name) AS categories
		FROM post
		LEFT JOIN comment ON post.id = comment.post_id
		LEFT JOIN post_rating AS pr ON post.id = pr.post_id
		LEFT JOIN post_category AS pc ON post.id = pc.post_id
		LEFT JOIN category ON pc.category_name = category.name
		GROUP BY post.id
		ORDER BY post.created DESC;
		`
	}

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []models.Post{}
	for rows.Next() {
		var post models.Post
		var categories sql.NullString // Use sql.NullString instead of string
		err = rows.Scan(
			&post.ID,
			&post.User_ID,
			&post.Author,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Comments,
			&post.Likes,
			&post.Dislikes,
			&categories, // Scan as sql.NullString
		)
		if err != nil {
			return nil, err
		}
		post.Created = post.CreatedAt.Format("02-01-2006 15:04:05")
		if post.UpdatedAt != nil {
			uptime := post.UpdatedAt.Format("02-01-2006 15:04:05")
			post.Updated = &uptime
		}
		if categories.Valid {
			post.Category = strings.Split(categories.String, ",")
		} else {
			post.Category = []string{} // Set an empty slice for NULL values
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

func (r *PostSQL) UpdatePost(p models.Post) error {
	stmt, err := r.db.Prepare("UPDATE post SET title = ?, content = ?, updated = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(p.Title, p.Content, p.UpdatedAt, p.ID)
	return err
}

func (r *PostSQL) DeletePost(user_id, post_id int) error {
	stmt, err := r.db.Prepare("DELETE FROM post WHERE user_id = ? AND id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user_id, post_id)
	return err
}

func (r *PostSQL) LikeDis(rate models.RatePost) error {
	var oldIslike int8

	err := r.db.QueryRow("SELECT islike FROM post_rating WHERE user_id = ? AND post_id = ?", rate.User_ID, rate.Post_ID).Scan(&oldIslike)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	}
	if oldIslike == rate.IsLike {
		stmt, err := r.db.Prepare("DELETE FROM post_rating WHERE user_id = ? AND post_id = ?")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(rate.User_ID, rate.Post_ID)
		return err
	}
	stmt, err := r.db.Prepare(`INSERT INTO post_rating (user_id, post_id, islike) 
	VALUES (?, ?, ?) 
	ON CONFLICT(user_id, post_id) DO UPDATE 
	SET islike = excluded.islike`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(rate.User_ID, rate.Post_ID, rate.IsLike)
	return err
}

func (r *PostSQL) GetFilteredByUserPosts(user_id int, filter string) ([]models.Post, error) {
	var query string
	if filter == UserLikedPosts {
		query = `
		SELECT post.*,
			COUNT(DISTINCT comment.id) AS comment_count,
			SUM(CASE WHEN pr.islike = 1 THEN 1 ELSE 0 END) AS like_count,
       		SUM(CASE WHEN pr.islike = -1 THEN 1 ELSE 0 END) AS dislike_count,
			GROUP_CONCAT(DISTINCT category.name) AS categories
		FROM post
		LEFT JOIN comment ON post.id = comment.post_id
		LEFT JOIN post_rating AS pr ON post.id = pr.post_id
		LEFT JOIN post_category AS pc ON post.id = pc.post_id
		LEFT JOIN category ON pc.category_name = category.name
		WHERE pr.user_id = ? AND islike = 1
		GROUP BY post.id
		ORDER by created DESC;
		`
	} else {
		query = `
		SELECT post.*,
		COUNT(DISTINCT comment.id) AS comment_count,
			COUNT(DISTINCT CASE WHEN pr.islike = 1 THEN pr.user_id || '-' || pr.post_id END) AS like_count,
			COUNT(DISTINCT CASE WHEN pr.islike = -1 THEN pr.user_id || '-' || pr.post_id END) AS dislike_count,
			GROUP_CONCAT(DISTINCT category.name) AS categories
		FROM post
		LEFT JOIN comment ON post.id = comment.post_id
		LEFT JOIN post_rating AS pr ON post.id = pr.post_id
		LEFT JOIN post_category AS pc ON post.id = pc.post_id
		LEFT JOIN category ON pc.category_name = category.name
		WHERE post.user_id = ?
		GROUP BY post.id
		ORDER by created DESC;
		`
	}

	rows, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := []models.Post{}
	for rows.Next() {
		var post models.Post
		var categories sql.NullString // Use sql.NullString instead of string
		err = rows.Scan(
			&post.ID,
			&post.User_ID,
			&post.Author,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.Comments,
			&post.Likes,
			&post.Dislikes,
			&categories, // Scan as sql.NullString
		)
		if err != nil {
			return nil, err
		}
		post.Created = post.CreatedAt.Format("02-01-2006 15:04:05")
		if post.UpdatedAt != nil {
			uptime := post.UpdatedAt.Format("02-01-2006 15:04:05")
			post.Updated = &uptime
		}
		if categories.Valid {
			post.Category = strings.Split(categories.String, ",")
		} else {
			post.Category = []string{} // Set an empty slice for NULL values
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
