package client

import (
	"errors"

	"github.com/vanamelnik/gophkeeper/client/repo"
	"github.com/vanamelnik/gophkeeper/models"
)

func (c *Client) MergeItems(dataVersion uint64, items []models.Item) {
	for _, item := range items {
		err := c.repo.MergeItem(item)
		if err != nil {
			c.processConflictResolving(item, err)
		}
	}
	c.repo.StoreDataVersion(dataVersion)
}

func (c *Client) processConflictResolving(receivedItem models.Item, err error) {
	if err == nil {
		return
	}
	var ce repo.ErrMergeConflict
	if errors.As(err, &ce) {
		// ask the user which item he chooses
		if c.conflictResolveFn(receivedItem, ce.LocalEntry) { // if true - the user chose to merge the received item
			c.repo.ForceMergeItem(receivedItem)
			return
		}
		// else use the local item, but replace the version number and publish it again.
		mergedItem := ce.LocalEntry.Item
		mergedItem.Version = receivedItem.Version
		c.repo.ForceMergeItem(mergedItem)
		c.PublishEvent(models.Event{ // send this item again
			Operation: models.OpUpdate,
			Item:      mergedItem,
		})
	}
}
