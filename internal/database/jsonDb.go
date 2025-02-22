package database

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// jsonDatabase holds a collection of activities
type jsonDatabase struct {
	filename string
	tokens   []UserToken
}

func newJsonDatabase() *jsonDatabase {
	db := &jsonDatabase{filename: "tokens.json"}
	db.Load()
	return db
}

// Load reads the database from the JSON file
func (db *jsonDatabase) Load() error {
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

func (db *jsonDatabase) Save() error {
	data, err := json.MarshalIndent(db.tokens, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(db.filename, data, 0644)
}

func (db *jsonDatabase) AddToken(token UserToken) {
	if db.FindTokenById(token.AthleteId) != nil {
		db.UpdateToken(token)
	} else {
		db.tokens = append(db.tokens, token)
		db.Save()
	}
}

func (db *jsonDatabase) ListTokens() {
	fmt.Println("Tokens:")
	for _, token := range db.tokens {
		fmt.Printf("%v\n", token)
	}
}

func (db *jsonDatabase) FindTokenById(athleteId int64) *UserToken {
	for i, token := range db.tokens {
		if token.AthleteId == athleteId {
			return &db.tokens[i]
		}
	}
	return nil
}

func (db *jsonDatabase) UpdateToken(token UserToken) {
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

func (db *jsonDatabase) DeleteToken(athleteId int64) {
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
