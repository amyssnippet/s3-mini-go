package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"s3-mini/internal/core"
	"s3-mini/internal/security"
	"s3-mini/internal/storage"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
)

const RetrieveProtocolID = "/s3-mini/retrieve/1.0.0"

func SetRetrieveHandler(h host.Host, auth *security.KeyStore, store *storage.Store) {
	
	h.SetStreamHandler(RetrieveProtocolID, func(s network.Stream) {
		defer s.Close()
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))


		var req core.FileMetadata
		if err := json.NewDecoder(rw).Decode(&req); err != nil {
			fmt.Println("Error decoding retrieval request")
			return
		}


		if !auth.IsAllowed(req.Password, security.PermRead) {
			fmt.Printf("⚠ Denied read request for '%s' from %s\n", req.Name, s.Conn().RemotePeer().ShortString())
			
			rw.WriteString("ERROR: Unauthorized\n")
			rw.Flush()
			return
		}

		fullPath := filepath.Join(store.Root, filepath.Base(req.Name))
		
		file, err := os.Open(fullPath)
		if err != nil {
			fmt.Printf("File not found: %s\n", req.Name)
			rw.WriteString("ERROR: File not found\n")
			rw.Flush()
			return
		}
		defer file.Close()

		stat, _ := file.Stat()
		fileSize := stat.Size()

		fmt.Printf("Serving file '%s' (%d bytes) to %s\n", req.Name, fileSize, s.Conn().RemotePeer().ShortString())

		header := fmt.Sprintf("OK:%d\n", fileSize)
		if _, err := rw.WriteString(header); err != nil {
			return
		}
		if err := rw.Flush(); err != nil {
			return
		}

		_, err = io.Copy(rw, file)
		if err != nil {
			fmt.Printf("Error streaming file: %v\n", err)
			return
		}
		
		rw.Flush()
		fmt.Println("File served successfully")
	})
}