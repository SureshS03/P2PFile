package main

import (
	"errors"
	"fmt"
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

func FileAlreadyExits(data *MetaData, id *string) (int, error) {
	var r int
	for _, file := range data.Files {
		if file.Id == *id {
			fmt.Println("File is already Exits\nDo You wanna do it again?")
			fmt.Println("YES means: 1\nNO means: 2")
			_, err := fmt.Scanln(&r)
			if err != nil {
				return -1, err
			}
			return r, nil
		}
	}
	return 0, nil
}

func SignUp() (string, string,error) {
	//TODO add sign up msg
	fmt.Println("Enter Your Email :")
	var email string
	_, err := fmt.Scanln(&email)
	if err != nil {
		fmt.Println("in scan", err)
		return "", "", err
	}
	fmt.Println("Email :", email)
	fmt.Println("Enter Your Email Password :")
	var pass string
	x, err := fmt.Scanln(&pass)
	fmt.Println("Length of the password is:",x)
	if err != nil {
		fmt.Println("in scan", err)
		return "", "", err
	}
	fmt.Println("Password :", pass)
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
