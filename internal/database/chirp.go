package database

import (
	"errors"
	"sort"
)

type Chirp struct {
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
	Body     string `json:"body"`
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(author_id int, body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	uniqueId := len(dbStructure.Chirps) + 1

	chirp := Chirp{
		Id:       uniqueId,
		AuthorId: author_id,
		Body:     body,
	}

	dbStructure.Chirps[chirp.Id] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
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

// DeleteChirpById deletes chirp with matching id in the database
func (db *DB) DeleteChirpById(i int) (Chirp, error) {
	chirp := Chirp{}
	dbStructure, err := db.loadDB()
	if err != nil {
		return chirp, err
	}

	chirp, ok := dbStructure.Chirps[i]
	if !ok {
		return chirp, errors.New("Chirp not found")
	}
	delete(dbStructure.Chirps, i)

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}
