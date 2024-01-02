package database

import (
	"time"
)

type Revocation struct {
	Token     string    `json:"token"`
	RevokedAt time.Time `json:"revoked_at"`
}

func (db *DB) RevokeToken(token string) error {
	dbStucture, err := db.loadDB()
	if err != nil {
		return err
	}

	revocation := Revocation{
		Token:     token,
		RevokedAt: time.Now().UTC(),
	}
	dbStucture.Revocations[token] = revocation

	err = db.writeDB(dbStucture)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) IsTokenRevoked(token string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}

	revocation, ok := dbStructure.Revocations[token]
	if !ok {
		return false, nil
	}

	if revocation.RevokedAt.IsZero() {
		return false, nil
	}

	return true, nil
}
