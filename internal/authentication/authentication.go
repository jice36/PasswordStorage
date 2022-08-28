package authentication

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/jice36/PasswordStorage/internal/internalErrors"
	"golang.org/x/crypto/pbkdf2"
	"sync"
)

/* table auth user hash */
func InitUser(login, userPassword string, dbname string) ([]byte, error) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	hashDB := ""
	hashUser := ""
	var key []byte

	go func(hashUser *string) {
		*hashUser, key = hashPasswordUser(login, userPassword)
		wg.Done()
	}(&hashUser)


	go func(hashDB *string, err  *error) {
		*hashDB, *err = getHashDB(db, login)
		wg.Done()
	}(&hashDB, &err)

	wg.Wait()

	if err != nil {
		return nil, err
	}
	if hashDB != hashUser {
		return nil, internalErrors.NewErrorIncorrectPassword()
	}
	return key, nil
}

func getHashDB(db *sql.DB, login string) (string, error) {
	var hashRow sql.NullString

	row := db.QueryRow("select hash from auth where login = $1", login)
	err := row.Scan(&hashRow)
	if err != nil {
		return "", err
	}

	if !hashRow.Valid {
		return "", errors.New("hash is null")
	}
	return hashRow.String, nil
}

func hashPasswordUser(login, userPassword string) (string, []byte) {
	key := pbkdf2.Key([]byte(userPassword), []byte(login), 10, 32, sha256.New)

	sha := sha256.New()
	sha.Write(key)
	hash := sha.Sum(nil)
	return hex.EncodeToString(hash), key
}
