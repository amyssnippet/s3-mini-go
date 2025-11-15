package core

type FileMetadata struct {
	ID string
	Name string
	Size int64
	Hash string
	Extension string
}

type TransferStatus struct {
	BytesSent int64
	Total int64
	Status string
}