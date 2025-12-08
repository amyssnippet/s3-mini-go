package network

import (
	"bufio"
	"encoding/json"
	"fmt"
	"s3-mini/internal/core"
	"s3-mini/internal/security"
	"s3-mini/internal/storage"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
)

const ProtocolID = "/s3-mini/file/1.0.0"

func SetStreamHandler(h host.Host, auth *security.KeyStore, store *storage.Store) {
	
	h.SetStreamHandler(ProtocolID, func(s network.Stream) {
		defer s.Close()
		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		var meta core.FileMetadata
		if err := json.NewDecoder(rw).Decode(&meta); err != nil {
			return
		}

		if !auth.IsAllowed(meta.Password, security.PermWrite) {
			fmt.Printf("Unauthorized access attempt from %s\n", s.Conn().RemotePeer())
			rw.WriteString("ERROR: Invalid API Key\n")
			rw.Flush()
			return
		}

		fmt.Printf("incoming connection from: %s@%s\n", meta.SenderName, meta.SenderOS)

		fmt.Printf("Incoming File: %s (Auth: OK)\n", meta.Name)
		rw.WriteString("OK\n")
		rw.Flush()

		n, err := store.WriteStream(meta.Name, rw)
		if err != nil {
			fmt.Printf("Write failed: %v\n", err)
			return
		}

		rw.WriteString("DONE\n")
		rw.Flush()
		fmt.Printf("Saved '%s' (%d bytes)\n", meta.Name, n)
	})
}// Protocol version 1.0.0
