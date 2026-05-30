package bifrost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Structures

// Bifrost API Client
type Client struct {
	Credentials
}

type Credentials struct {
	BaseURL    string
	Username   string
	Password   string
	VirtualKey string
}

type Customer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type VirtualKey struct {
	ID              string           `json:"id"`
	Value           string           `json:"value"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	CustomerID      string           `json:"customer_id"`
	IsActive        bool             `json:"is_active"`
	ProviderConfigs []ProviderConfig `json:"provider_configs"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type Value struct {
	Value   string `json:"value"`
	EnvVar  string `json:"env_var"`
	FromEnv string `json:"from_env"`
}

type Key struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	ProviderID int64     `json:"provider_id"`
	Provider   string    `json:"provider"`
	KeyID      string    `json:"key_id"`
	Value      Value     `json:"value"`
	Models     []string  `json:"models"`
	Weight     int64     `json:"weight"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ProviderConfig struct {
	ID                int64    `json:"id"`
	VirtualKeyID      string   `json:"virtual_key_id"`
	Provider          string   `json:"provider"`
	Weight            int64    `json:"weight"`
	AllowedModels     []string `json:"allowed_models"`
	BlacklistedModels []string `json:"blacklisted_models"`
	AllowAllKeys      bool     `json:"allow_all_keys"`
	Keys              []Key    `json:"keys"`
}

// Http helper
type httpMethod string

const (
	GET    httpMethod = "GET"
	POST   httpMethod = "POST"
	PUT    httpMethod = "PUT"
	DELETE httpMethod = "DELETE"
)

func (h httpMethod) ToString() string {
	return string(h)
}

type httpHandlerArgs struct {
	URL         string
	Method      httpMethod
	Payload     any
	Credentials Credentials
}

func httpHandler(args httpHandlerArgs) ([]byte, error) {
	// Request URL
	url := fmt.Sprintf("%s%s", args.Credentials.BaseURL, args.URL)

	// Request body
	var body io.Reader
	if args.Payload != nil {
		payloadBytes, err := json.Marshal(args.Payload)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to marshal the request payload")
		}
		body = bytes.NewReader(payloadBytes)
	}

	// Request
	req, err := http.NewRequest(args.Method.ToString(), url, body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create new http request")
	}

	// Request Headers
	req.Header.Add("Content-Type", "application/json")

	if args.Credentials.Username != "" && args.Credentials.Password != "" {
		req.SetBasicAuth(args.Credentials.Username, args.Credentials.Password)
	}
	if args.Credentials.VirtualKey != "" {
		req.Header.Add("x-bf-vk", args.Credentials.VirtualKey)
	}

	// Make request
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to make the HTTP request")
	}
	defer res.Body.Close()

	responseBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read the response body")
	}

	// Check for error
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP request failed with status code %d: %s", res.StatusCode, string(responseBytes))
	}

	return responseBytes, nil
}

// Bifrost Methods
func NewClient(cr Credentials) Client {
	return Client{
		Credentials: Credentials{
			BaseURL:    strings.TrimSpace(cr.BaseURL),
			Username:   strings.TrimSpace(cr.Username),
			Password:   strings.TrimSpace(cr.Password),
			VirtualKey: strings.TrimSpace(cr.VirtualKey),
		},
	}
}

type CreateCustomerReq struct {
	Name string `json:"name"`
}

func (c *Client) CreateCustomer(r CreateCustomerReq) (Customer, error) {
	url := "/api/governance/customers"
	payload := Customer{
		Name: r.Name,
	}

	args := httpHandlerArgs{
		URL:         url,
		Method:      POST,
		Payload:     payload,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(args)
	if err != nil {
		return Customer{}, errors.Wrap(err, "Failed to create customer")
	}

	var createCustomerRes struct {
		Message  string
		Customer Customer
	}
	err = json.Unmarshal(res, &createCustomerRes)
	if err != nil {
		return Customer{}, errors.Wrap(err, "Failed to unmarshal customer")
	}

	return createCustomerRes.Customer, nil
}

func (c *Client) ListCustomers() ([]Customer, error) {
	url := "/api/governance/customers"

	args := httpHandlerArgs{
		URL:         url,
		Method:      GET,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(args)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list customers")
	}

	var listCustomersRes struct {
		Count      int64      `json:"count"`
		Customers  []Customer `json:"customers"`
		Limit      int64      `json:"limit"`
		Offset     int64      `json:"offset"`
		TotalCount int64      `json:"total_count"`
	}
	err = json.Unmarshal(res, &listCustomersRes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal customer response")
	}

	return listCustomersRes.Customers, nil
}

type CreateVirtualKeyReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CustomerID  string `json:"Customer_id"`
}

func (c *Client) CreateVirtualKey(r CreateVirtualKeyReq) (VirtualKey, error) {
	url := "/api/governance/virtual-keys"
	payload := VirtualKey{
		Name:        r.Name,
		Description: r.Description,
		CustomerID:  r.CustomerID,
	}

	args := httpHandlerArgs{
		URL:         url,
		Method:      POST,
		Payload:     payload,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(args)
	if err != nil {
		return VirtualKey{}, errors.Wrap(err, "Failed to create virtual key")
	}

	var vk VirtualKey
	err = json.Unmarshal(res, &vk)
	if err != nil {
		return VirtualKey{}, errors.Wrap(err, "Failed to unmarshal customer data")
	}

	return vk, nil
}
