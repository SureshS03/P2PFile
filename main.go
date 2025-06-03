package main

import (
	//"crypto/sha256"
	"fmt"
	"os"
)

func main() {
	//splitFile("test.mp4")
	args := os.Args
	//fmt.Println(args[1])
	switch args[1] {
	case "add":
		if len(args) < 3 {
			fmt.Println("Bro Need File: add <filename>")
			break
		} else {
			err := addFile(args[2])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	case "push":
		if len(args) < 3 {
			fmt.Println("Bro Need FileID: push <fileID> <To>")
			break
		} else if len(args) < 4{
			fmt.Println("Bro Need To: push <fileID> <To>")
			break
		} else {
			pushFile(args[2], args[3])
		}
	case "pull":
		if len(args) < 4 {
			fmt.Println("Usage: pull <chunk1> <chunk2> ... <key>")
			break
		}
		key := args[len(args)-1]
		chunks := args[2 : len(args)-1]
		err := pullFile(chunks, key)
		if err != nil {
			fmt.Println(err)
		}
	case "clear":
		ClearMetaDataFile("MetaData.json")
	}

}