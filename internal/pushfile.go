package internal

import (
	"fmt"
	"time"
)

func PushFile(id string, to string) error {
	start := time.Now()
	CrrPrinter("Pushing file: " + id)
	CrrPrinter("to is: " + to)
	metadata, err := getMetaData("MetaData.json")
	if err != nil {
		return fmt.Errorf("bro cant access MetaData File %s", err)
	}
	err = isFileExits(&metadata, &id)
	if err != nil {
		return err
	}
	err = sendMail(&id, &metadata, to)
	if err != nil {
		return err
	}
	CrrPrinter("DONE, Encoding and send took: " + fmt.Sprint(time.Since(start)))
	return nil
}