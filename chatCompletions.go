package bifrost

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
)

type CreateAChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CreateAChatCompletionReq struct {
	Model    string                         `json:"model"`
	Messages []CreateAChatCompletionMessage `json:"messages"`
}

type CreateAChatCompletionResChoice struct {
	Index        int64                        `json:"index"`
	FinishReason string                       `json:"finish_reason"`
	Message      CreateAChatCompletionMessage `json:"message"`
}

type CreateAChatCompletionRes struct {
	ID      string                           `json:"id"`
	Choices []CreateAChatCompletionResChoice `json:"choices"`
}

func (c *Client) CreateAChatCompletion(ctx context.Context, r CreateAChatCompletionReq) (CreateAChatCompletionRes, error) {
	url := "/v1/chat/completions"

	payload := r
	args := httpHandlerArgs{
		URL:         url,
		Method:      POST,
		Payload:     payload,
		Credentials: c.Credentials,
	}
	res, err := httpHandler(ctx, args)
	if err != nil {
		return CreateAChatCompletionRes{}, errors.Wrap(err, "Failed to create a chat completion")
	}

	var chatRes CreateAChatCompletionRes
	err = json.Unmarshal(res, &chatRes)
	if err != nil {
		return CreateAChatCompletionRes{}, errors.Wrap(err, "Failed to unmarshal chat completion")
	}

	return chatRes, nil
}
