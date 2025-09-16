package internal

import (
	"errors"
	"fmt"
	"net/smtp"
)

func sendMail(id *string, data *MetaData, to string) error {
	if data.Mail == "" {
		WarnPrinter("Mail id not found, Please Sign up")
		signUp()
	}
	var file FileMetaData
	for i := 0; i < len(data.Files); i++ {
		if data.Files[i].Id == *id {
			file = data.Files[i]
		}
	}
	body, err := addMIME(file, data.Mail, to, file.Id)
	if err != nil {
		return err
	}
	auth := smtp.PlainAuth("", data.Mail, data.Pass, "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, data.Mail, []string{to}, body)
	if err != nil {
		return errors.New("error in Sending mail:\n" + fmt.Sprint(err.Error()) )
	}
	return nil
}
