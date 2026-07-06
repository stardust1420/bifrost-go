package bifrost

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

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

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProviderConfig struct {
	ID                int64    `json:"id"`
	VirtualKeyID      string   `json:"virtual_key_id"`
	Provider          string   `json:"provider"`
	Weight            *int64   `json:"weight"`
	AllowedModels     []string `json:"allowed_models"`
	BlacklistedModels []string `json:"blacklisted_models"`
	AllowAllKeys      bool     `json:"allow_all_keys"`
	Keys              []Key    `json:"keys"`
}

type CreateCustomerReq struct {
	Name string `json:"name"`
}

func (c *Client) CreateCustomer(ctx context.Context, r CreateCustomerReq) (Customer, error) {
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
	res, err := httpHandler(ctx, args)
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

func (c *Client) ListCustomers(ctx context.Context) ([]Customer, error) {
	url := "/api/governance/customers"

	args := httpHandlerArgs{
		URL:         url,
		Method:      GET,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(ctx, args)
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

func (c *Client) GetCustomer(ctx context.Context, r GetCustomerReq) (Customer, error) {
	url := fmt.Sprintf("/api/governance/customers/%s", r.ID)

	args := httpHandlerArgs{
		URL:         url,
		Method:      GET,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(ctx, args)
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
	CustomerID  string `json:"customer_id"`
	IsActive    bool   `json:"is_active"`
}

func (c *Client) CreateVirtualKey(ctx context.Context, r CreateVirtualKeyReq) (VirtualKey, error) {
	url := "/api/governance/virtual-keys"
	payload := VirtualKey{
		Name:        r.Name,
		Description: r.Description,
		CustomerID:  r.CustomerID,
		IsActive:    r.IsActive,
	}

	args := httpHandlerArgs{
		URL:         url,
		Method:      POST,
		Payload:     payload,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(ctx, args)
	if err != nil {
		return VirtualKey{}, errors.Wrap(err, "Failed to create virtual key")
	}

	var createVirtualKeyRes struct {
		Message    string     `json:"message"`
		VirtualKey VirtualKey `json:"virtual_key"`
	}
	err = json.Unmarshal(res, &createVirtualKeyRes)
	if err != nil {
		return VirtualKey{}, errors.Wrap(err, "Failed to unmarshal virtual key data")
	}

	return createVirtualKeyRes.VirtualKey, nil
}

type GetVirtualKeyReq struct {
	ID string `json:"id"`
}

func (c *Client) GetVirtualKey(ctx context.Context, r GetVirtualKeyReq) (VirtualKey, error) {
	url := fmt.Sprintf("/api/governance/virtual-keys/%s", r.ID)

	args := httpHandlerArgs{
		URL:         url,
		Method:      GET,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(ctx, args)
	if err != nil {
		return VirtualKey{}, errors.Wrap(err, "Failed to get virtual key")
	}

	var getVirtualKeyRes struct {
		VirtualKey VirtualKey `json:"virtual_key"`
	}
	err = json.Unmarshal(res, &getVirtualKeyRes)
	if err != nil {
		return VirtualKey{}, errors.Wrap(err, "Failed to unmarshal virtual key response")
	}

	return getVirtualKeyRes.VirtualKey, nil
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

func (c *Client) UpdateVirtualKey(ctx context.Context, r UpdateVirtualKeyReq) (VirtualKey, error) {
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
	res, err := httpHandler(ctx, args)
	if err != nil {
		return VirtualKey{}, errors.Wrap(err, "Failed to update virtual key")
	}

	var updateVirtualKeyRes struct {
		Message    string     `json:"message"`
		VirtualKey VirtualKey `json:"virtual_key"`
	}
	err = json.Unmarshal(res, &updateVirtualKeyRes)
	if err != nil {
		return VirtualKey{}, errors.Wrap(err, "Failed to unmarshal virtual key data")
	}

	return updateVirtualKeyRes.VirtualKey, nil
}
