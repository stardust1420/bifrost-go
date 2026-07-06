package bifrost

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

// Http helper
type httpMethod string

const (
	GET    httpMethod = "GET"
	POST   httpMethod = "POST"
	PUT    httpMethod = "PUT"
	DELETE httpMethod = "DELETE"
)

func (h httpMethod) String() string {
	return string(h)
}

type httpHandlerArgs struct {
	URL         string
	Method      httpMethod
	Payload     any
	Credentials Credentials
}

func httpHandler(ctx context.Context, args httpHandlerArgs) ([]byte, error) {
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
	req, err := http.NewRequestWithContext(ctx, args.Method.String(), url, body)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create new http request")
	}

	// Request Headers
	req.Header.Add("Content-Type", "application/json")

	credentialsPresent := false

	if args.Credentials.Username != "" && args.Credentials.Password != "" {
		req.SetBasicAuth(args.Credentials.Username, args.Credentials.Password)
		credentialsPresent = true

	}
	if args.Credentials.VirtualKey != "" {
		req.Header.Add("x-bf-vk", args.Credentials.VirtualKey)
		credentialsPresent = true
	}

	if !credentialsPresent {
		return nil, errors.Wrap(err, "Credentials not found")
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
