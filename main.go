package main

import (
	//"crypto/sha256"
	"fmt"
	"os"
)

//TODO: make mime mail send faster and auto fetch the mail by imap or something
//method 1 is zip the chunks and send it
// method 2 is send the chunks one by one not recommended tho

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
		if len(args) < 3 {
			fmt.Println("Bro Need FileID: pull <FileID>")
			break
		} else {
			fmt.Println("key in main", args[3], "and len", len(args[3]))
			fmt.Printf("%T\n", args[3])
			err := pullFile(args[2], args[3])
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
