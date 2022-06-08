package ui

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/client"
)

func (ui *UserInterface) passwordsView() error {
	for {
		passwords, err := ui.c.GetPasswords()
		if err != nil {
			return err
		}
		if len(passwords) == 0 {
			ui.noPasswords()
			continue
		}

		fmt.Printf("\nYou have %d passwords in the storage.\n", len(passwords))
		choices := make([]selectItem, len(passwords))
		for i := range passwords {
			choices[i] = selectItem{
				ID: selectID(fmt.Sprint(i + 1)),
				Text: fmt.Sprintf("%s\t*****\t%v\t%s",
					passwords[i].Login, passwords[i].CreatedAt, passwords[i].Notes), // TODO: implement formatted table output
			}
		}
		choices = append(choices, selectItemSeparator, selectItem{cmdNewItem, "Create a new password"}, selectItemBack)
		usersChoice := choose("Select the password to view, modify or delete, 'n' to create a new password "+
			"or 'b' to return to the previous menu:", choices)
		n, err := strconv.Atoi(string(usersChoice))
		if err != nil {
			switch usersChoice {
			case cmdNewItem:
				if err := ui.newPassword(); err != nil {
					fmt.Println(err)
					continue
				}
				continue
			case cmdBack:
				return nil
			}
		}
		ui.workWithPassword(passwords[n])
	}
}

func (ui *UserInterface) noPasswords() {
	const (
		cmdNewPassword = "1"
	)
	for {
		choice := choose("You have no passwords in the storage. Select what do you want to do:", []selectItem{
			{cmdNewPassword, "Create a new password"},
			selectItemBack,
		})
		if choice == cmdBack {
			return
		}
		err := ui.newPassword()
		if err == nil {
			return
		}
		fmt.Printf("Could not create the password: %s\n", err)
	}
}

func (ui UserInterface) newPassword() error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("could not generate password id: %w", err)
	}
	p := client.Password{
		ID:        id,
		Version:   0,
		CreatedAt: time.Now(),
	}
	fmt.Println("Create a new password")
	fmt.Print("Enter login: ")
	p.Login = inputWord()
	pwd, err := enterNewPassword()
	if err != nil {
		return err
	}
	p.Password = pwd

	fmt.Print("Enter text notes: ")
	p.Notes = inputString()

	return ui.c.CreatePassword(p)
}

func (ui *UserInterface) workWithPassword(p client.Password) {
	isChanged := false
	for {
		fmt.Printf("Login:\t%s\n", p.Login)
		fmt.Printf("Password:\t%s\n", p.Password)
		fmt.Printf("Created:\t%v\n", p.CreatedAt)
		fmt.Printf("Notes:\t%s\n", p.Notes)
		fmt.Println()
		const (
			cmdModifyLogin    = "1"
			cmdModifyPassword = "2"
			cmdModifyNotes    = "3"
		)
		choices := []selectItem{
			{cmdModifyLogin, "Modify login"},
			{cmdModifyPassword, "Modify password"},
			{cmdModifyNotes, "Modify notes"},
			{cmdDelete, "Delete password"},
		}
		if isChanged {
			choices = append(choices, selectItemSeparator, selectItemSave)
		}
		choices = append(choices, selectItemSeparator, selectItemBack)
		choice := choose("Select operation with your password:", choices)
		switch choice {
		case cmdModifyLogin:
			fmt.Print("Enter new login: ")
			p.Login = inputWord()
		case cmdModifyPassword:
			fmt.Print("Enter new password: ")
			password, err := enterNewPassword()
			if err != nil {
				continue
			}
			p.Password = password
		case cmdModifyNotes:
			fmt.Print("Enter new notes: ")
			p.Notes = inputWord()
		case cmdDelete:
			if err := ui.c.DeletePassword(p); err != nil {
				fmt.Printf("Could not delete the password: %s", err)
				continue
			}
			fmt.Println("Password has been successfully deleted.")
			return
		case cmdSave:
			if err := ui.c.UpdatePassword(p); err != nil {
				fmt.Printf("Could not update the password: %s", err)
				continue
			}
		case cmdBack:
			return
		}
	}
}

func enterNewPassword() (string, error) {
	fmt.Print("Enter password: ")
	password1, err := inputPassword()
	if err != nil {
		return "", err
	}
	fmt.Print("Re-enter password: ")
	password2, err := inputPassword()
	if err != nil {
		return "", err
	}
	if password1 != password2 {
		return "", errors.New("passwords doesn't equal")
	}
	return password1, nil
}
