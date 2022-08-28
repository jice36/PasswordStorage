package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/jice36/PasswordStorage/internal/initDB"
	"github.com/jice36/PasswordStorage/internal/internalErrors"
	"log"

	"github.com/jice36/PasswordStorage/config"
	"github.com/jice36/PasswordStorage/internal/authentication"
	"github.com/jice36/PasswordStorage/internal/registration"
	"github.com/jice36/PasswordStorage/internal/service"
	"github.com/jice36/PasswordStorage/internal/userInterface"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "config.yml", "path to config file")
}

func main() {
	flag.Parse()

	conf, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = initDB.CheckDB(conf.Database.Dbname)
	if err != nil {
		log.Fatal(err.Error())
	}

	var login, pass string
	var key []byte
	fmt.Println("Login ")
	fmt.Scanln(&login)
	fmt.Println("Password ")
	fmt.Scanln(&pass)
	key, err = authentication.InitUser(login, pass, conf.Database.Dbname)
	if err != nil {
		if errors.Is(err, internalErrors.NewErrorIncorrectPassword()) {
			log.Fatal("wrong password or login")
		} else if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("User is not found\n Register? Y/N")
			reg := ""
			fmt.Scanln(&reg)
			switch reg {
			case "Y":
				key, err = registration.RegistrationUser(login, pass, conf.Database.Dbname)
				if err != nil {
					log.Fatal(err.Error())
				}
			case "N":
				log.Fatal("quit program")
			default:
				log.Fatal("incorrect command")
			}
		} else{
			log.Fatal(internalErrors.NewErrorAuthentication(err))
		}
	}

	fmt.Println("Hello!")

	s, err := service.NewService(conf, key, login)
	if err != nil {
		log.Fatal(err.Error())
	}

	ctxP := context.Background()
	ctx, cancel := context.WithCancel(ctxP)
	defer func() {
		s.Db.Close()
		cancel()
	}()

	ui := userInterface.NewUI(s)
	ui.CmdLine(ctx)
}
