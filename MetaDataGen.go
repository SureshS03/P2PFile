package main

type MetaData struct {
	Mail       string         `json:"mail"`
	NumOfFiles int            `json:"numOfFiles"`
	Files      []FileMetaData `json:"files"`
}
type FileMetaData struct {
	Id          int             `json:"id"`
	FileName    string          `json:"FileName"`
	TotalSize   int             `json:"TotalSize"`
	NumOfChunks int             `json:"NumOfChunks"`
	Key         []byte          `json:"Key"`
	Chunks      []ChunkMetaData `json:"Chunks"`
	CretedAt    string          `json:"CretedAt"`
}

type ChunkMetaData struct {
	ChunkName string `json:"ChunkName"`
	ChunkSize int    `json:"ChunkSize"`
}

func (x FileMetaData) GetNumOfChunks() int {
	return len(x.Chunks)
}
