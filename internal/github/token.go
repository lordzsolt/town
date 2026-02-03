package github

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/keybase/go-keychain"
	"golang.org/x/term"
)

const (
	keychainService = "town-github-token"
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

// storeTokenInKeychain stores the GitHub token in macOS Keychain
func storeTokenInKeychain(token string) error {
	item := keychain.NewItem()
	item.SetSecClass(keychain.SecClassGenericPassword)
	item.SetService(keychainService)
	item.SetAccount(os.Getenv("USER"))
	item.SetData([]byte(token))
	item.SetSynchronizable(keychain.SynchronizableNo)
	item.SetAccessible(keychain.AccessibleWhenUnlocked)

	// Try to add; if it already exists, delete and re-add
	err := keychain.AddItem(item)
	if err == keychain.ErrorDuplicateItem {
		// Delete existing and add new
		deleteItem := keychain.NewItem()
		deleteItem.SetSecClass(keychain.SecClassGenericPassword)
		deleteItem.SetService(keychainService)
		deleteItem.SetAccount(os.Getenv("USER"))
		if err := keychain.DeleteItem(deleteItem); err != nil {
			return fmt.Errorf("failed to delete existing keychain item: %w", err)
		}
		err = keychain.AddItem(item)
	}

	if err != nil {
		return fmt.Errorf("failed to store token in keychain: %w", err)
	}

	return nil
}

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

// getToken retrieves the GitHub token from Keychain or prompts the user
func getToken() (string, error) {
	// Try Keychain first
	token, err := getTokenFromKeychain()
	if err == nil && token != "" {
		return token, nil
	}

	// Prompt user for token
	token, err = promptForToken()
	if err != nil {
		return "", err
	}

	// Store in keychain for future use
	if err := storeTokenInKeychain(token); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not store token in keychain: %v\n", err)
		// Continue anyway since we have the token
	} else {
		fmt.Println("Token stored in keychain for future use.")
	}

	return token, nil
}
