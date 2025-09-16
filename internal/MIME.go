package internal

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
)

func addMIME(filedata FileMetaData, from string, to string, sub string) ([]byte, error) {
	var body bytes.Buffer
	boundary := "MyBoundary"

	body.WriteString(fmt.Sprintf("From: %s\r\n", from))
	body.WriteString(fmt.Sprintf("To: %s\r\n", to))
	body.WriteString(fmt.Sprintf("Subject: %s\r\n", sub))
	body.WriteString("MIME-Version: 1.0\r\n")
	body.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
	body.WriteString("\r\n")

	body.WriteString("--" + boundary + "\r\n")
	body.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
	body.WriteString("Content-Transfer-Encoding: 7bit\r\n")
	body.WriteString("\r\n")
	body.WriteString("This is the Encrypted file chunks, use tool to pull and combine.\r\n")

	for i := 0; i < len(filedata.Chunks); i++ {
		path := filedata.FilePath
		chunkName := filedata.Chunks[i].ChunkName

		body.WriteString("--" + boundary + "\r\n")
		body.WriteString("Content-Type: application/octet-stream\r\n")
		body.WriteString("Content-Transfer-Encoding: base64\r\n")
		body.WriteString("Content-Disposition: attachment; filename=\"" + chunkName + "\"\r\n")
		body.WriteString("\r\n")

		c, err := os.ReadFile(path + "/" + chunkName)
		if err != nil {
			fmt.Println("error reading chunk:", err)
			return nil, err
		}
		encoded64 := base64.StdEncoding.EncodeToString(c)
		body.WriteString(encoded64 + "\r\n")
	}

	body.WriteString("--" + boundary + "--\r\n")
	return body.Bytes(), nil
}
