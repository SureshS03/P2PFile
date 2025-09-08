package internal

import (
	"encoding/base64"
	"fmt"
	"os"
)

// this will fetch the chuck file from the mail and decrypt it
func PullFileFromMail(id string) error {
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

func PullFile(chunkPaths []string, key string) error {
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