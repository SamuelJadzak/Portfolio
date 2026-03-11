package routes

import (
	"database/sql"
	"encoding/json"
	"example/data-access/auth"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

type UserClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
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
			refreshToken, err := auth.CreateToken(creds.Username, "refresh")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			accessToken, err := auth.CreateToken(creds.Username, "access")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "token",
				Value:    refreshToken,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
				Expires:  time.Now().Add(time.Hour * 24 * 7),
			})
			w.Header().Set("Authorization", "Bearer "+accessToken)
		default:
			http.NotFound(w, r)
		}
	}

	refreshHandler := func(w http.ResponseWriter, r *http.Request) {
		refreshToken, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		_, claims, err := auth.ValidateToken(refreshToken.Value, "refresh")
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		accessToken, err := auth.CreateToken(claims.Username, "access")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Authorization", "Bearer "+accessToken)
	}

	mux.HandleFunc("/api/posts/", postsHandler)
	mux.HandleFunc("/api/posts/{id}", singlePostHandler)
	mux.HandleFunc("/api/login", loginHandler)
	mux.HandleFunc("/api/refresh", refreshHandler)
	mux.Handle("/", http.NotFoundHandler())
}
