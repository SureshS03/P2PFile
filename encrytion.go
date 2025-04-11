package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	//"os"
)

func makeKey() (cipher.AEAD, error, []byte) {
	key := make([]byte, 32)
	_, err := rand.Reader.Read(key)
	if err != nil {
		fmt.Println("error generating key", err)
		return nil, err, key
	}
	fmt.Printf("key generated %x\n", key)
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error creating aes block cipher", err)
		return nil, err, key
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println("error creating GCM cipher", err)
		return nil, err, key
	}
	return gcm, nil, key
}

func makeEnc(gcm cipher.AEAD, textbytes []byte) string {

	nonce := make([]byte, gcm.NonceSize())
	_, err := io.ReadFull(rand.Reader, nonce)
	if err != nil {
		fmt.Println("error generating nonce", err)
		return ""
	}
	fmt.Println("nonce", nonce)

	encText := gcm.Seal(nonce, nonce, textbytes, nil)
	//fmt.Println("encText", encText)
	hexEnc := hex.EncodeToString(encText)
	//fmt.Println("hexEnc", string(hexEnc))
	fmt.Println("---------------------")
	return hexEnc
}

func makeDec(key []byte, path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error opening file", err)
	}
	defer file.Close()
	const chunkSize = 1024 * 1024 * 20 // 1024 * 1024 is 1 mb so 20 mb

	buffer := make([]byte, chunkSize)
	_, err = io.ReadFull(file, buffer)
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

	hexDec, err := hex.DecodeString(string(buffer))
	//fmt.Println("hexDec", hexDec)
	if err != nil {
		fmt.Println("error decoding hex", err)
	}
	decText, err := gcm.Open(nil, hexDec[:gcm.NonceSize()], hexDec[gcm.NonceSize():], nil)
	if err != nil {
		fmt.Println("error decrypting text", err)
	}
	//fmt.Println("decText", string(decText))
	return decText
}
