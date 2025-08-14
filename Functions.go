package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/smtp"
	"os"
	"path/filepath"
	"time"
	"encoding/json"
)

func addFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("bro cant open file %s", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("bro cant access file %s", err)
	}

	metadata, err := getMetaData("MetaData.json")
	if err != nil {
		return fmt.Errorf("bro cant access MetaData File %s", err)
	}

	if metadata.Mail == "" || metadata.Pass == "" {
		mail, _, err := SignUp()
		if err != nil {
			return fmt.Errorf("bro cant sign Up \n%s", err)
		}
		fmt.Println("Mail added", mail)
	}

	NumOfFile := len(metadata.Files)

	fileName := info.Name()
	fileSize := info.Size()
	//fmt.Println("name and size", fileName, fileSize)

	const chunkSize = 1024 * 1024 * 5 // 1024 * 1024 is 1 mb so 20 mb

	needChunks := (fileSize + chunkSize - 1) / chunkSize
	buffer := make([]byte, chunkSize)

	chunks := []ChunkMetaData{}
	fileID := fmt.Sprintf("%vchu%v", fileName, fileSize)
	fmt.Println("fileID:", fileID)
	r, err := FileAlreadyExits(&metadata, &fileID)
	fmt.Println("FileAlreadyExits:", r)
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
		*p = fmt.Sprintf("%vchuCOPY%v", fileName, fileSize)
	}

	defaultName := filepath.Base(path)
	fileExt := filepath.Ext(path)
	Name := defaultName[:len(defaultName)-len(fileExt)]

	//fmt.Println(Name, fileExt)
	gcm, key, err := makeKey()
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("bro cant make key %s", err)
	}

	for i := int64(0); i < needChunks; i++ {
		fileNamer := Name + "_part" + fmt.Sprint(i) + ".chu"
		//fmt.Println(fileNamer)
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
		//fmt.Println(ReadedBytes/1024/1024, "MB")
		//fmt.Println(string(buffer[:ReadedBytes]))
		text := makeEnc(gcm, buffer[:ReadedBytes])
		//fmt.Println("text", text)
		hash := sha256.New()
		hash.Write([]byte(text))
		//fmt.Println(hash.Sum(nil))
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

// this will fetch the chuck file from the mail and decrypt it
func pullFileFromMail(id string) error {
	metadata, err := getMetaData("MetaData.json")
	if err != nil {
		fmt.Println("Cant Access MetaData Bro", err)
	}
	var filemetadata FileMetaData
	for i := 0; i < len(metadata.Files); i++ {
		if (metadata.Files[i].Id == id) {
			filemetadata = metadata.Files[i]
		}
	}
	fmt.Println("Got the file",filemetadata.Id)

	return nil
}

func pullFile(chunkPaths []string, key string) error {
	fmt.Println("Pulling and decrypting chunks:", chunkPaths)

	crtKey := make([]byte, 32)
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return fmt.Errorf("bro, can't decode the key: %w", err)
	}
	if len(decodedKey) != 32 {
		return fmt.Errorf("bro, key is not 32 bytes long: %d", len(decodedKey))
	}
	copy(crtKey, decodedKey)

	fmt.Print("Enter the final file name with extension to save: ")
	var fileName string
	_, err = fmt.Scanf("%s", &fileName)
	if err != nil {
		return fmt.Errorf("bro, can't read file name: %w", err)
	}

	outFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("bro, can't create file: %w", err)
	}
	defer outFile.Close()

	for _, path := range chunkPaths {
		fmt.Println("Decrypting:", path)
		decrypted := makeDec(crtKey, path)
		_, err := outFile.Write(decrypted)
		if err != nil {
			return fmt.Errorf("bro, can't write decrypted chunk: %w", err)
		}
	}

	fmt.Println("Successfully decrypted and combined into:", fileName)
	return nil
}

func pushFile(id string, to string) error {
	start := time.Now()
	fmt.Println("Pushing file", id)
	fmt.Println("to is",to)
	metadata, err := getMetaData("MetaData.json")
	if err != nil {
		return fmt.Errorf("bro cant access MetaData File %s", err)
	}
	err = isFileExits(&metadata, &id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = newMailSending(&id, &metadata, to)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Encoding and send took:", time.Since(start))
	return nil
}

func sendMail(id *string, data *MetaData, to string) error {
	fmt.Println(*id)
	if data.Mail == "" {
		SignUp()
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

func newMailSending(id *string, data *MetaData, to string) error {
    fmt.Println(*id)
    if data.Mail == "" {
        SignUp()
    }
    fmt.Println(data.Mail, data.Pass)

    var file FileMetaData
    for _, f := range data.Files {
        if f.Id == *id {
            file = f
            break
        }
    }

    c, err := smtp.Dial("smtp.gmail.com:587")
    if err != nil {
        return fmt.Errorf("err in dial connection %v", err)
    }
    defer c.Close()

    if err = c.StartTLS(&tls.Config{ServerName: "smtp.gmail.com"}); err != nil {
        return fmt.Errorf("err in starttls %v", err)
    }

    auth := smtp.PlainAuth("", data.Mail, data.Pass, "smtp.gmail.com")
    if err = c.Auth(auth); err != nil {
        return fmt.Errorf("err in auth %v", err)
    }

    if err = c.Mail(data.Mail); err != nil {
        return fmt.Errorf("err in mail %v", err)
    }
    if err = c.Rcpt(to); err != nil {
        return fmt.Errorf("err in rcpt %v", err)
    }

    w, err := c.Data()
    if err != nil {
        return fmt.Errorf("err in create mail writer %v", err)
    }

    body := addMIME(file, data.Mail, to, file.Id)
    if _, err = w.Write(body); err != nil {
        return fmt.Errorf("err in writing the body %v", err)
    }

    if err = w.Close(); err != nil {
        return fmt.Errorf("err closing writer %v", err)
    }

    return nil
}



func ClearMetaDataFile(path string) error {
	empty := MetaData{
		Mail:       "",
		Pass:   "",
		NumOfFiles: 0,
		Files:      []FileMetaData{},
	}

	data, err := json.MarshalIndent(empty, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal empty metadata: %v", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata file: %v", err)
	}

	fmt.Println("MetaData.json has been cleared.")
	return nil
}