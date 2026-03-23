package routes

import (
	"database/sql"
	"encoding/json"
	"example/data-access/auth"
	"example/data-access/helpers"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PostStore struct {
	Db       *sql.DB
	Location string
}

type PostMeta struct {
	ID    int
	Title string
}

type Post struct {
	PostMeta
	Body []string
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
			r0 := helpers.Try(getAllPosts(r.Context(), postStore))
			r1 := helpers.Bind(r0, func(posts []PostMeta) helpers.Result[[]byte] { return helpers.Try(helpers.JsonMarshal(posts)) })
			matchJson(w, r1)
		default:
			http.NotFound(w, r)
		}
	}

	singlePostHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			r0 := helpers.Try(strconv.Atoi(r.PathValue("id")))
			r1 := helpers.Bind(r0, func(id int) helpers.Result[Post] { return helpers.Try(getSinglePost(postStore, id)) })
			r2 := helpers.Bind(r1, func(post Post) helpers.Result[[]byte] { return helpers.Try(helpers.JsonMarshal(post)) })
			matchJson(w, r2)
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
			r0 := helpers.TryVoid(json.NewDecoder(r.Body).Decode(&creds))
			r1 := helpers.Then(r0, func() helpers.Result[struct{}] {
				return helpers.TryVoid(auth.CheckCredentials(authStore, creds.Username, creds.Password))
			})
			r2 := helpers.Then(r1, func() helpers.Result[string] { return helpers.Try(auth.CreateToken(creds.Username, "refresh")) })
			r3 := helpers.Bind(r2, func(refreshToken string) helpers.Result[[]byte] {
				return helpers.Bind(helpers.Try(auth.CreateToken(creds.Username, "access")), func(accessToken string) helpers.Result[[]byte] {
					return helpers.Try(helpers.JsonMarshal(map[string]string{"access_token": accessToken}))
				})
			})
			helpers.Match(r3, func(body []byte) {
				r2Success := r2.(helpers.Success[string])
				setCookie(w, r2Success.Data)
				if err := helpers.WriteJSON(w, append(body, '\n')); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}, func(err helpers.Error) { http.Error(w, err.Error(), helpers.ErrorCodeToHttp(err.Code)) })

		default:
			http.NotFound(w, r)
		}
	}

	refreshHandler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			r0 := helpers.Try(r.Cookie("token"))
			r1 := helpers.Bind(r0, func(cookie *http.Cookie) helpers.Result[*auth.CustomClaims] {
				return helpers.Try(auth.ValidateToken(cookie.Value, "refresh"))
			})
			r2 := helpers.Bind(r1, func(claims *auth.CustomClaims) helpers.Result[string] {
				return helpers.Try(auth.CreateToken(claims.Username, "access"))
			})
			r3 := helpers.Bind(r2, func(token string) helpers.Result[[]byte] {
				return helpers.Try(helpers.JsonMarshal(map[string]string{"access_token": token}))
			})
			matchJson(w, r3)
		default:
			http.NotFound(w, r)
		}
	}

	mux.HandleFunc("/api/posts/", postsHandler)
	mux.HandleFunc("/api/posts/{id}", singlePostHandler)
	mux.HandleFunc("/api/login", loginHandler)
	mux.HandleFunc("/api/refresh", refreshHandler)
	mux.Handle("/", http.NotFoundHandler())
}

func setCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
	})
}

func matchJson(w http.ResponseWriter, r helpers.Result[[]byte]) {
	helpers.Match(r, func(body []byte) {
		if err := helpers.WriteJSON(w, append(body, '\n')); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}, func(err helpers.Error) {
		http.Error(w, err.Error(), helpers.ErrorCodeToHttp(err.Code))
	})
}
