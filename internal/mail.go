package internal

import (
	"fmt"
	"net/smtp"
)

func sendMail(id *string, data *MetaData, to string) error {
	fmt.Println(*id)
	if data.Mail == "" {
		signUp()
	}
	fmt.Println(data.Mail,data.Pass)
	var file FileMetaData
	for i := 0; i < len(data.Files); i++ {
		if data.Files[i].Id == *id {
			file = data.Files[i]
		}
	}
	auth := smtp.PlainAuth("", data.Mail, data.Pass, "smtp.gmail.com")
	body := addMIME(file, data.Mail, to, file.Id)
	err := smtp.SendMail("smtp.gmail.com:587", auth, data.Mail, []string{to}, body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("done")
	return nil
}
