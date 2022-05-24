package repo

import "github.com/vanamelnik/gophkeeper/models"

// MergeItems performs merging the items received from the server in accordance with the rules (see README.md).
func (r *Repo) MergeItems(dataVersion uint64, items []models.Item) {
	r.RLock()
	defer r.RUnlock()
	for _, receivedItem := range items {
		for i, entry := range r.entries {
			if entry.Item.ID == receivedItem.ID {
				if receivedItem.Version != entry.Item.Version &&
					(entry.Pending && compareItemsData(receivedItem, entry.Item) || !entry.Pending) {
					r.entries[i] = Entry{
						Item:    receivedItem,
						Pending: false,
					}
					continue
				}
				if r.resolveConflict(receivedItem, entry) {
					r.entries[i].Item.Version = receivedItem.Version
					r.entries[i].Pending = true

				}

			}
		}
	}
}

func compareItemsData(item1, item2 models.Item) bool {
	return item1.Payload == item2.Payload &&
		item1.Meta == item2.Meta &&
		item1.DeletedAt == item2.DeletedAt
}
