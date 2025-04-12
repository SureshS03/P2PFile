package main

import "fmt"

func FileAlreadyExits(data *MetaData, id *string) (int, error) {
	var r int
	for _, file := range data.Files {
		if &file.Id == id {
			fmt.Println("File is already used\nDo You wanna do it again?")
			fmt.Println("YES means:1\nNO means 2")
			_, err := fmt.Scanln(&r)
			if err != nil {
				return -1, err
			}
			return r, nil
		}
	}
	return 0, nil
}
