package email

import (
	"context"
	"fmt"
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"github.com/kwhite17/Neighbors/pkg/managers"
	"gopkg.in/gomail.v2"
)

const SENDGRID_SENDER_EMAIL = "neighbors@massally.org"

type EmailSender interface {
	DeliverEmail(ctx context.Context, previousItem *managers.Item, currentItem *managers.Item, userSession *managers.UserSession) error
	DeliverPasswordResetEmail(ctx context.Context, user *managers.User, temporaryPassword string) error
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

func (ls *LocalSender) DeliverPasswordResetEmail(ctx context.Context, recipient *managers.User, temporaryPassword string) error {
	passwordReset := BuildPasswordReset(recipient, temporaryPassword)
	return ls.sendPasswordResetEmail(passwordReset)
}

func (ls *LocalSender) sendPasswordResetEmail(passwordReset *PasswordReset) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", "kwhite@hubspot.com", "kwhite")
	m.SetAddressHeader("To", passwordReset.Recipient.Email, passwordReset.Recipient.Name)
	m.SetHeader("Subject", "Neighbors Password Reset")
	m.SetBody("text/plain", formatPasswordResetEmailBody(passwordReset))

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

func (ss *SendGridSender) DeliverPasswordResetEmail(ctx context.Context, recipient *managers.User, temporaryPassword string) error {
	passwordReset := BuildPasswordReset(recipient, temporaryPassword)
	return ss.sendPasswordResetEmail(passwordReset)
}

func (ss *SendGridSender) sendPasswordResetEmail(passwordReset *PasswordReset) error {
	from := mail.NewEmail("Neighbors", SENDGRID_SENDER_EMAIL)
	to := mail.NewEmail(passwordReset.Recipient.Name, passwordReset.Recipient.Email)
	plainTextContent := formatPasswordResetEmailBody(passwordReset)
	htmlContent := "<div>" + plainTextContent + "</div>"
	message := mail.NewSingleEmail(from, "Neighbors Password Reset", to, plainTextContent, htmlContent)
	response, err := ss.Client.Send(message)
	if response.StatusCode > 299 {
		log.Println(response)
	}
	return err
}

func (ss *SendGridSender) sendEmail(itemUpdate *ItemUpdate) error {
	from := mail.NewEmail(itemUpdate.Updater.Name+" ("+itemUpdate.Updater.Email+") via "+SENDGRID_SENDER_EMAIL, SENDGRID_SENDER_EMAIL)
	to := mail.NewEmail(itemUpdate.Recipient.Name, itemUpdate.Recipient.Email)
	plainTextContent := formatEmailBody(itemUpdate)
	htmlContent := "<div>" + plainTextContent + "</div>"
	message := mail.NewSingleEmail(from, "Item Updated by "+itemUpdate.Updater.Name+"!", to, plainTextContent, htmlContent)
	response, err := ss.Client.Send(message)
	if response.StatusCode > 299 {
		log.Println(response)
	}
	return err
}
