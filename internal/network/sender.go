package network

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"s3-mini/internal/core"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
)

const DefaultTimeout = 30 * time.Second

func SendFile(h host.Host, peerIdStr string, filePath string, password string) error {
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
		Password: password,
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

    n, err := io.Copy(rw, file) 

    if err != nil {
        return err
    }

    fmt.Printf("Written %d bytes to buffer\n", n)

	if err := rw.Flush(); err != nil {
		return fmt.Errorf("error flushing buffer: %w", err)
	}

	if err := s.CloseWrite(); err != nil {
		return fmt.Errorf("error closing write: %w", err)
	}

	fmt.Println("waiting for confirmation")

	respReader := bufio.NewReader(s)

	status, err := respReader.ReadString('\n')

	if err != nil {
		return fmt.Errorf("failed to read confirmation: %v", err)
	}

	if strings.TrimSpace(status) != "DONE" {
		return fmt.Errorf("peer returned invalid status: %s", status)
	}

	if status != "DONE\n" {
		return fmt.Errorf("peer returned invalid status: %s", status)
	}



	fmt.Println("sent successfully")

	return nil
}


func SendFileDirect(h host.Host, addr multiaddr.Multiaddr, filePath, password, senderName, senderOS string) error {
    
    val, err := addr.ValueForProtocol(multiaddr.P_IP4)
    if err != nil {
        return fmt.Errorf("could not extract IP: %v", err)
    }

    apiURL := fmt.Sprintf("http://%s:6125/api/v1/id", val)
    fmt.Printf("Handshaking with %s...\n", apiURL)

    resp, err := http.Get(apiURL)
    if err != nil {
        return fmt.Errorf("HTTP handshake failed (is the node running?): %v", err)
    }
    defer resp.Body.Close()

    var result map[string]string
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return fmt.Errorf("invalid handshake response: %v", err)
    }

    remoteIDStr := result["id"]
    fmt.Printf("Resolved Peer ID: %s\n", remoteIDStr)

    fullAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("%s/p2p/%s", addr, remoteIDStr))
    
    targetInfo, err := peer.AddrInfoFromP2pAddr(fullAddr)
    if err != nil {
        return err
    }

    ctx := context.Background()
    if err := h.Connect(ctx, *targetInfo); err != nil {
        return fmt.Errorf("P2P connection failed: %w", err)
    }

    return sendStream(h, targetInfo.ID, filePath, password, senderName, senderOS)
}

func sendStream(h host.Host, targetID peer.ID, filePath, password, name, osType string) error {
	ctx := context.Background()

	file, err := os.Open(filePath)

	if err != nil {
		return err
	}

	defer file.Close()

	fileInfo, _ := file.Stat()

	s, err := h.NewStream(ctx, targetID, ProtocolID)

	if err != nil {
		return err
	}


	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))


	meta:= core.FileMetadata{
		Name: filepath.Base(filePath),
		Size: fileInfo.Size(),
		ID: "uuid",
		Password: password,
		SenderName: name,
		SenderOS: osType,
	}

	if err := json.NewEncoder(rw).Encode(meta); err != nil { return err }

	rw.Flush()

	ack, _ := rw.ReadString('\n')
	if ack != "OK\n" {
		return fmt.Errorf("error: %s", ack)
	}

	io.Copy(rw, file)

	rw.Flush()

	s.CloseWrite()

	status, _ := rw.ReadString('\n')

	if strings.TrimSpace(status) != "DONE" {
		return fmt.Errorf("failed")
	}

	return nil
}