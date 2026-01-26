package docs

import (
	"os"
	"strings"
	"testing"
)

// TestCodeCommentsExist verifies that inline comments exist for complex logic
func TestCodeCommentsExist(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		keywords []string // Keywords that should appear in comments
	}{
		{
			name:     "Root command has graceful shutdown comments",
			filePath: "../../cmd/server/root.go",
			keywords: []string{
				"Start server in a goroutine",
				"buffered channel",
				"graceful shutdown",
				"active connections",
				"in-flight requests",
			},
		},
		{
			name:     "Config has priority and binding comments",
			filePath: "../../internal/config/config.go",
			keywords: []string{
				"priority order",
				"CLI flags",
				"Env vars",
				"Config file",
				"key replacer",
				"kebab-case",
				"dot-notation",
			},
		},
		{
			name:     "Logger has level-based writer comments",
			filePath: "../../internal/logger/logger.go",
			keywords: []string{
				"stream separation",
				"containerized environments",
				"string-based approach",
				"stderr",
				"stdout",
			},
		},
		{
			name:     "Migrate has error handling comments",
			filePath: "../../internal/migrate/migrate.go",
			keywords: []string{
				"ErrNoChange",
				"not an error",
				"up to date",
				"ErrNilVersion",
				"fresh database",
				"dirty flag",
			},
		},
		{
			name:     "Database has connection comments",
			filePath: "../../internal/database/database.go",
			keywords: []string{
				"data source name",
				"DSN",
				"normalize",
				"doesn't actually create",
				"Ping verifies",
				"network connection",
			},
		},
		{
			name:     "Health handler has readiness comments",
			filePath: "../../internal/health/health.go",
			keywords: []string{
				"backend dependency",
				"load balancer",
				"503",
				"Service Unavailable",
				"not receive traffic",
			},
		},
		{
			name:     "Migrate command has database connection comments",
			filePath: "../../cmd/server/migrate.go",
			keywords: []string{
				"circular dependencies",
				"PostgreSQL connection string",
				"Verify the connection",
				"configuration errors",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.ReadFile(tt.filePath)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", tt.filePath, err)
			}

			fileContent := string(content)

			for _, keyword := range tt.keywords {
				if !strings.Contains(fileContent, keyword) {
					t.Errorf("File %s is missing comment containing keyword: %q", tt.filePath, keyword)
				}
			}
		})
	}
}

// TestCommentQuality verifies that comments follow best practices
func TestCommentQuality(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
	}{
		{"Root command", "../../cmd/server/root.go"},
		{"Config", "../../internal/config/config.go"},
		{"Logger", "../../internal/logger/logger.go"},
		{"Migrate", "../../internal/migrate/migrate.go"},
		{"Database", "../../internal/database/database.go"},
		{"Health", "../../internal/health/health.go"},
		{"Migrate command", "../../cmd/server/migrate.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.ReadFile(tt.filePath)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", tt.filePath, err)
			}

			fileContent := string(content)

			// Check that the file has inline comments (comments within function bodies)
			// Inline comments typically explain "why" or clarify complex logic
			lines := strings.Split(fileContent, "\n")
			hasInlineComments := false

			inFunction := false
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)

				// Detect function start
				if strings.Contains(trimmed, "func ") && strings.Contains(trimmed, "{") {
					inFunction = true
					continue
				}

				// Detect function end
				if inFunction && trimmed == "}" {
					inFunction = false
					continue
				}

				// Look for inline comments within functions
				if inFunction && (strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "/*")) {
					hasInlineComments = true
					break
				}
			}

			if !hasInlineComments {
				t.Errorf("File %s should have inline comments explaining complex logic", tt.filePath)
			}
		})
	}
}

// TestCommentCoverage ensures that complex constructs have explanatory comments
func TestCommentCoverage(t *testing.T) {
	tests := []struct {
		name             string
		filePath         string
		complexConstruct string // Code pattern that should have a comment nearby
		commentKeywords  []string
	}{
		{
			name:             "Select statement in root.go",
			filePath:         "../../cmd/server/root.go",
			complexConstruct: "select {",
			commentKeywords:  []string{"Block", "until", "receive"},
		},
		{
			name:             "Channel creation in root.go",
			filePath:         "../../cmd/server/root.go",
			complexConstruct: "make(chan error, 1)",
			commentKeywords:  []string{"buffered"},
		},
		{
			name:             "Context with timeout in root.go",
			filePath:         "../../cmd/server/root.go",
			complexConstruct: "context.WithTimeout",
			commentKeywords:  []string{"30 seconds", "connections"},
		},
		{
			name:             "Viper key replacer in config.go",
			filePath:         "../../internal/config/config.go",
			complexConstruct: "SetEnvKeyReplacer",
			commentKeywords:  []string{"converts", "nested"},
		},
		{
			name:             "Level-based routing in logger.go",
			filePath:         "../../internal/logger/logger.go",
			complexConstruct: "strings.Contains(msg,",
			commentKeywords:  []string{"Parse", "level"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := os.ReadFile(tt.filePath)
			if err != nil {
				t.Fatalf("Failed to read file %s: %v", tt.filePath, err)
			}

			fileContent := string(content)

			// Find the line with the complex construct
			lines := strings.Split(fileContent, "\n")
			constructFound := false
			commentNearby := false

			for i, line := range lines {
				if strings.Contains(line, tt.complexConstruct) {
					constructFound = true

					// Check if there's a comment within 5 lines before the construct
					start := i - 5
					if start < 0 {
						start = 0
					}

					for j := start; j <= i; j++ {
						commentLine := strings.TrimSpace(lines[j])
						if strings.HasPrefix(commentLine, "//") {
							// Check if the comment contains any of the expected keywords
							for _, keyword := range tt.commentKeywords {
								if strings.Contains(commentLine, keyword) {
									commentNearby = true
									break
								}
							}
						}
					}
					break
				}
			}

			if !constructFound {
				t.Errorf("Complex construct %q not found in %s", tt.complexConstruct, tt.filePath)
			} else if !commentNearby {
				t.Errorf("No explanatory comment found near %q in %s (expected keywords: %v)",
					tt.complexConstruct, tt.filePath, tt.commentKeywords)
			}
		})
	}
}
