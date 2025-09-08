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

func signUp() (string, string,error) {
	fmt.Println("Enter Your Email :")
	var email string
	_, err := fmt.Scanln(&email)
	if err != nil {
		fmt.Println("in scan", err)
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
	if len(email) == 0 {
		return "", "", errors.New("bro, Its empty")
	}
	if !strings.Contains(email, "@") {
		return "", "", errors.New("bro, Its dont have @ in it")
	} else if !strings.Contains(email, ".") {
		return "", "", errors.New("bro, Its not an vaild email")
	} else {
		err := JsonWriter("MetaData.json", MetaData{
			NumOfFiles: 0,
		})
		if err != nil {
			return "", "", err
		}
		return email, pass, nil
	}
}
