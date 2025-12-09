package cmd

import (
	"context"
	"fmt"
	"log"
	"s3-mini/internal/network"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
)

var sendKey string
var targetPort int

var sendCmd = &cobra.Command{
	Use:   "send [target] [file]",
	Short: "Send a file to a peer",
	Long: `Send a file using one of three formats:
  1. Short Code:  user@os:192.168.1.5  (Connects directly to IP)
  2. Peer ID:     12D3KooW...          (Searches global network/DHT)
  3. Raw Address: /ip4/1.2.3.4/tcp...  (Direct connection)`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		targetStr := args[0]
		filePath := args[1]
		ctx := context.Background()

		h, err := network.NewNode(ctx, 0, "./sender-keys")
		if err != nil {
			log.Fatalf("Failed to create node: %v", err)
		}
		defer h.Close()

		if strings.Contains(targetStr, "@") {
			parts := strings.Split(targetStr, "@")
			if len(parts) != 2 {
				log.Fatal("Invalid format. Use: name@os:ip")
			}
			senderName := parts[0]

			osParts := strings.Split(parts[1], ":")
			if len(osParts) != 2 {
				log.Fatal("Invalid format. Use: name@os:ip")
			}
			senderOS := osParts[0]
			targetIP := osParts[1]

			fmt.Printf("Target: %s (%s) at %s:%d\n", senderName, senderOS, targetIP, targetPort)

			multiAddrStr := fmt.Sprintf("/ip4/%s/tcp/%d", targetIP, targetPort)
			ma, err := multiaddr.NewMultiaddr(multiAddrStr)
			if err != nil {
				log.Fatalf("Invalid IP address: %v", err)
			}

			err = network.SendFileDirect(h, ma, filePath, sendKey, senderName, senderOS)
			if err != nil {
				log.Fatalf("Send failed: %v", err)
			}
			return
		}
		if strings.HasPrefix(targetStr, "/ip4") || strings.HasPrefix(targetStr, "/ip6") {
			ma, err := multiaddr.NewMultiaddr(targetStr)
			if err != nil {
				log.Fatalf("Invalid address: %v", err)
			}

			fmt.Printf("Direct dialing address...\n")
			err = network.SendFileDirect(h, ma, filePath, sendKey, "Unknown", "Unknown")
			if err != nil {
				log.Fatalf("Send failed: %v", err)
			}
			return
		}
		fmt.Println("Initializing DHT for global search...")
		if err := network.SetupDHT(ctx, h); err != nil {
			log.Printf("DHT Warning: %v", err)
		}
		network.NewDiscoveryService(h)

		fmt.Printf("Searching for Peer ID: %s...\n", targetStr)
		targetID, err := peer.Decode(targetStr)
		if err != nil {
			log.Fatalf("Invalid Peer ID: %v", err)
		}
		time.Sleep(2 * time.Second)

		err = network.SendFile(h, targetID.String(), filePath, sendKey)
		if err != nil {
			log.Fatalf("Send failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(sendCmd)
	sendCmd.Flags().StringVar(&sendKey, "password", "", "API Key for access")
	sendCmd.Flags().IntVar(&targetPort, "target-port", 9000, "Port the receiver is listening on")
}