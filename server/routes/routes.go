package routes

import (
	"database/sql"
	"net/http"
	"strconv"
)

type PostStore struct {
	Db *sql.DB
}

type post struct {
	ID    int
	Title string
	Body  []string
}

func AddRoutes(
	mux *http.ServeMux,
	postStore *PostStore,
) {
	postsHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			posts, err := getAllPosts(postStore)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			for _, post := range posts {
				_, _ = w.Write([]byte(post.Title + "\n"))
			}
		default:
			http.NotFound(w, r)
		}
	}

	singlePostHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				http.NotFound(w, r)
				return
			}
			post, err := getSinglePost(postStore, id)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			_, _ = w.Write([]byte(post.Title + "\n"))
			for _, line := range post.Body {
				_, _ = w.Write([]byte(line + "\n"))
			}
		default:
			http.NotFound(w, r)
		}
	}

	mux.HandleFunc("/api/posts/", postsHandler)
	mux.HandleFunc("/api/posts/{id}", singlePostHandler)
	mux.Handle("/", http.NotFoundHandler())
}
