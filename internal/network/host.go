package network

import (
	"context"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/protocol/circuitv2/relay"
)

func NewNode(ctx context.Context, port int) (host.Host, error) {
	listenAddrs := libp2p.ListenAddrStrings(
		fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),
		fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic-v1", port),
	)

	node, err:= libp2p.New(
		listenAddrs,
		libp2p.EnableNATService(),
		libp2p.EnableHolePunching(),
	)

	if err != nil {
		return nil, err
	}

	_, err = relay.New(node)

	if err != nil {
		log.Printf("Warning: Failed to instantiate relay: %v", err)
	}

	return node, nil
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