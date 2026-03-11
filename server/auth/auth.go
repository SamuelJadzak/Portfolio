package auth

import (
	"database/sql"
	"example/data-access/env"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthStore struct {
	Db       *sql.DB
	Location string
}

type CustomClaims struct {
	Username string `json:"username"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

func CreateToken(username string, tokenType string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		Username: username,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	secret, err := getSecret(tokenType)
	if err != nil || secret == "" {
		return "", fmt.Errorf("failed to get secret: %w", err)
	}
	key := []byte(secret)

	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateToken(tokenString string, tokenType string) (*jwt.Token, *CustomClaims, error) {
	cleanTokenString := strings.TrimPrefix(tokenString, "Bearer ")
	validMethods := []string{"HS256"}
	parser := jwt.NewParser(jwt.WithValidMethods(validMethods))
	claims := &CustomClaims{}
	token, err := parser.ParseWithClaims(cleanTokenString, claims, myKeyFunc(tokenType))
	if err != nil {
		return nil, nil, err
	}
	return token, claims, nil
}

func myKeyFunc(tokenType string) func(token *jwt.Token) (any, error) {
	secret, err := getSecret(tokenType)
	key := []byte(secret)

	return func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if err != nil || secret == "" {
			return nil, fmt.Errorf("invalid secret for %s: %w", tokenType, err)
		}
		return key, nil
	}
}

func getSecret(tokenType string) (string, error) {
	var secret string
	switch tokenType {
	case "access":
		secret = env.AccessTokenSecret.GetValue()
	case "refresh":
		secret = env.RefreshTokenSecret.GetValue()
	default:
		return "", fmt.Errorf("invalid token type: %s", tokenType)
	}
	return secret, nil
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
