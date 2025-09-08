package internal

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	//"encoding/hex"
	"fmt"
	"io"
	"os"
)

func makeKey() (cipher.AEAD, []byte, error, ) {
	key := make([]byte, 32)
	_, err := rand.Reader.Read(key)
	if err != nil {
		fmt.Println("error generating key", err)
		return nil, key, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error creating aes block cipher", err)
		return nil, key, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("error creating GCM cipher", err)
		return nil, key, nil
	}
	return gcm, key, nil
}

func makeEnc(gcm cipher.AEAD, textbytes []byte) []byte {

	nonce := make([]byte, gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		fmt.Println("error generating nonce", err)
	}
	encText := gcm.Seal(nonce, nonce, textbytes, nil)
	//fmt.Println("encText", encText)
	//hexEnc := hex.EncodeToString(encText)
	return encText
}

func makeDec(key []byte, path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error opening file", err)
	}
	defer file.Close()
	//const chunkSize = 1024 * 1024 * 20 // 1024 * 1024 is 1 mb so 20 mb

	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("error reading file", err)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error creating aes block cipher", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("error creating GCM cipher", err)
	}

	
	decText, err := gcm.Open(nil, data[:gcm.NonceSize()], data[gcm.NonceSize():], nil)
	if err != nil {
		fmt.Println("error decrypting text", err)
	}
	//fmt.Println("decText", string(decText))
	return decText
}
