package ui

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/vanamelnik/gophkeeper/client"
	"github.com/vanamelnik/gophkeeper/client/repo"
	"github.com/vanamelnik/gophkeeper/models"
	pb "github.com/vanamelnik/gophkeeper/proto"
)

var (
	// ErrQuit is thrown when the user chooses to quit the application.
	ErrQuit = errors.New("user chose to quit the application")
)

const (
	cmdNewItem = "n"
	cmdBack    = "b"
	cmdLogout  = "x"
	cmdQuit    = "q"
	cmdSave    = "S"
	cmdDelete  = "D"
)

var (
	selectItemQuit      = selectItem{cmdQuit, "Quit"}
	selectItemLogout    = selectItem{cmdLogout, "Log out"}
	selectItemBack      = selectItem{cmdBack, "Back to previous menu"}
	selectItemSave      = selectItem{cmdSave, "Save changes"}
	selectItemSeparator = selectItem{}
)

type (
	UserInterface struct {
		ctx      context.Context
		repo     *repo.Repo
		c        *client.Client
		pbClient pb.GophkeeperClient
	}

	selectItem struct {
		ID   selectID
		Text string
	}

	selectID string
)

func NewUI(ctx context.Context, repo *repo.Repo, c *client.Client, pbClient pb.GophkeeperClient) UserInterface {
	return UserInterface{
		ctx:      ctx,
		repo:     repo,
		c:        c,
		pbClient: pbClient,
	}
}

// Run is the application's main loop.
func (ui *UserInterface) Run() {
	fmt.Println("GophKeeper client")

	for {
		if !ui.c.IsLoggedIn() {
			if err := ui.signInView(); err != nil {
				if errors.Is(err, ErrQuit) {
					return
				}
				log.Println(err)
				return
			}
		}
		err := ui.clientSession()
		if errors.Is(err, client.ErrReloginNeeded) {
			// nolint: errcheck
			ui.c.LogOut()
			continue
		}
		if errors.Is(err, client.ErrSessionInactive) {
			continue
		}
		if err != nil {
			log.Println(err)
		}
		return
	}
}

// signInView perform login or signup procedure.
// The Access and Refresh tokens are stored in the local repository.
// Returns ErrQuit if the user chose to quit.
func (ui *UserInterface) signInView() error {
	const (
		cmdSignUp = "1"
		cmdLogin  = "2"
	)
	for {
		choice := choose(`What do you want to do?`, []selectItem{
			{cmdSignUp, "Register a new user"},
			{cmdLogin, "Log in existing user"},
			selectItemQuit,
		})
		switch choice {
		case cmdSignUp:
			if err := ui.signUp(); err != nil {
				fmt.Printf("Could not sign up: %v\n", err)
				continue
			}
			return nil
		case cmdLogin:
			if err := ui.logIn(); err != nil {
				fmt.Printf("Could not log in: %v\n", err)
				continue
			}
			return nil
		case cmdQuit:
			return ErrQuit
		}
	}
}

// signUp registers a new user, creates a new user session and stores the token pair
// in the local repository.
func (ui *UserInterface) signUp() error {
	email, password := getCredentials()
	if err := ui.c.SignUp(email, password); err != nil {
		return err
	}
	fmt.Println("Welcome, new user!")
	return nil
}

// logIn authenticates the user on the server and stores a token pair in the local repository..
func (ui *UserInterface) logIn() error {
	email, password := getCredentials()
	if err := ui.c.LogIn(email, password); err != nil {
		return err
	}
	fmt.Println("Welcome back!")
	return nil

}

func getCredentials() (email, password string) {
	for {
		fmt.Print("login (e-mail address): ")
		email = inputWord()
		// TODO: validate email with regexp
		fmt.Print("password: ")
		var err error
		password, err = inputPassword()
		if err != nil {
			fmt.Printf("Wrong input: %s Try again.", err)
			continue
		}
		// TODO: validate the password
		return
	}
}

// choose prints the message to the stdout and waits for the user to enter a digit or a letter matching the selected item.
//
// Example:
//	What do you want to do?
//		1	Sign up
//		2	Log In
//		q	Quit
func choose(message string, choices []selectItem) selectID {
	if message != "" {
		fmt.Println(message)
	}
	for _, element := range choices {
		fmt.Printf("\t%s\t%s\n", element.ID, element.Text)
	}
	for {
		fmt.Print(">>")
		in := inputWord()
		for _, element := range choices {
			if in == string(element.ID) {
				return element.ID
			}
		}
		fmt.Println("Syntax error, try again")
	}
}

func inputPassword() (string, error) {
	// TODO: input password without echo
	return inputString(), nil
}

func inputString() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// inputWord scans a single word from stdin. If several words are entered, all but the first words are ignored.
func inputWord() string {
	str := inputString()
	if str == "" {
		return str
	}
	return strings.Fields(str)[0]
}

func ResolveConflict(recievedItem models.Item, localEntry repo.Entry) (userChoseReceivedItem bool) {
	//TODO: ...
	return true
}
