package routes

import (
	"database/sql"
	"encoding/json"
	"example/data-access/auth"
	"net/http"
	"strconv"
	"time"
)

type PostStore struct {
	Db       *sql.DB
	Location string
}

type post struct {
	ID    int
	Title string
	Body  []string
}

func AddRoutes(
	mux *http.ServeMux,
	postStore *PostStore,
	authStore *auth.AuthStore,
) {
	postsHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			posts, err := getAllPosts(r.Context(), postStore)
			if err != nil {
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

	loginHandler := func(w http.ResponseWriter, r *http.Request) {
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		switch r.Method {
		case "POST":
			if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			if !auth.CheckCredentials(authStore, creds.Username, creds.Password) {
				http.Error(w, "Invalid credentials", http.StatusUnauthorized)
				return
			}
			token, err := auth.CreateToken(creds.Username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    token,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().Add(time.Hour),
			})

		default:
			http.NotFound(w, r)
		}
	}

	mux.HandleFunc("/api/posts/", postsHandler)
	mux.HandleFunc("/api/posts/{id}", singlePostHandler)
	mux.HandleFunc("/api/login", loginHandler)
	mux.Handle("/", http.NotFoundHandler())
}
