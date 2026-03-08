package auth

import (
	"database/sql"
	"example/data-access/env"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthStore struct {
	Db       *sql.DB
	Location string
}

func CreateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	signedToken, err := token.SignedString([]byte(env.JwtSecret.GetValue()))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	validMethods := []string{"HS256"}
	parser := jwt.NewParser(jwt.WithValidMethods(validMethods))
	return parser.ParseWithClaims(tokenString, &jwt.MapClaims{}, myKeyFunc)
}

func myKeyFunc(token *jwt.Token) (any, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
	}
	return []byte(env.JwtSecret.GetValue()), nil
}

func comparePasswords(hashedPwd []byte, plainPwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPwd, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func CheckCredentials(authStore *AuthStore, username string, password string) bool {
	var hashedPwd string
	err := authStore.Db.QueryRow("SELECT password FROM "+authStore.Location+" WHERE username = $1", username).Scan(&hashedPwd)
	if err != nil {
		return false
	}
	return comparePasswords([]byte(hashedPwd), []byte(password))
}

func InitAdminPwd(authStore *AuthStore) {
	pwd := env.ServerAdminPassword.GetValue()
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	_, err = authStore.Db.Exec("INSERT INTO "+authStore.Location+" (username, password) VALUES ($1, $2)", env.ServerAdminUsername.GetValue(), string(hashedPwd))
	if err != nil {
		log.Fatal(err)
	}
}
