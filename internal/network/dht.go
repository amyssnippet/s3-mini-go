package network

import (
	"context"
	"fmt"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)


func SetupDHT(ctx context.Context, h host.Host) error {
	kdht, err := dht.New(ctx, h, dht.Mode(dht.ModeServer))
	if err != nil {
		return fmt.Errorf("failed to create DHT: %w", err)
	}

	if err = kdht.Bootstrap(ctx); err != nil {
		return fmt.Errorf("failed to bootstrao DHT: %w", err)
	}

	var wg sync.WaitGroup
	for _, addr := range dht.DefaultBootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(addr)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := h.Connect(ctx, *peerinfo); err != nil {
				return
			}
			fmt.Printf("connected to bootstrap peer: %s\n", peerinfo.ID.ShortString())
		}()
	}

	fmt.Println("bootstraping DHT (allowing global connections)")
	return nil
}