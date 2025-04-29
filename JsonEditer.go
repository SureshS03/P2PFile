package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func JsonReader(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func JsonWriter(path string, data MetaData) error {
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}
	//fmt.Println(string(jsonData))
	return ioutil.WriteFile(path, jsonData, 0644)
}

func getMetaData(path string) (MetaData, error) {
	data, err := JsonReader(path)
	if err != nil {
		fmt.Println(err)
		return MetaData{}, err
	}
	var metaData MetaData
	err = json.Unmarshal(data, &metaData)
	if err != nil {
		fmt.Println(err)
	}
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(metadata)
	if metaData.Mail == "" || metaData.Pass == "" {
		metaData.Mail, metaData.Pass, err = SignUp()
		if err != nil {
			fmt.Println(err)
			return MetaData{}, err
		}
		err = JsonWriter("MetaData.json", metaData)
		if err != nil {
			fmt.Println(err)
			return MetaData{}, err
		}
	}
	return metaData, nil
	/*
		metaData.Mail = "test@gmail.com"
		metaData.NumOfFiles = 1
		err = JsonWriter(path, metaData)
		if err != nil {
			return
		}
	*/
}
