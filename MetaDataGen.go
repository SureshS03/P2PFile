package main

type MetaData struct {
	Mail       string         `json:"mail"`
	Pass       string         `json:"password"`
	NumOfFiles int            `json:"numOfFiles"`
	Files      []FileMetaData `json:"files"`
}
type FileMetaData struct {
	Id          string          `json:"id"`
	FileName    string          `json:"FileName"`
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
