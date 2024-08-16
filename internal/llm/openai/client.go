package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/monochromegane/aiforall/internal/payload"
)

const (
	API_ENDPOINT              = "https://api.openai.com"
	API_CHAT_COMPLETIONS_PATH = "/v1/chat/completions"
)

type Client struct {
	Endpoint string
}

func NewClient() *Client {
	return &Client{
		Endpoint: API_ENDPOINT,
	}
}

func (c *Client) ChatCompletion(request *payload.Request, ctx context.Context) (*payload.Response, error) {
	repacked := c.repackRequest(request)
	req, err := c.newJsonRequest(ctx, repacked)
	if err != nil {
		return nil, err
	}

	resp, err := c.postJson(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var output Response
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}
	return c.repackResponse(&output), nil
}

func (c *Client) newJsonRequest(ctx context.Context, request *Request) (*http.Request, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&request); err != nil {
		return nil, err
	}

	endpoint, err := url.JoinPath(c.Endpoint, API_CHAT_COMPLETIONS_PATH)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, &buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if apiKey, ok := ctx.Value("openai-api-key").(string); ok {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	}

	return req, nil
}

func (c *Client) postJson(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(req)
}

func (c *Client) repackRequest(request *payload.Request) *Request {
	messages := make([]*Message, len(request.Messages))
	for i, message := range request.Messages {
		messages[i] = &Message{
			Role:    message.Role,
			Content: message.Content,
		}
	}
	return &Request{
		Model:    request.Model,
		Messages: messages,
	}
}

func (c *Client) repackResponse(response *Response) *payload.Response {
	var message payload.Message
	if len(response.Choices) > 0 {
		message.Role = response.Choices[0].Message.Role
		message.Content = response.Choices[0].Message.Content
	}
	return &payload.Response{
		Message: &message,
	}
}

type Request struct {
	Model    string     `json:"model"`
	Messages []*Message `json:"messages"`
	Stream   bool       `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Refusal string `json:"refusal,omitempty"`
}

type Response struct {
	Choices []*Choice `json:"choices"`
}

type Choice struct {
	Index        int      `json:"index"`
	FinishReason string   `json:"finish_reason"`
	Message      *Message `json:"message"`
}
