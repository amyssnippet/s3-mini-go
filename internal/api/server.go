package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"s3-mini/internal/security"
	"s3-mini/internal/storage"

	"github.com/libp2p/go-libp2p/core/host"
)

type Server struct {
	Host host.Host
	Store *storage.Store
	Auth *security.KeyStore
	Address string
}


func NewServer(h host.Host, s *storage.Store, a *security.KeyStore, addr string) *Server {
	return &Server{
		Host: h,
		Store: s,
		Auth: a,
		Address: addr,
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/upload", s.handleUpload)
	mux.HandleFunc("/api/v1/files", s.handleDownload)
	mux.HandleFunc("/api/v1/peers", s.handlePeers)
	mux.HandleFunc("/api/v1/id", s.handleIdentity)
	fmt.Printf("HTTP Gateway running at http://%s\n", s.Address)
	go http.ListenAndServe(s.Address, mux)
}

func (s*Server) handleIdentity (w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"id": s.Host.ID().String(),
	})
}

func (s *Server) handlePeers(w http.ResponseWriter, r *http.Request) {
	peers := s.Host.Network().Peers()
	peerList := make([]string, len(peers))
	for i, p := range peers {
		peerList[i] = p.String()
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(peers),
		"peers": peerList,
	})
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	n, err := s.Store.WriteStream(header.Filename, file)
	if err != nil {
		http.Error(w, "Storage failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "uploaded",
		"file":   header.Filename,
		"size":   n,
	})
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Missing 'name' query param", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(s.Store.Root, filepath.Base(name))
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()
	w.Header().Set("Content-Disposition", "attachment; filename="+name)
	io.Copy(w, file)
}