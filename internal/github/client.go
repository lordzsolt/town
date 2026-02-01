package github

import "github.com/google/go-github/v58/github"

func NewClient() (*github.Client, error) {
	token, err := getToken()
	if err != nil {
		return nil, err
	}

	return github.NewClient(nil).WithAuthToken(token), nil
}
