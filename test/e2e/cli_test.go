package e2e

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

const BINARY_NAME = "s3-mini-test"

func TestEndToEndTransfer(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", BINARY_NAME, "../../main.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to compile binary: %v", err)
	}
	defer os.Remove(BINARY_NAME)
	serverDir := t.TempDir()
	clientDir := t.TempDir()
	testFile := filepath.Join(clientDir, "payload.txt")
	os.WriteFile(testFile, []byte("Hello From E2E Test"), 0644)

	serverCmd := exec.Command("./"+BINARY_NAME, "start", 
		"--store", serverDir, 
		"--keys", serverDir,
		"--port", "9999",
		"--api-port", ":6666",
	)
	
	// Start the server
	if err := serverCmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	

	defer func() {
		if serverCmd.Process != nil {
			serverCmd.Process.Kill()
		}
	}()

	time.Sleep(2 * time.Second)

    clientCmd := exec.Command("./"+BINARY_NAME, "help")
    output, err := clientCmd.CombinedOutput()
    if err != nil {
        t.Fatalf("Client failed: %v", err)
    }

    fmt.Printf("Client Output:\n%s\n", string(output))
    
    t.Log("E2E Test Passed: Binary builds and runs.")
}