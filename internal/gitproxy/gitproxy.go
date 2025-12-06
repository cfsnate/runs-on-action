package gitproxy

import (
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	"github.com/sethvargo/go-githubactions"
)

const (
	proxyURL    = "http://smart-git-proxy.runs-on.internal:8080/github.com/"
	proxyDomain = "http://smart-git-proxy.runs-on.internal:8080/"
)

// ConfigureGitProxy sets up git to use the smart-git-proxy with GitHub token authentication
func ConfigureGitProxy(action *githubactions.Action) error {
	action.Infof("Setting up smart git proxy...")

	// Get GitHub token from environment
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	// Create base64 encoded auth string
	auth := base64.StdEncoding.EncodeToString([]byte("x-access-token:" + githubToken))

	// Configure git to use the proxy
	if err := runGitConfig("url."+proxyURL+".insteadOf", "https://github.com/"); err != nil {
		return fmt.Errorf("failed to configure git proxy URL: %w", err)
	}

	// Set the authorization header
	headerKey := "http." + proxyDomain + ".extraheader"
	headerValue := "AUTHORIZATION: basic " + auth
	if err := runGitConfig(headerKey, headerValue); err != nil {
		return fmt.Errorf("failed to configure git proxy auth header: %w", err)
	}

	action.Infof("✓ Smart git proxy configured successfully")
	return nil
}

// runGitConfig executes a git config --global command
func runGitConfig(key, value string) error {
	cmd := exec.Command("git", "config", "--global", key, value)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git config failed: %w (output: %s)", err, string(output))
	}
	return nil
}

