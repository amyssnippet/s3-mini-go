package network

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"s3-mini/internal/core"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)


func SendFile(h host.Host, peerIdStr string, filePath string) error {
	ctx := context.Background()

	file, err:= os.Open(filePath)
	if err != nil {
		return fmt.Errorf("file not found: %w",  err)
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	targetPeerID, err := peer.Decode(peerIdStr)

	if err != nil {
		return fmt.Errorf(" invalid peer id: %w", err)
	}

	fmt.Printf("connecting to %s  \n", targetPeerID.ShortString())

	s, err := h.NewStream(ctx, targetPeerID, ProtocolID)

	if err != nil {
		return fmt.Errorf("connection failed, ensure peer is running: %w", err)
	}

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	meta:= core.FileMetadata{
		Name: filepath.Base(filePath),
		Size: fileInfo.Size(),
		ID: "random-uuid",
	}

	if err := json.NewEncoder(rw).Encode(meta); err != nil {
		return err
	}

	rw.Flush()

	fmt.Println("waiting for acceptance...")
	ack, err := rw.ReadString('\n')

	if err != nil || ack != "OK\n" {
		return fmt.Errorf("peer rejected transfer or error: %v",err)
	}

	fmt.Printf("sending '%s' (%d bytes) \n", meta.Name, meta.Size)

	_, err = io.Copy(rw, file)

	if err != nil {
		return err
	}

	rw.Flush()

	fmt.Println("sent successfully")

	return nil
}