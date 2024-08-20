package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Hien-Trinh/chirpy/internal/database"
	"github.com/golang-jwt/jwt/v5"
)

// GetUserByJWT returns a user by JWT
func GetUserByJWT(db *database.DB, jwtSecret, token string) (database.User, error) {
	token_parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return database.User{}, fmt.Errorf("couldn't parse token: %s", err)
	}

	claims, ok := token_parsed.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return database.User{}, fmt.Errorf("couldn't parse claims")
	}

	token_expiration_time, err := claims.GetExpirationTime()
	if err != nil {
		return database.User{}, fmt.Errorf("couldn't get expiration time: %s", err)
	}

	if token_expiration_time.Before(time.Now().UTC()) {
		return database.User{}, fmt.Errorf("token has expired")
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return database.User{}, fmt.Errorf("couldn't get subject: %s", err)
	}

	user_id, err := strconv.Atoi(subject)
	if err != nil {
		return database.User{}, fmt.Errorf("couldn't parse user ID")
	}

	user, err := db.GetUserById(user_id)
	if err != nil {
		return database.User{}, fmt.Errorf("couldn't get user: %s", err)
	}

	return user, nil
}
