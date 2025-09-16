package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

func isFileExits(data *MetaData, id *string) error {
	for _, file := range data.Files {
		if file.Id == *id {
			return nil
		}
	}
	return errors.New("bro, file ID not exists")
}

//after added cobra need to make this in Yes/No
func fileAlreadyExits(data *MetaData, id *string) (int, error) {
	var r int
	for _, file := range data.Files {
		if file.Id == *id {
			WarnPrinter("File is already Exits\nDo You wanna do it again?")
			WarnPrinter("YES means: 1\nNO means: 2")
			_, err := fmt.Scanln(&r)
			if err != nil {
				return -1, err
			}
			return r, nil
		}
	}
	return 0, nil
}

func signUp() (string, string, error) {
	fmt.Printf("Welcome to...\n" + ` ____  ____  ____ 
(  _ \(___ \(  _ \
 ) __/ / __/ ) __/
(__)  (____)(__)  
` + "\n")
	fmt.Println("Enter Your Email :")
	var email string
	_, err := fmt.Scanln(&email)
	if err != nil {
		return "", "", errors.New("error in reading mail")
	}
	err = IsValidMail(email)
	if err != nil {
		return "", "", err
	}
	CrrPrinter("Email :" + email)
	passReader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter Your Email Password :")
	pass, err := passReader.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("error at reading password bro")
	}
	pass = strings.TrimSpace(pass)
	CrrColorString("Password :"+ pass)
	err = JsonWriter("MetaData.json", MetaData{
		NumOfFiles: 0,
	})
	if err != nil {
		return "", "", err
	}
	return email, pass, nil
}

func IsValidMail(mail string) error {
	if len(mail) == 0 {
		return errors.New("bro, mail is empty\nPlease Enter vaild mail")
	}
	if !strings.Contains(mail, "@") {
		return errors.New("bro, mail dont have @ in it\nPlease Enter vaild mail")
	}
	if !strings.Contains(mail, ".") {
		return errors.New("bro, mail not an vaild email\nPlease Enter vaild mail")
	}

	return nil
}