package routes

import (
	"context"

	"github.com/lib/pq"
)

func getAllPosts(ctx context.Context, s *PostStore) ([]post, error) {
	posts := []post{}
	rows, err := s.Db.QueryContext(ctx, "select * from "+s.Location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post post
		var _throwaway []string
		if err := rows.Scan(&post.ID, &post.Title, pq.Array(&_throwaway)); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
