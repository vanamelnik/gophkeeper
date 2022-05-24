package repo

import "github.com/vanamelnik/gophkeeper/models"

// MergeItems performs merging the items received from the server in accordance with the rules (see README.md).
func (r *Repo) MergeItems(dataVersion uint64, items []models.Item) {
	r.RLock()
	defer r.RUnlock()
	for _, receivedItem := range items {

	}
}
