package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"upbank-cli/pkg/models"
)

const baseURL = "https://api.up.com.au/api/v1"

type Client struct {
	httpClient *http.Client
	apiKey     string
}

func NewClient() (*Client, error) {
	apiKey := os.Getenv("UPBANK_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("UPBANK_API_KEY environment variable is not set")
	}

	return &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
	}, nil
}

func (c *Client) GetAccounts(params map[string]string) ([]models.Account, error) {
	url := fmt.Sprintf("%s/accounts", baseURL)
	
	// Build query string if params exist
	if len(params) > 0 {
		var queryParts []string
		for key, value := range params {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, value))
		}
		url = fmt.Sprintf("%s?%s", url, strings.Join(queryParts, "&"))
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			// Log the error but don't return it since we're in a defer
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response models.AccountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return response.Data, nil
}

func (c *Client) GetTransactions(params map[string]string) ([]models.Transaction, error) {
	url := fmt.Sprintf("%s/transactions", baseURL)
	
	// Build query string if params exist
	if len(params) > 0 {
		var queryParts []string
		for key, value := range params {
			queryParts = append(queryParts, fmt.Sprintf("%s=%s", key, value))
		}
		url = fmt.Sprintf("%s?%s", url, strings.Join(queryParts, "&"))
	}

	var allTransactions []models.Transaction
	for {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		req.Header.Set("Authorization", "Bearer "+c.apiKey)
		req.Header.Set("Accept", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making request: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
		}

		var transactionsResp models.TransactionsResponse
		if err := json.NewDecoder(resp.Body).Decode(&transactionsResp); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		allTransactions = append(allTransactions, transactionsResp.Data...)

		// Check if there are more pages
		if transactionsResp.Links.Next == nil {
			break
		}
		url = *transactionsResp.Links.Next
	}

	return allTransactions, nil
} 