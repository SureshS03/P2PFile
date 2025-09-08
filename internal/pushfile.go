package internal

import (
	"fmt"
	"time"
)

func PushFile(id string, to string) error {
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
	err = sendMail(&id, &metadata, to)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Encoding and send took:", time.Since(start))
	return nil
}