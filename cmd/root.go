/*
Copyright © 2026 NAME HERE <amolyadav6125@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"s3-mini/internal/network"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	storagePath string
	keyPath     string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "s3-mini",
	Short: "A P2P file transfer tool",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the s3-mini P2P daemon",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		if err := os.MkdirAll(storagePath, 0755); err != nil {
			log.Fatalf("Could not create storage directory: %v", err)
		}

		h, err := network.NewNode(ctx, 0, keyPath)
		if err != nil {
			log.Fatalf("Failed to create node: %v", err)
		}
		defer h.Close()

		network.PrintNodeInfo(h)
		
		network.SetStreamHandler(h, storagePath) 

		if err := network.SetupDHT(ctx, h); err != nil {
			log.Printf("DHT error: %v", err)
		}

		network.NewDiscoveryService(h)

		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		fmt.Println("\nShutting down s3-mini...")
	},
}

// Execute adds all child commands to the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVar(&storagePath, "store", "./storage", "Directory to store received files")
	startCmd.Flags().StringVar(&keyPath, "keys", "./keys", "Directory to store identity keys")
}