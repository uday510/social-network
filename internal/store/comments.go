package store

import (
	"context"
	"database/sql"
	"log"
)

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

func (s *CommentsStore) GetByPostID(ctx context.Context, postID int64) ([]Comment, error) {

	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username, u.id FROM comments c
		JOIN users u on u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
		`

	rows, err := s.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var comments []Comment
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
		`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

type CommentsStore struct {
	db *sql.DB
}
