package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"

	//"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	//splitFile("test.mp4")
	metadata, err := getMetaData("MetaData.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(metadata)
	if metadata.Mail == "" {
		metadata.Mail, err = SignUp()
		if err != nil {
			return
		}
		err = JsonWriter("MetaData.json", metadata)
		if err != nil {
			return
		}
	}

	/*
		var name [4]string = [4]string{
			"./test_part0.chu",
			"./test_part1.chu",
			"./test_part2.chu",
			"./test_part3.chu",
		}
		_, err, key := makeKey()
		if err != nil {

		}
		fileName := "op.mp4"
		_, err = os.Create(fileName)
		f, err := os.Open(fileName)
		if err != nil {
			return
		}
		const chunkSize = 1024 * 1024 * 20 // 1024 * 1024 is 1 mb so 20 mb

		for i := 0; i < 4; i++ {
			path := name[i]
			x := makeDec(key, path)
			_, err := f.WriteAt(x, int64(i*chunkSize))
			if err != nil {
				return
			}
		}
	*/
}

func splitFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	fileNamer := info.Name()
	fileSize := info.Size()
	fmt.Println("name and size", fileNamer, fileSize)

	const chunkSize = 1024 * 1024 * 20 // 1024 * 1024 is 1 mb so 20 mb

	needChunks := (fileSize + chunkSize - 1) / chunkSize
	buffer := make([]byte, chunkSize)

	defaultName := filepath.Base(path)
	fileExt := filepath.Ext(path)
	Name := defaultName[:len(defaultName)-len(fileExt)]

	fmt.Println(Name, fileExt)
	gcm, err, _ := makeKey()
	if err != nil {
		fmt.Println(err)
	}

	for i := int64(0); i < needChunks; i++ {
		fileNamer := Name + "_part" + fmt.Sprint(i) + ".chu"
		fmt.Println(fileNamer)
		_, err := file.Seek(i*chunkSize, 0)
		if err != nil {
			fmt.Println(err)
			return
		}
		ReadedBytes, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return
		}
		fmt.Println(ReadedBytes/1024/1024, "MB")
		//fmt.Println(string(buffer[:ReadedBytes]))
		text := makeEnc(gcm, buffer[:ReadedBytes])
		//fmt.Println("text", text)
		hash := sha256.New()
		hash.Write([]byte(text))
		fmt.Println(hash.Sum(nil))
		partFile, err := os.Create(fileNamer)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = partFile.Write(text)
		if err != nil {
			partFile.Close()
			fmt.Println(err)
			return
		}
		partFile.Close()
	}
}

func getMetaData(path string) (MetaData, error) {
	data, err := JsonReader(path)
	if err != nil {
		fmt.Println(err)
		return MetaData{}, err
	}
	var metaData MetaData
	err = json.Unmarshal(data, &metaData)
	if err != nil {
		fmt.Println(err)
	}
	return metaData, nil
	/*
		metaData.Mail = "test@gmail.com"
		metaData.NumOfFiles = 1
		err = JsonWriter(path, metaData)
		if err != nil {
			return
		}
	*/
}

func SignUp() (string, error) {
	fmt.Println("Enter Your Email :")
	var email string
	_, err := fmt.Scanln(&email)
	if err != nil {
		return "", err
	}
	fmt.Println("Email :", email)
	if email == "" {
		return "", errors.New("bro, Its empty")
	}
	for i := 0; i < len(email); i++ {
		if email[i] != '@' {
			return "", errors.New("bro, Its don't Have '@' in it")
		}
	}
	return email, nil
}
