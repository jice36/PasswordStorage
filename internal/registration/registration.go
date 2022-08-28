package registration

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
)

func RegistrationUser(login, password string, dbname string) ([]byte, error) {
	key, err := addUser(login, password, dbname)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func addUser(login, password string, dbname string) ([]byte, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	key, hashPassword := getHashPassword(login, password)

	_, err = db.Exec("insert into auth(login, hash) values($1,$2)", login, hashPassword)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func getHashPassword(login, password string) ([]byte, string) {
	key := pbkdf2.Key([]byte(password), []byte(login), 10, 32, sha256.New)

	sha := sha256.New()
	sha.Write(key)
	hash := sha.Sum(nil)
	return key, hex.EncodeToString(hash)
}
