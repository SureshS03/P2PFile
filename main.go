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
			fmt.Println("Bro Need FileID: push <fileID>")
			break
		} else {

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
