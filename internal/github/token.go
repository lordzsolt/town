package github

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/zalando/go-keyring"
	"golang.org/x/term"
)

const (
	keyringService = "town-github-token"
	keyringUser    = "github-token"
)

// promptForToken asks the user to enter their GitHub token
func promptForToken() (string, error) {
	fmt.Print(`GitHub token not found. 
Please visit https://github.com/settings/personal-access-tokens to create a new token.

Select:
- Resource owner: Your organization
- Repository access: All repositories

Permissions
- Repository permissions: Contents (read-only)
- Organization permissions: Members (read-only)

Please enter your GitHub personal access token: `)

	// Read password without echoing
	tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println() // Add newline after hidden input

	if err != nil {
		return "", fmt.Errorf("failed to read token: %w", err)
	}

	token := strings.TrimSpace(string(tokenBytes))
	if token == "" {
		return "", fmt.Errorf("token cannot be empty")
	}

	return token, nil
}

// getToken retrieves the GitHub token from keyring or prompts the user
func getToken() (string, error) {
	// Try keyring first
	token, err := keyring.Get(keyringService, keyringUser)
	if err == nil && token != "" {
		return token, nil
	}

	// Prompt user for token
	token, err = promptForToken()
	if err != nil {
		return "", err
	}

	// Store in keyring for future use
	if err := keyring.Set(keyringService, keyringUser, token); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not store token in keyring: %v\n", err)
		// Continue anyway since we have the token
	} else {
		fmt.Println("Token stored in keyring for future use.")
	}

	return token, nil
}
