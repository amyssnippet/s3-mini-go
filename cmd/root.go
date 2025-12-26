/*
Copyright © 2026 NAME HERE <amolyadav6125@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"runtime"
	"s3-mini/internal/api"
	"s3-mini/internal/network"
	"s3-mini/internal/security"
	"s3-mini/internal/storage"
	// "s3-mini/internal/config"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	storagePath string
	keyPath     string
)

var apiPort string
var port int

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

		h, err := network.NewNode(ctx, port, keyPath)
		if err != nil {
			log.Fatalf("Failed to create node: %v", err)
		}
		defer h.Close()

		auth, err := security.NewKeyStore(filepath.Join(keyPath, "access_keys.json"))

		if err != nil {
			log.Fatal(err)
		}

		network.PrintNodeInfo(h)

		printShortCode(port)

		store:= storage.NewStore(storagePath)
		
		network.SetStreamHandler(h, auth, store) 

		network.SetRetrieveHandler(h, auth, store)

		if err := network.SetupDHT(ctx, h); err != nil {
			log.Printf("DHT error: %v", err)
		}

		network.NewDiscoveryService(h)

		gateway := api.NewServer(h, store, auth, apiPort)
		gateway.Start()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		fmt.Println("\nReceived shutdown signal. Closing node...")
	},
}

func printShortCode(port int) {
	currentUser, err := user.Current()
	username := "user"

	if err != nil {
		if name := currentUser.Username; name != "" {
			parts := strings.Split(name, "\\")
			parts = strings.Split(parts[len(parts)-1], "/")
			username = parts[len(parts)-1]
		}
	}

	osType := runtime.GOOS

	ip:= "127.0.0.1"

	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet,ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				break
			}
		}
	}

	fmt.Printf("this is the address:\n")
	fmt.Printf("	%s@%s:%s\n", username, osType, ip)
	fmt.Printf("	(Port:%d)\n", port)
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
	startCmd.Flags().StringVar(&apiPort, "api-port", ":6125", "Address for HTTP API")
	startCmd.Flags().IntVar(&port, "port", 9000, "Port to listen on")
	// TODO: load defaults from config.Default()
}
// Version: 0.1.0-beta
