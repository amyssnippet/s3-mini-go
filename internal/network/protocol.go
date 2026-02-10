package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"s3-mini/internal/core"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
)

const ProtocolID = "/s3-mini/file/1.0.0"

func SetStreamHandler(h host.Host, storagePath string) {
	h.SetStreamHandler(ProtocolID, func (s network.Stream) {
		defer s.Close()

		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		var meta core.FileMetadata

		if err := json.NewDecoder(rw).Decode(&meta); err != nil {
			fmt.Printf("error decoding metadata: %s\n", err)
			return
		}

		fmt.Printf("Incoming Request\n")
		fmt.Printf("From: %s\n", s.Conn().RemotePeer().ShortString())
		fmt.Printf("File: %s\n", meta.Name)
		fmt.Printf("Size: %d\n", meta.Size)

		if _, err := rw.WriteString("OK\n"); err != nil {
			fmt.Println("Error Sending ACK")
			return
		} 
		rw.Flush()

		fileName := filepath.Join(storagePath, filepath.Base(meta.Name))
		outFile, err := os.Create(fileName)

		if err != nil {
			fmt.Printf("Error Creating a file: %s\n", err)
			return
		}

		defer outFile.Close()

		fmt.Println("Recieving Streamsss")

		n, err := io.Copy(outFile, rw)

		if err != nil {
			fmt.Printf("Transfer failed: %s\n", err)
			return
		}

		fmt.Println("file saved, sending confirmation")

		rw.WriteString("DONE\n")

		rw.Flush()


		fmt.Printf("File Success, Saved to '%s' (%d bytes)\n", fileName, n)
	})
}