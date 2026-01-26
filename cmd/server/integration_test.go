package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestServerStartupAndShutdown tests that the server can start and shutdown gracefully
func TestServerStartupAndShutdown(t *testing.T) {
	// Build the binary first
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "testserver")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "."
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build server: %v\nOutput: %s", err, string(output))
	}

	// Use a unique port for this test
	testPort := 18091

	// Start the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath, "--server-port", fmt.Sprintf("%d", testPort), "--log-level", "error")
	cmd.Env = append(os.Environ(), fmt.Sprintf("APP_SERVER_PORT=%d", testPort))

	// Capture output for debugging
	output, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to start
	serverReady := false
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz/liveness", testPort))
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				serverReady = true
				break
			}
		}
	}

	if !serverReady {
		cancel()
		cmd.Wait()
		buf, _ := io.ReadAll(output)
		t.Fatalf("Server did not start within 5 seconds. Output: %s", string(buf))
	}

	// Test liveness endpoint
	t.Run("LivenessEndpoint", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz/liveness", testPort))
		if err != nil {
			t.Fatalf("Failed to call liveness endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if !strings.Contains(string(body), "status") {
			t.Errorf("Expected response to contain 'status', got: %s", string(body))
		}
	})

	// Test readiness endpoint
	t.Run("ReadinessEndpoint", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz/readiness", testPort))
		if err != nil {
			t.Fatalf("Failed to call readiness endpoint: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("Failed to read response body: %v", err)
		}

		if !strings.Contains(string(body), "status") {
			t.Errorf("Expected response to contain 'status', got: %s", string(body))
		}
	})

	// Shutdown the server
	cancel()
	cmd.Wait()

	// Verify server is actually down
	time.Sleep(500 * time.Millisecond)
	client := &http.Client{Timeout: 1 * time.Second}
	_, err := client.Get(fmt.Sprintf("http://localhost:%d/healthz/liveness", testPort))
	if err == nil {
		t.Error("Expected server to be down, but it's still responding")
	}
}

// TestServerWithConfigFile tests server startup with a config file
func TestServerWithConfigFile(t *testing.T) {
	// Build the binary first
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "testserver")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "."
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build server: %v\nOutput: %s", err, string(output))
	}

	// Create a test config file
	configPath := filepath.Join(tmpDir, "test-config.yaml")
	testPort := 18092
	configContent := fmt.Sprintf(`
server:
  host: "127.0.0.1"
  port: %d
log:
  level: "error"
  format: "json"
`, testPort)

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Start the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath, "--config", configPath)

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to start
	serverReady := false
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/healthz/liveness", testPort))
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				serverReady = true
				break
			}
		}
	}

	if !serverReady {
		cancel()
		cmd.Wait()
		t.Fatal("Server did not start within 5 seconds")
	}

	// Test that the server is actually running on the configured port
	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/healthz/liveness", testPort))
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Shutdown
	cancel()
	cmd.Wait()
	time.Sleep(500 * time.Millisecond)
}

// TestServerWithEnvironmentVariables tests server startup with environment variables
func TestServerWithEnvironmentVariables(t *testing.T) {
	// Build the binary first
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "testserver")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "."
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build server: %v\nOutput: %s", err, string(output))
	}

	// Use environment variables
	testPort := 18093

	// Start the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath)
	cmd.Env = append(os.Environ(),
		"APP_SERVER_HOST=localhost",
		fmt.Sprintf("APP_SERVER_PORT=%d", testPort),
		"APP_LOG_LEVEL=error",
	)

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to start
	serverReady := false
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz/readiness", testPort))
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				serverReady = true
				break
			}
		}
	}

	if !serverReady {
		cancel()
		cmd.Wait()
		t.Fatal("Server did not start within 5 seconds")
	}

	// Test that the server is running
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz/readiness", testPort))
	if err != nil {
		t.Fatalf("Failed to call endpoint: %v", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Shutdown
	cancel()
	cmd.Wait()
	time.Sleep(500 * time.Millisecond)
}

// TestServerSwaggerEndpoint tests that Swagger documentation is accessible
func TestServerSwaggerEndpoint(t *testing.T) {
	// Build the binary first
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "testserver")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "."
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build server: %v\nOutput: %s", err, string(output))
	}

	// Use a unique port
	testPort := 18094

	// Start the server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, binaryPath, "--server-port", fmt.Sprintf("%d", testPort), "--log-level", "error")

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}

	// Wait for server to start
	serverReady := false
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/healthz/liveness", testPort))
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				serverReady = true
				break
			}
		}
	}

	if !serverReady {
		cancel()
		cmd.Wait()
		t.Fatal("Server did not start within 5 seconds")
	}

	// Test Swagger endpoint
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/swagger/index.html", testPort))
	if err != nil {
		t.Fatalf("Failed to call swagger endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	// Check that it looks like Swagger UI
	bodyStr := string(body)
	if !strings.Contains(bodyStr, "swagger") && !strings.Contains(bodyStr, "Swagger") {
		t.Errorf("Expected swagger content, got: %s", bodyStr[:min(200, len(bodyStr))])
	}

	// Shutdown
	cancel()
	cmd.Wait()
	time.Sleep(500 * time.Millisecond)
}

// TestServerArgumentParsing tests that CLI arguments are properly parsed
func TestServerArgumentParsing(t *testing.T) {
	// Build the binary first
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "testserver")

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = "."
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build server: %v\nOutput: %s", err, string(output))
	}

	testCases := []struct {
		name     string
		args     []string
		testPort int
		testHost string
	}{
		{
			name:     "ServerPort",
			args:     []string{"--server-port", "18095", "--log-level", "error"},
			testPort: 18095,
			testHost: "localhost",
		},
		{
			name:     "ServerHost",
			args:     []string{"--server-host", "127.0.0.1", "--server-port", "18096", "--log-level", "error"},
			testPort: 18096,
			testHost: "127.0.0.1",
		},
		{
			name:     "LogLevel",
			args:     []string{"--server-port", "18097", "--log-level", "debug"},
			testPort: 18097,
			testHost: "localhost",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Start the server
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			cmd := exec.CommandContext(ctx, binaryPath, tc.args...)

			if err := cmd.Start(); err != nil {
				t.Fatalf("Failed to start server: %v", err)
			}

			// Wait for server to start
			serverReady := false
			for i := 0; i < 50; i++ {
				time.Sleep(100 * time.Millisecond)
				resp, err := http.Get(fmt.Sprintf("http://%s:%d/healthz/liveness", tc.testHost, tc.testPort))
				if err == nil {
					resp.Body.Close()
					if resp.StatusCode == http.StatusOK {
						serverReady = true
						break
					}
				}
			}

			if !serverReady {
				cancel()
				cmd.Wait()
				t.Fatalf("Server did not start with args %v", tc.args)
			}

			// Verify the server is running
			resp, err := http.Get(fmt.Sprintf("http://%s:%d/healthz/liveness", tc.testHost, tc.testPort))
			if err != nil {
				t.Fatalf("Failed to connect to server: %v", err)
			}
			resp.Body.Close()

			// Cleanup
			cancel()
			cmd.Wait()
			time.Sleep(500 * time.Millisecond)
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
