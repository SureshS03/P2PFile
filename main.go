package main

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"strings"
	"time"

	//"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	//splitFile("test.mp4")
	args := os.Args
	fmt.Println(args[1])
	switch args[1] {
	case "add":
		if len(args) < 3 {
			fmt.Println("Bro Need File: add <filename>")
			break
		} else {
			err := splitFile(args[2])
			if err != nil {
				fmt.Println(err)
				return
			}
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

func splitFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return errors.New(fmt.Sprintf("Bro cant open file %s", err))
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return errors.New(fmt.Sprintf("Bro cant access file %s", err))
	}

	metadata, err := getMetaData("MetaData.json")
	if err != nil {
		fmt.Println(err)
		return errors.New(fmt.Sprintf("Bro cant access MetaData File %s", err))
	}

	if metadata.Mail == "" {
		mail, err := SignUp()
		if err != nil {
			return errors.New(fmt.Sprintf("Bro cant sign Up %s", err))
		}
		fmt.Println("Mail added", mail)
	}

	NumOfFile := metadata.NumOfFiles

	fileName := info.Name()
	fileSize := info.Size()
	fmt.Println("name and size", fileName, fileSize)

	const chunkSize = 1024 * 1024 * 20 // 1024 * 1024 is 1 mb so 20 mb

	needChunks := (fileSize + chunkSize - 1) / chunkSize
	buffer := make([]byte, chunkSize)

	chunks := []ChunkMetaData{}
	fileID := fmt.Sprintf("%dchu%v", fileName, fileSize)
	r, err := FileAlreadyExits(&metadata, &fileID)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if r == 2 {
		return errors.New("STOPPED...file Already Exit")
	}
	if r == 1 {
		fmt.Println("Doing Again...")
		p := &fileID
		*p = fmt.Sprintf("%dchuCOPY%v", fileName, fileSize)
	}

	defaultName := filepath.Base(path)
	fileExt := filepath.Ext(path)
	Name := defaultName[:len(defaultName)-len(fileExt)]

	fmt.Println(Name, fileExt)
	gcm, err, key := makeKey()
	if err != nil {
		fmt.Println(err)
		return errors.New(fmt.Sprintf("Bro cant make key %s", err))
	}

	for i := int64(0); i < needChunks; i++ {
		fileNamer := Name + "_part" + fmt.Sprint(i) + ".chu"
		fmt.Println(fileNamer)
		_, err := file.Seek(i*chunkSize, 0)
		if err != nil {
			fmt.Println(err)
			return errors.New(fmt.Sprint("Bro, cant seek into the file:", err))
		}
		ReadedBytes, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println(err)
			return errors.New(fmt.Sprint("Bro, cant read the file:", err))
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
			return errors.New(fmt.Sprint("Bro, cant create the file chunks:", err))
		}
		_, err = partFile.Write(text)
		if err != nil {
			partFile.Close()
			fmt.Println(err)
			return errors.New(fmt.Sprint("Bro, cant write the file chunks:", err))
		}
		partFile.Close()
		chunks = append(chunks, ChunkMetaData{
			ChunkName: fileNamer,
			ChunkSize: fmt.Sprintf("%d MB", ReadedBytes/1024/1024),
		})
	}
	metadata.NumOfFiles = NumOfFile + 1
	filemetadata := FileMetaData{
		Id:          fileID,
		FileName:    fileName,
		TotalSize:   fmt.Sprintf("%d MB", fileSize/1024/1024),
		NumOfChunks: needChunks,
		Key:         key,
		Chunks:      chunks,
		CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
	}
	metadata.Files = append(metadata.Files, filemetadata)
	err = JsonWriter("MetaData.json", metadata)
	if err != nil {
		fmt.Println("Bro cant add Metadata details:", err)
		fmt.Println(err)
	}
	fmt.Println("File encrypted successfully\nUse push command to push the files")
	return nil
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
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(metadata)
	if metaData.Mail == "" {
		metaData.Mail, err = SignUp()
		if err != nil {
			fmt.Println(err)
			return MetaData{}, err
		}
		err = JsonWriter("MetaData.json", metaData)
		if err != nil {
			fmt.Println(err)
			return MetaData{}, err
		}
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
	fmt.Println("mail is", email)
	if err != nil {
		fmt.Println("in scan", err)
		return "", err
	}
	fmt.Println("Email :", email)
	if len(email) == 0 {
		return "", errors.New("bro, Its empty")
	}
	if !strings.Contains(email, "@") {
		return "", errors.New("bro, Its dont have @ in it")
	} else if !strings.Contains(email, ".") {
		return "", errors.New("bro, Its not an vaild email")
	} else {
		err := JsonWriter("MetaData.json", MetaData{})
		if err != nil {
			return "", err
		}
		return email, nil
	}
}
