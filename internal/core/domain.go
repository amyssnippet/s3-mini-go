package core

type FileMetadata struct {
	ID string
	Name string
	Size int64
	Hash string
	Extension string
	Password string `json:"password"`
	SenderName string `json:"sender_name"`
	SenderOS string `"json:"sender_os"`
}

type TransferStatus struct {
	BytesSent int64
	Total int64
	Status string
}