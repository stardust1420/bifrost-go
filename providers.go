package bifrost

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type Provider struct {
	Name string `json:"name"`
}

type Value struct {
	Value   string `json:"value"`
	EnvVar  string `json:"env_var"`
	FromEnv bool   `json:"from_env"`
}

type Key struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Value      Value     `json:"value"`
	ProviderID int64     `json:"provider_id"`
	Provider   string    `json:"provider"`
	KeyID      string    `json:"key_id"`
	Models     []string  `json:"models"`
	Weight     int64     `json:"weight"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c *Client) ListAllProviders(ctx context.Context) ([]Provider, error) {
	url := "/api/providers"

	args := httpHandlerArgs{
		URL:         url,
		Method:      GET,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(ctx, args)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list customers")
	}

	var listProvidersRes struct {
		Total     int64      `json:"total"`
		Providers []Provider `json:"providers"`
	}
	err = json.Unmarshal(res, &listProvidersRes)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal providers response")
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

func (c *Client) CreateAKeyForAProvider(ctx context.Context, r CreateAKeyForAProviderReq) (CreateAKeyForAProviderRes, error) {
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
	res, err := httpHandler(ctx, args)
	if err != nil {
		return CreateAKeyForAProviderRes{}, errors.Wrap(err, "Failed to create provider key")
	}

	var keyRes CreateAKeyForAProviderRes
	err = json.Unmarshal(res, &keyRes)
	if err != nil {
		return CreateAKeyForAProviderRes{}, errors.Wrap(err, "Failed to unmarshal provider key data")
	}

	return keyRes, nil
}

type GetASpecificKeyForAProviderReq struct {
	Provider string `json:"provider"`
	KeyID    string `json:"key_id"`
}

// Use the below struct instead of Key because the below API returns KeyID in id which have different types
type GetASpecificKeyForAProviderRes struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Value   Value    `json:"value"`
	Models  []string `json:"models"`
	Weight  int64    `json:"weight"`
	Enabled bool     `json:"enabled"`
}

func (c *Client) GetASpecificKeyForAProvider(ctx context.Context, r GetASpecificKeyForAProviderReq) (GetASpecificKeyForAProviderRes, error) {
	url := fmt.Sprintf("/api/providers/%s/keys/%s", r.Provider, r.KeyID)

	args := httpHandlerArgs{
		URL:         url,
		Method:      GET,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(ctx, args)
	if err != nil {
		return GetASpecificKeyForAProviderRes{}, errors.Wrap(err, "Failed to get provider key")
	}

	var key GetASpecificKeyForAProviderRes
	err = json.Unmarshal(res, &key)
	if err != nil {
		return GetASpecificKeyForAProviderRes{}, errors.Wrap(err, "Failed to unmarshal provider key data")
	}

	return key, nil
}

type UpdateAKeyForAProviderReq struct {
	Provider string                         `json:"provider"`
	KeyID    string                         `json:"key_id"`
	Key      GetASpecificKeyForAProviderRes `json:"key"`
}

func (c *Client) UpdateAKeyForAProvider(ctx context.Context, r UpdateAKeyForAProviderReq) (GetASpecificKeyForAProviderRes, error) {
	url := fmt.Sprintf("/api/providers/%s/keys/%s", r.Provider, r.KeyID)

	args := httpHandlerArgs{
		URL:         url,
		Method:      PUT,
		Payload:     r.Key,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(ctx, args)
	if err != nil {
		return GetASpecificKeyForAProviderRes{}, errors.Wrap(err, "Failed to update a key for a provider")
	}

	var key GetASpecificKeyForAProviderRes
	err = json.Unmarshal(res, &key)
	if err != nil {
		return GetASpecificKeyForAProviderRes{}, errors.Wrap(err, "Failed to unmarshal provider key data")
	}

	return key, nil
}

type DeleteAKeyFromAProviderReq struct {
	Provider string `json:"provider"`
	KeyID    string `json:"key_id"`
}

func (c *Client) DeleteAKeyFromAProvider(ctx context.Context, r DeleteAKeyFromAProviderReq) (GetASpecificKeyForAProviderRes, error) {
	url := fmt.Sprintf("/api/providers/%s/keys/%s", r.Provider, r.KeyID)

	args := httpHandlerArgs{
		URL:         url,
		Method:      DELETE,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(ctx, args)
	if err != nil {
		return GetASpecificKeyForAProviderRes{}, errors.Wrap(err, "Failed to delete a key from a provider")
	}

	var key GetASpecificKeyForAProviderRes
	err = json.Unmarshal(res, &key)
	if err != nil {
		return GetASpecificKeyForAProviderRes{}, errors.Wrap(err, "Failed to unmarshal provider key data")
	}

	return key, nil
}
