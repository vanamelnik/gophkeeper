package repo

import "github.com/vanamelnik/gophkeeper/models"

type ErrMergeConflict struct {
	LocalEntry Entry
}

func (err ErrMergeConflict) Error() string {
	return "merge conflict with a local entry"
}

// MergeItems performs merging the items received from the server in accordance with the rules (see README.md).
func (r *Repo) MergeItem(receivedItem models.Item) error {
	for i, entry := range r.entries {
		if entry.Item.ID == receivedItem.ID {
			// merge condition boolean variables
			remoteVersionIsNewer := receivedItem.Version != entry.Item.Version // just in case we use '!=' instead of '>'
			confirmedPendingItem := entry.Pending && compareItemsData(receivedItem, entry.Item)
			itemChangedRemotely := !entry.Pending

			if remoteVersionIsNewer &&
				(confirmedPendingItem || itemChangedRemotely) {
				// replace local item with received one
				r.entries[i] = Entry{
					Item:    receivedItem,
					Pending: false,
				}
				return nil
			}

			// all other cases - we have a merge conflict
			return ErrMergeConflict{entry}
		}
	}
	// item not found - merge it!
	r.entries = append(r.entries, Entry{
		Item:    receivedItem,
		Pending: false,
	})
	return nil
}

// ForceMergeItem is used to merge items received from the server into the local repository.
// The method replaces the local item by the item provided or creates the new entry in the local repository.
// 'Pending' flag is unset.
// Contract: repo must be locked.
func (r *Repo) ForceMergeItem(item models.Item) {
	r.RLock()
	defer r.RUnlock()
	for i, entry := range r.entries {
		if entry.Item.ID == item.ID {
			r.entries[i].Item = item
			r.entries[i].Pending = false
			return
		}
	}
	// if there are no entry with such item, create it.
	r.entries = append(r.entries, Entry{
		Item:    item,
		Pending: false,
	})
}

// compareItemsData return true if the payload, DeleteAt fields and the Meta fields are the same for both items.
func compareItemsData(item1, item2 models.Item) bool {
	return item1.Payload == item2.Payload &&
		item1.Meta == item2.Meta &&
		item1.DeletedAt == item2.DeletedAt
}
