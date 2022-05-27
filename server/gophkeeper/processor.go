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
			s.Close()
			break processorLoop
		case p := <-s.eventCh:
			if err := s.processUserData(p.ctx, p.userID, p.events); err != nil {
				log.Printf("gophkeeper processor: %s", err)
			}
		}
	}
	if s.eventCh != nil {
		close(s.eventCh)
		s.eventCh = nil
	}
	log.Println("GophKeeper processor is stopped")
}

func (s Service) processUserData(ctx context.Context, userID uuid.UUID, events []models.Event) error {
	tx, err := s.storage.NewUserTransaction(ctx, userID)
	if err != nil {
		return err
	}
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

	return nil
}
