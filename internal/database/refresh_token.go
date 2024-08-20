package database

import (
	"errors"
	"sort"
	"time"
)

type RefreshToken struct {
	Id        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// CreateRefreshToken creates a new refresh token and saves it to disk
func (db *DB) CreateRefreshToken(user_id int, refresh_token_string string, refresh_token_expires_at time.Time) (RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, err
	}

	uniqueId := len(dbStructure.RefreshTokens) + 1

	refresh_token := RefreshToken{
		Id:        uniqueId,
		UserID:    user_id,
		Token:     refresh_token_string,
		ExpiresAt: refresh_token_expires_at,
	}

	dbStructure.RefreshTokens[refresh_token.Id] = refresh_token

	err = db.writeDB(dbStructure)
	if err != nil {
		return RefreshToken{}, err
	}

	return refresh_token, nil
}

// GetRefreshTokens returns all refresh tokens in the database
func (db *DB) GetRefreshTokens() ([]RefreshToken, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	refresh_tokens := make([]RefreshToken, 0, len(dbStructure.RefreshTokens))
	for _, refresh_token := range dbStructure.RefreshTokens {
		refresh_tokens = append(refresh_tokens, refresh_token)
	}

	sort.Slice(refresh_tokens, func(i, j int) bool { return refresh_tokens[i].Id < refresh_tokens[j].Id })

	return refresh_tokens, nil
}

func (db *DB) RevokeRefreshToken(i int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	_, ok := dbStructure.RefreshTokens[i]
	if !ok {
		return errors.New("refresh token not found")
	}

	delete(dbStructure.RefreshTokens, i)

	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	return nil
}
