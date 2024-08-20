package database

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"sort"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}
type DBStructure struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RefreshTokens map[int]RefreshToken `json:"refresh_tokens"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshToken struct {
	Id        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {
	db := &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}

	err := db.ensureDB()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	uniqueId := len(dbStructure.Chirps) + 1

	chirp := Chirp{
		Id:   uniqueId,
		Body: body,
	}

	dbStructure.Chirps[chirp.Id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// CreateUser creates a new user and saves it to disk
func (db *DB) CreateUser(email, password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	uniqueId := len(dbStructure.Users) + 1

	user := User{
		Id:       uniqueId,
		Email:    email,
		Password: password,
	}

	dbStructure.Users[user.Id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
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

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	sort.Slice(chirps, func(i, j int) bool { return chirps[i].Id < chirps[j].Id })

	return chirps, nil
}

// GetUsers returns all users in the database
func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool { return users[i].Id < users[j].Id })

	return users, nil
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

// GetChirpsById returns chirp with matching id in the database
func (db *DB) GetChirpById(i int) (Chirp, error) {
	chirp := Chirp{}
	dbStructure, err := db.loadDB()
	if err != nil {
		return chirp, err
	}

	chirp, ok := dbStructure.Chirps[i]
	if !ok {
		return chirp, errors.New("Chirp not found")
	}

	return chirp, nil
}

// GetUserById returns user with matching id in the database
func (db *DB) GetUserById(i int) (User, error) {
	user := User{}
	dbStructure, err := db.loadDB()
	if err != nil {
		return user, err
	}

	user, ok := dbStructure.Users[i]
	if !ok {
		return user, errors.New("User not found")
	}

	return user, nil
}

func (db *DB) UpdateUser(i int, new_email, new_password string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	_, ok := dbStructure.Users[i]
	if !ok {
		return User{}, errors.New("User not found")
	}

	new_user := User{
		Id:       i,
		Email:    new_email,
		Password: new_password,
	}
	dbStructure.Users[i] = new_user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return new_user, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if errors.Is(err, os.ErrNotExist) || *dbg {
		dbStructure := DBStructure{
			Chirps:        make(map[int]Chirp),
			Users:         make(map[int]User),
			RefreshTokens: make(map[int]RefreshToken),
		}
		return db.writeDB(dbStructure)
	}

	return err
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	dbStructure := DBStructure{}

	file, err := os.ReadFile(db.path)
	if err != nil {
		return dbStructure, err
	}

	err = json.Unmarshal(file, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	file, err := json.MarshalIndent(dbStructure, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("database.json", file, 0644)

	return err
}
