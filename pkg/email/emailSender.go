package email

import (
	"context"
	"fmt"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	DeliverEmail(ctx context.Context, previousItem *managers.Item, currentItem *managers.Item, userSession *managers.UserSession) error
}

type LocalSender struct {
	Dialer      *gomail.Dialer
	UserManager *managers.UserManager
}

type SendGridSender struct {
	Client      *sendgrid.Client
	UserManager *managers.UserManager
}

func (ls *LocalSender) DeliverEmail(ctx context.Context, previousItem *managers.Item, currentItem *managers.Item, userSession *managers.UserSession) error {
	var recipient *managers.User
	var err error
	if userSession.UserType == managers.SHELTER {
		if currentItem.SamaritanID < 1 {
			return fmt.Errorf("No samaritan to update")
		}

		recipient, err = ls.UserManager.GetUser(ctx, currentItem.SamaritanID)
		if err != nil {
			return err
		}
	} else {
		recipient, err = ls.UserManager.GetUser(ctx, currentItem.ShelterID)
		if err != nil {
			return err
		}
	}

	updater, err := ls.UserManager.GetUser(ctx, userSession.UserID)
	if err != nil {
		return err
	}

	updater.Email = "kwhite@hubspot.com"
	itemUpdate := BuildItemUpdate(previousItem, currentItem, recipient, updater)
	return ls.sendEmail(itemUpdate)
}

func (ls *LocalSender) sendEmail(itemUpdate *ItemUpdate) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", itemUpdate.Updater.Email, "kwhite")
	m.SetAddressHeader("To", itemUpdate.Recipient.Email, itemUpdate.Recipient.Name)
	m.SetHeader("Subject", "Item Updated by "+itemUpdate.Updater.Email+"!")
	m.SetBody("text/plain", formatEmailBody(itemUpdate))

	return ls.Dialer.DialAndSend(m)
}

func (ss *SendGridSender) DeliverEmail(ctx context.Context, previousItem *managers.Item, currentItem *managers.Item, userSession *managers.UserSession) error {
	var recipient *managers.User
	var err error
	if userSession.UserType == managers.SHELTER {
		if currentItem.SamaritanID < 1 {
			return fmt.Errorf("No samaritan to update")
		}

		recipient, err = ss.UserManager.GetUser(ctx, currentItem.SamaritanID)
		if err != nil {
			return err
		}
	} else {
		recipient, err = ss.UserManager.GetUser(ctx, currentItem.ShelterID)
		if err != nil {
			return err
		}
	}

	updater, err := ss.UserManager.GetUser(ctx, userSession.UserID)
	if err != nil {
		return err
	}

	itemUpdate := BuildItemUpdate(previousItem, currentItem, recipient, updater)
	return ss.sendEmail(itemUpdate)
}

func (ss *SendGridSender) sendEmail(itemUpdate *ItemUpdate) error {
	from := mail.NewEmail(itemUpdate.Updater.Name, itemUpdate.Updater.Email)
	to := mail.NewEmail(itemUpdate.Recipient.Name, itemUpdate.Recipient.Email)
	plainTextContent := formatEmailBody(itemUpdate)
	htmlContent := "<div>" + plainTextContent + "</div>"
	message := mail.NewSingleEmail(from, "Item Updated by "+itemUpdate.Updater.Name+"!", to, plainTextContent, htmlContent)
	_, err := ss.Client.Send(message)
	return err
}