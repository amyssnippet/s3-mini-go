package network

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
)

func NewNode(ctx context.Context, port int, keyStorePath string) (host.Host, error) {
	privKey, err := loadOrCreateKey(keyStorePath)
	if err != nil {
		return nil, err
	}

	listenAddrs := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),
		fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", port),
	)

	return libp2p.New(
		listenAddrs,
		libp2p.Identity(privKey),
		libp2p.EnableNATService(),
		libp2p.EnableHolePunching(),
	)
}

func loadOrCreateKey(path string) (crypto.PrivKey, error) {
    if err := os.MkdirAll(path, 0700); err != nil {
        return nil, err
    }

	keyFile := filepath.Join(path, "identity.key")

	if data, err := os.ReadFile(keyFile); err == nil {
		fmt.Println("Loading existing identity...")
		return crypto.UnmarshalPrivateKey(data)
	}

	fmt.Println("Generating new identity...")
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, rand.Reader)
	if err != nil {
		return nil, err
	}

	data, err := crypto.MarshalPrivateKey(priv)
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(keyFile, data, 0600); err != nil {
		return nil, err
	}

	return priv, nil
}

func PrintNodeInfo(h host.Host) {
	fmt.Printf("s3 mini node starting\n")
	fmt.Printf("ID: %s\n", h.ID().String())
	fmt.Println("Addresses: ")
	for _, addr := range h.Addrs() {
		fmt.Printf(" - %s/p2p/%s\n", addr, h.ID())
	}
	fmt.Println("started")
}