package email

import (
	"strconv"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"github.com/kwhite17/Neighbors/pkg/retrievers"
)

type ItemUpdate struct {
	PreviousItem   *managers.Item
	CategoryUpdate string
	GenderUpdate   string
	QuantityUpdate string
	SizeUpdate     string
	StatusUpdate   string
	Recipient      *managers.User
	Updater        *managers.User
}

type PasswordReset struct {
	Recipient    *managers.User
	TempPassword string
}

func BuildPasswordReset(recipient *managers.User, tempPassword string) *PasswordReset {
	return &PasswordReset{Recipient: recipient, TempPassword: tempPassword}
}

func BuildItemUpdate(previousItem *managers.Item, updatedItem *managers.Item, recipient *managers.User, updater *managers.User) *ItemUpdate {
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

	itemUpdate.Recipient = recipient
	itemUpdate.Updater = updater
	itemUpdate.PreviousItem = previousItem
	return itemUpdate
}

func formatEmailBody(itemUpdate *ItemUpdate) string {
	emailBody := "Updates on current request for: " + strconv.Itoa(int(itemUpdate.PreviousItem.Quantity)) +
		" " + itemUpdate.PreviousItem.Category + "\n"
	if itemUpdate.CategoryUpdate != "" {
		emailBody = emailBody + "Category: " + itemUpdate.CategoryUpdate + "\n"
	}
	if itemUpdate.GenderUpdate != "" {
		emailBody = emailBody + "Gender: " + itemUpdate.GenderUpdate + "\n"
	}
	if itemUpdate.QuantityUpdate != "" {
		emailBody = emailBody + "Quantity: " + itemUpdate.QuantityUpdate + "\n"
	}
	if itemUpdate.SizeUpdate != "" {
		emailBody = emailBody + "Size: " + itemUpdate.SizeUpdate + "\n"
	}
	if itemUpdate.StatusUpdate != "" {
		emailBody = emailBody + "Status: " + itemUpdate.StatusUpdate + "\n"
	}

	return emailBody
}

func formatPasswordResetEmailBody(passwordReset *PasswordReset) string {
	return "Hello " + passwordReset.Recipient.Name + ",\n\n" + "We've reset your password as requested to: " + passwordReset.TempPassword + ". " +
		"We recommend you change this as soon as possible. Have a nice day!"
}
