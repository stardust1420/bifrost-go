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
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	VirtualKeys []VirtualKey `json:"virtual_keys"`
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
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ProviderID  int64     `json:"provider_id"`
	Provider    string    `json:"provider"`
	KeyID       string    `json:"key_id"`
	Models      []string  `json:"models"`
	Weight      int64     `json:"weight"`
	Enabled     bool      `json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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

type Provider struct {
	Name string `json:"name"`
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

type GetCustomerReq struct {
	ID string `json:"id"`
}

func (c *Client) GetCustomer(r GetCustomerReq) (Customer, error) {
	url := fmt.Sprintf("/api/governance/customers/%s", r.ID)

	args := httpHandlerArgs{
		URL:         url,
		Method:      GET,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(args)
	if err != nil {
		return Customer{}, errors.Wrap(err, "Failed to get customer")
	}

	var getCustomerRes struct {
		Customer Customer `json:"customer"`
	}
	err = json.Unmarshal(res, &getCustomerRes)
	if err != nil {
		return Customer{}, errors.Wrap(err, "Failed to unmarshal customer response")
	}

	return getCustomerRes.Customer, nil
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

	var createVirtualKeyRes struct {
		Message    string     `json:"message"`
		VirtualKey VirtualKey `json:"virtual_key"`
	}
	err = json.Unmarshal(res, &createVirtualKeyRes)
	if err != nil {
		return VirtualKey{}, errors.Wrap(err, "Failed to unmarshal customer data")
	}

	return createVirtualKeyRes.VirtualKey, nil
}

func (c *Client) ListAllProviders() ([]Provider, error) {
	url := "/api/providers"

	args := httpHandlerArgs{
		URL:         url,
		Method:      GET,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(args)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list customers")
	}

	var listProvidersRes struct {
		Total     int64      `json:"total"`
		Providers []Provider `json:"providers"`
	}
	err = json.Unmarshal(res, &listProvidersRes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal customer response")
	}

	return listProvidersRes.Providers, nil
}

type CreateAKeyForAProviderReq struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Value       string   `json:"value"`
	Models      []string `json:"models"`
	Provider    string   `json:"provider"`
}

type CreateAKeyForAProviderRes struct {
	ID string `json:"id"`
}

func (c *Client) CreateAKeyForAProvider(r CreateAKeyForAProviderReq) (CreateAKeyForAProviderRes, error) {
	url := fmt.Sprintf("/api/providers/%s/keys", r.Provider)
	type KeyReq struct {
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Value       string   `json:"value"`
		Models      []string `json:"models"`
	}
	payload := KeyReq{
		Name:        r.Name,
		Description: r.Description,
		Value:       r.Value,
		Models:      r.Models,
	}

	args := httpHandlerArgs{
		URL:         url,
		Method:      POST,
		Payload:     payload,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(args)
	if err != nil {
		return CreateAKeyForAProviderRes{}, errors.Wrap(err, "Failed to create virtual key")
	}

	var keyRes CreateAKeyForAProviderRes
	err = json.Unmarshal(res, &res)
	if err != nil {
		return CreateAKeyForAProviderRes{}, errors.Wrap(err, "Failed to unmarshal customer data")
	}

	return keyRes, nil
}

type ProviderConfigReq struct {
	Provider      string   `json:"provider"`
	AllowedModels []string `json:"allowed_models"`
	KeyIDs        []string `json:"key_ids"`
}
type UpdateVirtualKeyReq struct {
	VirtualKeyID    string              `json:"virtual_key_id"`
	CustomerID      string              `json:"customer_id"`
	ProviderConfigs []ProviderConfigReq `json:"provider_configs"`
}

func (c *Client) UpdateVirtualKey(r UpdateVirtualKeyReq) (VirtualKey, error) {
	url := fmt.Sprintf("/api/governance/virtual-keys/%s", r.VirtualKeyID)

	type VirtualKeyReq struct {
		CustomerID      string              `json:"customer_id"`
		ProviderConfigs []ProviderConfigReq `json:"provider_configs"`
	}
	payload := VirtualKeyReq{
		CustomerID:      r.CustomerID,
		ProviderConfigs: r.ProviderConfigs,
	}
	args := httpHandlerArgs{
		URL:         url,
		Method:      PUT,
		Payload:     payload,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(args)
	if err != nil {
		return VirtualKey{}, errors.Wrap(err, "Failed to update virtual key")
	}

	var updateVirtualKeyRes struct {
		Message    string     `json:"message"`
		VirtualKey VirtualKey `json:"virtual_key"`
	}
	err = json.Unmarshal(res, &updateVirtualKeyRes)
	if err != nil {
		return VirtualKey{}, errors.Wrap(err, "Failed to unmarshal customer data")
	}

	return updateVirtualKeyRes.VirtualKey, nil
}
