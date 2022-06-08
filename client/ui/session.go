package ui

import (
	"errors"

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
					aaaaaaaaaaaa
				}
			}
		case cmdViewTexts:
		case cmdViewCards:
		case cmdViewBlobs:
		case cmdLogout:
			return ui.c.LogOut()
		case cmdQuit:
		}

		return client.ErrReloginNeeded
	}
}
