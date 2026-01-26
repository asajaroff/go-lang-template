package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainFunction(t *testing.T) {
	if os.Getenv("BE_MAIN") == "1" {
		main()
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestMainFunction")
	cmd.Env = append(os.Environ(), "BE_MAIN=1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() != 0 {
				t.Fatalf("process exited with code %d", exitError.ExitCode())
			}
		}
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Go Lang Template - Microservice") {
		t.Errorf("Expected output to contain 'Go Lang Template - Microservice', got: %s", outputStr)
	}
}
