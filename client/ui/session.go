package ui

import (
	"errors"
	"fmt"

	"github.com/vanamelnik/gophkeeper/client"
)

func (ui *UserInterface) clientSession() error {
	// Log in or sign up if needed
	if err := ui.c.NewSession(); err != nil {
		return err
	}
	defer ui.c.CloseClientSession()

	const (
		cmdViewPasswords = "1"
		cmdViewTexts     = "2"
		cmdViewCards     = "3"
		cmdViewBlobs     = "4"
	)
	for {
		action := choose("Choose what you want to work with:", []selectItem{
			{cmdViewPasswords, "Passwords"},
			{cmdViewTexts, "Text notes"},
			{cmdViewCards, "Credit cards"},
			{cmdViewBlobs, "Binary data"},
			selectItemSeparator,
			selectItemLogout,
			selectItemQuit,
		})
		switch action {
		case cmdViewPasswords:
			if err := ui.passwordsView(); err != nil {
				if errors.Is(err, client.ErrSessionInactive) {
					return client.ErrReloginNeeded
				}
				return err
			}
		case cmdViewTexts:
		case cmdViewCards:
		case cmdViewBlobs:
		case cmdLogout:
			if err := ui.c.LogOut(); err != nil {
				fmt.Printf("Error while trying to log out: %s", err)
			}
			return nil
		case cmdQuit:
			return ErrQuit
		}
	}
}
