package email

import (
	"strconv"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

type ItemUpdate struct {
	CategoryUpdate string
	GenderUpdate   string
	QuantityUpdate string
	SizeUpdate     string
	StatusUpdate   string
}

func BuildItemUpdate(previousItem *managers.Item, updatedItem *managers.Item) *ItemUpdate {
	itemUpdate := &ItemUpdate{}
	if previousItem.Category != updatedItem.Category {
		itemUpdate.CategoryUpdate = previousItem.Category + " -> " + updatedItem.Category
	}

	if previousItem.Gender != updatedItem.Gender {
		itemUpdate.GenderUpdate = previousItem.Gender + " -> " + updatedItem.Gender
	}

	if previousItem.Quantity != updatedItem.Quantity {
		itemUpdate.QuantityUpdate = strconv.Itoa(int(previousItem.Quantity)) + " -> " + strconv.Itoa(int(updatedItem.Quantity))
	}

	if previousItem.Size != updatedItem.Size {
		itemUpdate.SizeUpdate = previousItem.Size + " -> " + updatedItem.Size
	}

	if previousItem.Status != updatedItem.Status {
		itemUpdate.StatusUpdate = retrievers.StatusAsString(previousItem.Status) + " -> " + retrievers.StatusAsString(updatedItem.Status)
	}

	return itemUpdate
}

func formatEmailBody(itemUpdate *ItemUpdate) string {
	emailBody := "Item Upate Details: \n"
	if itemUpdate.CategoryUpdate != "" {
		emailBody = emailBody + "Category: " + itemUpdate.CategoryUpdate + "\n"
	}
	if itemUpdate.GenderUpdate != "" {
		emailBody = emailBody + "Gender: " + itemUpdate.GenderUpdate + "\n"
	}
	if itemUpdate.CategoryUpdate != "" {
		emailBody = emailBody + "Category: " + itemUpdate.QuantityUpdate + "\n"
	}
	if itemUpdate.CategoryUpdate != "" {
		emailBody = emailBody + "Size: " + itemUpdate.SizeUpdate + "\n"
	}
	if itemUpdate.CategoryUpdate != "" {
		emailBody = emailBody + "Status: " + itemUpdate.StatusUpdate + "\n"
	}

	return emailBody
}
