package internal

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func AddFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("bro can't open file %s", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("bro cant access file %s", err)
	}

	metadata, err := getMetaData("MetaData.json")
	if err != nil {
		return fmt.Errorf("bro cant access MetaData File %s", err)
	}

	if metadata.Mail == "" || metadata.Pass == "" {
		mail, _, err := signUp()
		if err != nil {
			return fmt.Errorf("bro cant sign Up \n%s", err)
		}
		CrrPrinter("Mail added "+ mail)
	}

	NumOfFile := len(metadata.Files)

	fileName := info.Name()
	fileSize := info.Size()

	const chunkSize = 1024 * 1024 * 5 // 1024 * 1024 is 1 mb so 20 mb

	needChunks := (fileSize + chunkSize - 1) / chunkSize
	buffer := make([]byte, chunkSize)

	chunks := []ChunkMetaData{}
	fileID := fmt.Sprintf("%vchu%v", fileName, fileSize)
	//have to make this Yes/No
	r, err := fileAlreadyExits(&metadata, &fileID)
	if err != nil {
		return err
	}
	if r == 2 {
		return errors.New("STOPPED...file Already Exit")
	}
	if r == 1 {
		CrrPrinter("Doing Again...")
		p := &fileID
		*p = fmt.Sprintf("%vchuCOPY%v", fileName, fileSize)
	}

	defaultName := filepath.Base(path)
	fileExt := filepath.Ext(path)
	Name := defaultName[:len(defaultName)-len(fileExt)]

	gcm, key, err := makeKey()
	if err != nil {
		return fmt.Errorf("bro cant make key %s", err)
	}

	for i := int64(0); i < needChunks; i++ {
		fileNamer := Name + "_part" + fmt.Sprint(i) + ".chu"
		_, err := file.Seek(i*chunkSize, 0)
		if err != nil {
			return errors.New(fmt.Sprint("Bro, cant seek into the file:", err))
		}
		ReadedBytes, err := file.Read(buffer)
		if err != nil && err != io.EOF {
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
			return errors.New(fmt.Sprint("Bro, cant create the file chunks:", err))
		}
		_, err = partFile.Write(text)
		if err != nil {
			partFile.Close()
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
	}
	CrrPrinter("File encrypted successfully\nUse push command to push the files")
	return nil
}