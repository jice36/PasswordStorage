package userInterface

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jice36/PasswordStorage/internal/internalErrors"
	"github.com/jice36/PasswordStorage/internal/service"
)

type UI struct {
	s *service.Service
}

func NewUI(s *service.Service) *UI {
	return &UI{s: s}
}

// CmdLine serviceName, serviceLogin, password
func (ui *UI) CmdLine(ctxP context.Context) {
	ctx, cancel := context.WithCancel(ctxP)
	defer cancel()
	scan := bufio.NewScanner(os.Stdin)
loop:
	for {
		fmt.Println("Введите команду ")
		scan.Scan()
		command := scan.Text()

		args := parseArg(command)
		if len(args) == 0 && args[0] != "" {
			fmt.Println(internalErrors.NewErrorIncorrectNumberArgs().Error())
			continue
		}
		switch args[0] {
		case "add":
			if err := ui.add(args); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case "get":
			if err := ui.get(args); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case "delete":
			if err := ui.delete(args); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case "update":
			if err := ui.update(args); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case "all":
			if err := ui.all(args); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case "help":
			if err := ui.help(args); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case "quit":
			ui.quit(ctx, args)
			break loop
		default:
			fmt.Println("incorrect command")
		}
	}
}

func (ui *UI) add(args []string) error {
	fmt.Println(args)
	if len(args) != 4 {
		return internalErrors.NewErrorIncorrectNumberArgs()
	}
	err := ui.s.PutPassword(args[1], args[2],args[3])
	if err != nil {
		return err
	}
	return nil
}

func (ui *UI) get(args []string) error {
	if len(args) != 2 {
		return internalErrors.NewErrorIncorrectNumberArgs()
	}
	uai, err := ui.s.GetPassword(args[1])
	if err != nil {
		return err
	}
	fmt.Printf("login: %s password: %s date of creation: %s\n", uai.Login,uai.Password, uai.Date)
	return nil
}

func (ui *UI) delete(args []string) error {
	if len(args) != 2 {
		return internalErrors.NewErrorIncorrectNumberArgs()
	}
	err := ui.s.DeletePassword(args[1])
	if err != nil {
		return err
	}
	return nil
}

func (ui *UI) update(args []string) error {
	if len(args) != 4 {
		return internalErrors.NewErrorIncorrectNumberArgs()
	}
	err := ui.s.PutPassword(args[1], args[2],args[3])
	if err != nil {
		return err
	}
	return nil
}

func (ui *UI) all(args []string) error {
	if len(args) != 1 {
		return internalErrors.NewErrorIncorrectNumberArgs()
	}
	ap, err := ui.s.GetAllPasswords()
	if err != nil {
		return err
	}
	ap.Print()
	return nil
}

func (ui *UI) help(args []string) error {
	if len(args) != 1 {
		return internalErrors.NewErrorIncorrectNumberArgs()
	}

	h := "add service_name service_login password - добавить пароль \n" +
		"get service_name  - посмотреть пароль от сервиса\n" +
		"delete service_name - удалить пароль от сервиса\n" +
		"update service_name service_login password - обновить пароль от сервиса\n"+
		"quit - закрыть программу \n" +
		"all - вывести все пароли\n"
	fmt.Print(h)
	return nil
}

func (ui *UI) quit(ctxP context.Context, args []string) error {
	ctx, cancel := context.WithCancel(ctxP)
	if len(args) != 1 {
		return internalErrors.NewErrorIncorrectNumberArgs()
	}

	go func(ctx context.Context) {
		cancel()
	}(ctx)

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			fmt.Println("context timeout exceeded")
		case context.Canceled:
			fmt.Println("program quit")
		}
	}
	return nil
}

func parseArg(arg string) []string {
	return strings.Split(arg, " ")
}
