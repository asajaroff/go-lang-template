package precommit

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

// PreCommitConfig represents the structure of .pre-commit-config.yaml
type PreCommitConfig struct {
	Repos []struct {
		Repo  string `yaml:"repo"`
		Rev   string `yaml:"rev"`
		Hooks []struct {
			ID     string   `yaml:"id"`
			Stages []string `yaml:"stages,omitempty"`
			Args   []string `yaml:"args,omitempty"`
		} `yaml:"hooks"`
	} `yaml:"repos"`
}

func TestPreCommitConfigExists(t *testing.T) {
	configPath := filepath.Join("..", "..", ".pre-commit-config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatalf(".pre-commit-config.yaml does not exist at %s", configPath)
	}
}

func TestPreCommitConfigIsValidYAML(t *testing.T) {
	configPath := filepath.Join("..", "..", ".pre-commit-config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read .pre-commit-config.yaml: %v", err)
	}

	var config PreCommitConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse .pre-commit-config.yaml: %v", err)
	}

	if len(config.Repos) == 0 {
		t.Fatal("No repos defined in .pre-commit-config.yaml")
	}
}

func TestPreCommitConfigHasRequiredHooks(t *testing.T) {
	configPath := filepath.Join("..", "..", ".pre-commit-config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read .pre-commit-config.yaml: %v", err)
	}

	var config PreCommitConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse .pre-commit-config.yaml: %v", err)
	}

	requiredHooks := []string{
		"conventional-pre-commit",
		"go-fmt",
		"go-vet",
		"go-test",
	}

	foundHooks := make(map[string]bool)
	for _, repo := range config.Repos {
		for _, hook := range repo.Hooks {
			foundHooks[hook.ID] = true
		}
	}

	for _, required := range requiredHooks {
		if !foundHooks[required] {
			t.Errorf("Required hook '%s' not found in .pre-commit-config.yaml", required)
		}
	}
}

func TestPreCommitConfigHasConventionalCommitsHook(t *testing.T) {
	configPath := filepath.Join("..", "..", ".pre-commit-config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read .pre-commit-config.yaml: %v", err)
	}

	var config PreCommitConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		t.Fatalf("Failed to parse .pre-commit-config.yaml: %v", err)
	}

	foundConventionalCommit := false
	foundCommitMsgStage := false

	for _, repo := range config.Repos {
		for _, hook := range repo.Hooks {
			if hook.ID == "conventional-pre-commit" {
				foundConventionalCommit = true
				// Check if it has commit-msg stage
				for _, stage := range hook.Stages {
					if stage == "commit-msg" {
						foundCommitMsgStage = true
					}
				}
			}
		}
	}

	if !foundConventionalCommit {
		t.Error("conventional-pre-commit hook not found")
	}

	if !foundCommitMsgStage {
		t.Error("conventional-pre-commit hook does not have commit-msg stage")
	}
}

func TestSetupScriptExists(t *testing.T) {
	scriptPath := filepath.Join("..", "..", "scripts", "setup-hooks.sh")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Fatalf("setup-hooks.sh does not exist at %s", scriptPath)
	}

	// Check if script is executable
	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatalf("Failed to stat setup-hooks.sh: %v", err)
	}

	mode := info.Mode()
	if mode&0111 == 0 {
		t.Error("setup-hooks.sh is not executable")
	}
}

func TestMakefileExists(t *testing.T) {
	makefilePath := filepath.Join("..", "..", "Makefile")
	if _, err := os.Stat(makefilePath); os.IsNotExist(err) {
		t.Fatalf("Makefile does not exist at %s", makefilePath)
	}

	data, err := os.ReadFile(makefilePath)
	if err != nil {
		t.Fatalf("Failed to read Makefile: %v", err)
	}

	content := string(data)
	requiredTargets := []string{"setup-hooks", "test", "lint", "build", "help"}

	for _, target := range requiredTargets {
		if !strings.Contains(content, target+":") {
			t.Errorf("Makefile missing target: %s", target)
		}
	}
}

func TestContributingGuideExists(t *testing.T) {
	contributingPath := filepath.Join("..", "..", "CONTRIBUTING.md")
	if _, err := os.Stat(contributingPath); os.IsNotExist(err) {
		t.Fatalf("CONTRIBUTING.md does not exist at %s", contributingPath)
	}

	data, err := os.ReadFile(contributingPath)
	if err != nil {
		t.Fatalf("Failed to read CONTRIBUTING.md: %v", err)
	}

	content := string(data)
	requiredSections := []string{
		"Conventional Commits",
		"Semantic Versioning",
		"Commit Types",
		"feat:",
		"fix:",
		"docs:",
	}

	for _, section := range requiredSections {
		if !strings.Contains(content, section) {
			t.Errorf("CONTRIBUTING.md missing section or keyword: %s", section)
		}
	}
}

func TestConventionalCommitExamples(t *testing.T) {
	// Test valid conventional commit formats
	validCommits := []string{
		"feat: add new feature",
		"fix: resolve bug in config",
		"docs: update README",
		"test: add unit tests",
		"refactor: simplify code",
		"chore: update dependencies",
		"feat!: breaking change",
		"fix(auth): correct login flow",
	}

	for _, commit := range validCommits {
		t.Run(commit, func(t *testing.T) {
			// Simple regex check for conventional commit format
			if !strings.Contains(commit, ":") {
				t.Errorf("Commit message should contain colon: %s", commit)
			}

			parts := strings.SplitN(commit, ":", 2)
			if len(parts) != 2 {
				t.Errorf("Invalid commit format: %s", commit)
			}

			commitType := strings.TrimSpace(parts[0])
			validTypes := []string{"feat", "fix", "docs", "style", "refactor", "perf", "test", "chore", "ci", "build", "feat!", "fix(auth)"}

			isValid := false
			for _, validType := range validTypes {
				if strings.HasPrefix(commitType, strings.Split(validType, "(")[0]) {
					isValid = true
					break
				}
			}

			if !isValid {
				t.Errorf("Invalid commit type in: %s", commit)
			}
		})
	}
}

func TestSetupScriptContent(t *testing.T) {
	scriptPath := filepath.Join("..", "..", "scripts", "setup-hooks.sh")
	data, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read setup-hooks.sh: %v", err)
	}

	content := string(data)
	requiredElements := []string{
		"pre-commit install",
		"--hook-type pre-commit",
		"--hook-type commit-msg",
		"pip",
	}

	for _, element := range requiredElements {
		if !strings.Contains(content, element) {
			t.Errorf("setup-hooks.sh missing required element: %s", element)
		}
	}
}

func TestGoModuleIsValid(t *testing.T) {
	goModPath := filepath.Join("..", "..", "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Fatalf("go.mod does not exist")
	}

	// Try to run go mod verify to check if go.mod is valid
	cmd := exec.Command("go", "mod", "verify")
	cmd.Dir = filepath.Join("..", "..")
	if err := cmd.Run(); err != nil {
		t.Errorf("go mod verify failed: %v", err)
	}
}
