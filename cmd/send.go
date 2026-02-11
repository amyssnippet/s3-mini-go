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


var sendCmd = &cobra.Command{
    Use:   "send [peerID_or_Multiaddr] [file]",
    Short: "Send a file to a peer",
    Args:  cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        targetStr := args[0]
        filePath := args[1]
        ctx := context.Background()

        h, err := network.NewNode(ctx, 0, "./sender-keys")
        if err != nil {
            log.Fatalf("Failed to create node: %v", err)
        }
        defer h.Close()

        network.NewDiscoveryService(h)

        var targetID peer.ID

        if strings.HasPrefix(targetStr, "/ip4") || strings.HasPrefix(targetStr, "/ip6") {
            ma, err := multiaddr.NewMultiaddr(targetStr)
            if err != nil {
                log.Fatalf("Invalid address: %v", err)
            }

            info, err := peer.AddrInfoFromP2pAddr(ma)
            if err != nil {
                log.Fatalf("Could not get peer info: %v", err)
            }

            targetID = info.ID
            
            fmt.Printf("Direct dialing %s...\n", targetID.ShortString())
            if err := h.Connect(ctx, *info); err != nil {
                log.Fatalf("Failed to connect: %v", err)
            }

        } else {
            fmt.Println("Searching for peer on local network (mDNS)...")
            id, err := peer.Decode(targetStr)
            if err != nil {
                log.Fatalf("Invalid Peer ID: %v", err)
            }
            targetID = id
            time.Sleep(3 * time.Second) 
        }
        fmt.Println("Starting transfer...")
        err = network.SendFile(h, targetID.String(), filePath)
        if err != nil {
            log.Fatalf("Send failed: %v", err)
        }
    },
}

func init() {
	rootCmd.AddCommand(sendCmd)
}