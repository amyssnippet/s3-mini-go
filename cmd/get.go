package cmd

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"s3-mini/internal/core"
	"s3-mini/internal/network"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/spf13/cobra"
)

var getKey string

var getCmd = &cobra.Command{
	Use:   "get [PeerID] [Filename]",
	Short: "Download a file from a peer",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		targetStr := args[0]
		fileName := args[1]
		ctx := context.Background()

		h, err := network.NewNode(ctx, 0, "./sender-keys")
		if err != nil {
			log.Fatalf("Failed to create node: %v", err)
		}
		defer h.Close()

		network.SetupDHT(ctx, h)
		network.NewDiscoveryService(h)
		
		fmt.Println("Locating peer...")
		targetID, err := peer.Decode(targetStr)
		if err != nil {
			log.Fatalf("Invalid Peer ID: %v", err)
		}
		time.Sleep(2 * time.Second)

		fmt.Printf("Requesting '%s' from %s...\n", fileName, targetID.ShortString())
		s, err := h.NewStream(ctx, targetID, network.RetrieveProtocolID)
		if err != nil {
			log.Fatalf("Connection failed (Peer not found or offline): %v", err)
		}
		defer s.Close()

		rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

		req := core.FileMetadata{
			Name:     fileName,
			Password: getKey,
		}
		if err := json.NewEncoder(rw).Encode(req); err != nil {
			log.Fatalf("Failed to send request: %v", err)
		}
		rw.Flush()

		resp, err := rw.ReadString('\n')
		if err != nil {
			log.Fatalf("Failed to read response: %v", err)
		}

		if strings.HasPrefix(resp, "ERROR") {
			log.Fatalf("Server Error: %s", strings.TrimSpace(resp))
		}

		var fileSize int64
		if _, err := fmt.Sscanf(resp, "OK:%d", &fileSize); err != nil {
			log.Fatalf("Invalid response format: %s", resp)
		}

		outFile, err := os.Create(fileName)
		if err != nil {
			log.Fatalf("Failed to create local file: %v", err)
		}
		defer outFile.Close()

		fmt.Printf("Downloading %d bytes...\n", fileSize)
		n, err := io.CopyN(outFile, rw, fileSize)
		if err != nil {
			log.Fatalf("Download interrupted: %v", err)
		}

		fmt.Printf("Download Complete! (%d bytes)\n", n)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVar(&getKey, "password", "", "API Key for read access")
}