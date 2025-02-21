package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// Database holds a collection of activities
type Database struct {
	filename string
	tokens   []UserToken
}

func NewDatabase() *Database {
	db := &Database{filename: "tokens.json"}
	db.Load()
	return db
}

// Load reads the database from the JSON file
func (db *Database) Load() error {
	file, err := os.Open(db.filename)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create an empty database
			db.tokens = []UserToken{}
			return nil
		}
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &db.tokens)
}

func (db *Database) Save() error {
	data, err := json.MarshalIndent(db.tokens, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(db.filename, data, 0644)
}

func (db *Database) AddToken(token UserToken) {
	if db.FindTokenById(token.AthleteId) != nil {
		db.UpdateToken(token)
	} else {
		db.tokens = append(db.tokens, token)
		db.Save()
	}
}

func (db *Database) ListTokens() {
	fmt.Println("Tokens:")
	for _, token := range db.tokens {
		fmt.Printf("%v\n", token)
	}
}

func (db *Database) FindTokenById(athleteId int64) *UserToken {
	for i, token := range db.tokens {
		if token.AthleteId == athleteId {
			return &db.tokens[i]
		}
	}
	return nil
}

func (db *Database) UpdateToken(token UserToken) {
	for i, t := range db.tokens {
		if t.AthleteId == token.AthleteId {
			db.tokens[i] = token
			db.Save()
			fmt.Println("Token updated.")
			return
		}
	}
	fmt.Println("Token not found.")
}

func (db *Database) DeleteToken(athleteId int64) {
	for i, token := range db.tokens {
		if token.AthleteId == athleteId {
			db.tokens = append(db.tokens[:i], db.tokens[i+1:]...)
			db.Save()
			fmt.Println("Token deleted.")
			return
		}
	}
	fmt.Println("Token not found.")
}
