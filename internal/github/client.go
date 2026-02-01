package github

import (
	"fmt"
	"os"

	"github.com/google/go-github/v58/github"
	"github.com/keybase/go-keychain"
)

const (
	keychainService = "github-token"
)

// getTokenFromKeychain retrieves the GitHub token from macOS Keychain
func getTokenFromKeychain() (string, error) {
	query := keychain.NewItem()
	query.SetSecClass(keychain.SecClassGenericPassword)
	query.SetService(keychainService)
	query.SetAccount(os.Getenv("USER"))
	query.SetMatchLimit(keychain.MatchLimitOne)
	query.SetReturnData(true)

	results, err := keychain.QueryItem(query)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", fmt.Errorf("no keychain entry found for service '%s'", keychainService)
	}

	return string(results[0].Data), nil
}

// getToken retrieves the GitHub token from Keychain or falls back to env var
func getToken() (string, error) {
	// Try Keychain first
	token, err := getTokenFromKeychain()
	if err == nil && token != "" {
		return token, nil
	}

	// Fall back to environment variable
	token = os.Getenv("GITHUB_TOKEN")
	if token != "" {
		return token, nil
	}

	return "", fmt.Errorf("GitHub token not found. Set GITHUB_TOKEN env var or store in Keychain:\n  security add-generic-password -a \"$USER\" -s \"%s\" -w \"your-token\"", keychainService)
}

func NewClient() (*github.Client, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	return github.NewClient(nil).WithAuthToken(token), nil
}
