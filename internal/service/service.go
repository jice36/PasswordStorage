package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jice36/PasswordStorage/config"
	_ "github.com/mattn/go-sqlite3"

	ci "github.com/jice36/PasswordStorage/internal/cipher"
)

type Service struct {
	Db     *sql.DB
	cipher *ci.Cipher
	user   string
}

// NewService /* password_storage, nameService password */
func NewService(conf *config.Config, key []byte, login string) (*Service, error) {
	db, err := sql.Open("sqlite3", conf.Database.Dbname)
	if err != nil {
		return nil, err
	}
	cipher, err := ci.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &Service{
		Db:     db,
		cipher: cipher,
		user:   login,
	}, nil
}

func (s *Service) PutPassword(serviceName, serviceLogin, password string) error {
	encLogin := s.encryption(serviceLogin)
	encPass := s.encryption(password)
	q := `INSERT INTO storage_passwords (login, service_name, service_login,password,create_date)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (password) DO UPDATE
			SET password = $4 where login = $1;`
	_, err := s.Db.Exec(q, s.user, serviceName, encLogin, encPass, time.Now().Format(time.RFC850))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) encryption(data string) string {
	out := s.cipher.EncryptModeECB([]byte(data))
	return string(out)
}

type UserAccountInformation struct {
	Login    string
	Password string
	Date     string
}

func (s *Service) GetPassword(serviceName string) (*UserAccountInformation, error) {
	var encLogin, encPass, date sql.NullString
	row := s.Db.QueryRow("select service_login,password, create_date from storage_passwords where service_name = $1 AND login = $2", serviceName, s.user)

	err := row.Scan(&encLogin, &encPass, &date)
	if err != nil {
		return nil, err
	}

	if !encPass.Valid && !encLogin.Valid && !date.Valid{
		return nil, errors.New("failed to get info")
	}
	decLogin := s.decryption(encLogin.String)
	decPass := s.decryption(encPass.String)

	uai := &UserAccountInformation{
		Login:    decLogin,
		Password: decPass,
		Date:     date.String,
	}

	return uai, nil
}

func (s *Service) decryption(data string) string {
	out := s.cipher.DecryptModeECB([]byte(data))
	return string(out)
}

func (s *Service) GetAllPasswords() (*AllPasswords, error) {
	rows, err := s.Db.Query("select service_name, service_login, password, create_date  from storage_passwords where login = $1", s.user)
	if err != nil {
		return nil, err
	}

	names := make([]string, 0)
	logins := make([]string,0)
	pass := make([]string, 0)
	dates := make([]string,0)
	for rows.Next() {
		var serviceName, encLogin, encPass, date sql.NullString
		err = rows.Scan(&serviceName, &encLogin, &encPass, &date)
		if err != nil {
			return nil, err
		}

		if !(serviceName.Valid && encPass.Valid && encLogin.Valid && date.Valid) {
			return nil, errors.New("failed to get info")
		}

		names = append(names, serviceName.String)

		decLogin := s.decryption(encLogin.String)
		logins = append(logins, decLogin)

		decPass := s.decryption(encPass.String)
		pass = append(pass, decPass)

		dates = append(dates, date.String)
	}

	return &AllPasswords{
		services:  names,
		logins: logins,
		passwords: pass,
		dates: dates,
	}, nil
}

func (s *Service) DeletePassword(serviceName string) error {
	_, err := s.Db.Exec("delete from storage_passwords where service_name = $1 and where user = $2", serviceName, s.user)
	if err != nil {
		return err
	}
	return nil
}

type AllPasswords struct {
	services  []string
	logins []string
	passwords []string
	dates []string
}

func (ap *AllPasswords) Print() {
	length := len(ap.services)
	for i := 0; i < length-1; i++ {
		fmt.Printf("name: %s \n", ap.services[i])
		fmt.Printf("login: %s \n", ap.logins[i])
		fmt.Printf("password: %s \n", ap.passwords[i])
		fmt.Printf("date of creation: %s \n", ap.dates[i])
		fmt.Println()
	}
	fmt.Printf("name: %s \n", ap.services[length-1])
	fmt.Printf("login: %s \n", ap.logins[length-1])
	fmt.Printf("password: %s \n", ap.passwords[length-1])
	fmt.Printf("date of creation: %s \n", ap.dates[length-1])
}
