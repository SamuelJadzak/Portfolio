package routes

import (
	"github.com/lib/pq"
)

func getSinglePost(s *PostStore, id int) (*post, error) {
	var post post
	row := s.Db.QueryRow("SELECT * FROM "+s.Location+" WHERE id = $1", id)
	err := row.Scan(&post.ID, &post.Title, pq.Array(&post.Body))
	if err != nil {
		return nil, err
	}

	return &post, nil
}
