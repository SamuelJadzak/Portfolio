package routes

import (
	"context"

	"github.com/lib/pq"
)

func getAllPosts(ctx context.Context, s *PostStore) ([]PostMeta, error) {
	posts := []PostMeta{}
	rows, err := s.Db.QueryContext(ctx, "select * from "+s.Location)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		var _throwaway []string
		if err := rows.Scan(&post.ID, &post.Title, pq.Array(&_throwaway)); err != nil {
			return nil, err
		}
		posts = append(posts, PostMeta{ID: post.ID, Title: post.Title})
	}
	return posts, nil
}
