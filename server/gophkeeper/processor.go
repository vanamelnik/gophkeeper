package gophkeeper

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/vanamelnik/gophkeeper/models"
)

// processor collects event packages from the event channel and sends them to the storage
func (s Service) processor() {
processorLoop:
	for {
		select {
		case <-s.stopCh:
			break processorLoop
		case p := <-s.eventCh:
			go func(p eventsPack) {
				if err := s.processUserData(p.ctx, p.userID, p.events); err != nil {
					log.Printf("gophkeeper processor: %s", err)
				}
			}(p)
		}
	}
	log.Println("GophKeeper processor is stopped")
}

func (s Service) processUserData(ctx context.Context, userID uuid.UUID, events []models.Event) error {
	defer s.wg.Done()

	tx, err := s.storage.NewUserTransaction(ctx, userID)
	if err != nil {
		return err
	}
	// nolint: errcheck
	defer tx.RollBack()

	for _, event := range events {
		switch event.Operation {
		case models.OpCreate:
			if err := tx.CreateItem(ctx, event.Item); err != nil {
				return err
			}
		case models.OpUpdate:
			if err := tx.UpdateItem(ctx, event.Item); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}
