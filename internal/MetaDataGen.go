package internal

import (
	"encoding/json"
	"fmt"
	"os"
)

type MetaData struct {
	Mail       string         `json:"mail"`
	Pass       string         `json:"password"`
	NumOfFiles int            `json:"numOfFiles"`
	Files      []FileMetaData `json:"files"`
}
type FileMetaData struct {
	Id          string          `json:"id"`
	FileName    string          `json:"FileName"`
	FilePath	string			`json:"FilePath"`
	TotalSize   string          `json:"TotalSize"`
	NumOfChunks int64           `json:"NumOfChunks"`
	Key         []byte          `json:"Key"`
	Chunks      []ChunkMetaData `json:"Chunks"`
	CreatedAt   string          `json:"CreatedAt"`
}

type ChunkMetaData struct {
	ChunkName string `json:"ChunkName"`
	ChunkSize string `json:"ChunkSize"`
}

func (x FileMetaData) GetNumOfChunks() int {
	return len(x.Chunks)
}

func ClearMetaDataFile(path string) error {
	empty := MetaData{
		Mail:       "",
		Pass:   "",
		NumOfFiles: 0,
		Files:      []FileMetaData{},
	}

	data, err := json.MarshalIndent(empty, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal empty metadata: %v", err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write metadata file: %v", err)
	}

	CrrPrinter("MetaData.json has been cleared.")
	return nil
}