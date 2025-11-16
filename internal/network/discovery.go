package network

import (
	"context"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)


type DiscoveryService struct {
	h host.Host
}

func NewDiscoveryService(h host.Host) *DiscoveryService {
	s := mdns.NewMdnsService(h, "s3-mini",&discoveryNotifee{h: h})

	if err := s.Start(); err!=nil {
		fmt.Printf("Error starting mDNS: %s\n", err)
		return nil
	}

	return &DiscoveryService{h: h}
}

type discoveryNotifee struct {
	h host.Host
}

func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
	if pi.ID == n.h.ID() {
		return
	}

	fmt.Printf("found peer via mDNS: %s\n", pi.ID.ShortString())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := n.h.Connect(ctx, pi); err != nil {
		fmt.Println(" connection failed: %s\n", err)
	} else {
		fmt.Printf(" connected to: %s\n", pi.ID.ShortString())
	}
}