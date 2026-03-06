package routes

import (
	"log"

	"github.com/lib/pq"
)

func getAllPosts(s *PostStore) ([]post, error) {
	posts := []post{}
	rows, err := s.Db.Query("select * from posts")
	if err != nil || rows == nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var post post
		var _throwaway []string
		err := rows.Scan(&post.ID, &post.Title, pq.Array(&_throwaway))
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
