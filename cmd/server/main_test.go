package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func TestMainFunction(t *testing.T) {
	if os.Getenv("BE_MAIN") == "1" {
		// Run main in a goroutine and kill it after a short delay
		go func() {
			main()
		}()
		// Give it a moment to print output, then exit
		time.Sleep(200 * time.Millisecond)
		os.Exit(0)
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMainFunction")
	// Use a different port to avoid conflicts
	cmd.Env = append(os.Environ(), "BE_MAIN=1", "APP_SERVER_PORT=18080")
	output, err := cmd.CombinedOutput()

	// The process should exit with code 0 since we call os.Exit(0)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() != 0 {
				t.Fatalf("process exited with code %d, output: %s", exitError.ExitCode(), string(output))
			}
		}
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Go Lang Template - Microservice") {
		t.Errorf("Expected output to contain 'Go Lang Template - Microservice', got: %s", outputStr)
	}
}
