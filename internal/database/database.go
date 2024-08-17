package database

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"sort"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
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
func (db *DB) CreateUser(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	uniqueId := len(dbStructure.Users) + 1

	user := User{
		Id:    uniqueId,
		Email: email,
	}

	dbStructure.Users[user.Id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return User{}, err
	}

	return user, nil
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

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if errors.Is(err, os.ErrNotExist) || *dbg {
		dbStructure := DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
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
