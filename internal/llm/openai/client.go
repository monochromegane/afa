package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/monochromegane/afa/internal/payload"
)

const (
	API_ENDPOINT              = "https://api.openai.com"
	API_CHAT_COMPLETIONS_PATH = "/v1/chat/completions"
)

var (
	dataPrefix = []byte("data: ")
	dataSuffix = []byte("[DONE]")
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: Status Code %d.\n", resp.StatusCode)
	}

	var output Response
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return nil, err
	}

	if len(output.Choices) > 0 {
		if refusal := output.Choices[0].Message.Refusal; refusal != "" {
			return nil, fmt.Errorf("Refused to respond %s.", refusal)
		}
	}

	return c.repackResponse(&output), nil
}

func (c *Client) ChatCompletionStream(request *payload.Request, ctx context.Context, onData func(*payload.Response) error) error {
	repacked := c.repackRequest(request)
	repacked.Stream = true
	req, err := c.newJsonRequest(ctx, repacked)
	if err != nil {
		return err
	}

	resp, err := c.postJson(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error: Status Code %d.\n", resp.StatusCode)
		}
		return fmt.Errorf("Error: Status Code %d, Response Body: %s.\n", resp.StatusCode, string(body))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			return err
		}

		line = bytes.TrimSpace(line)
		if !bytes.HasPrefix(line, dataPrefix) {
			continue
		}
		line = bytes.TrimPrefix(line, dataPrefix)

		if bytes.HasSuffix(line, dataSuffix) {
			break
		}

		var output ResponseStream
		if err := json.Unmarshal(line, &output); err != nil {
			return err
		}

		if len(output.Choices) > 0 {
			if refusal := output.Choices[0].Delta.Refusal; refusal != "" {
				return fmt.Errorf("Refused to respond %s.", refusal)
			}
		}

		if err := onData(c.repackResponseStream(&output)); err != nil {
			return err
		}
	}

	return nil
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
	repacked := &Request{
		Model:    request.Model,
		Messages: messages,
	}

	if request.JsonSchema != nil {
		repacked.ResponseFormat = &ResponseFormat{
			Type: "json_schema",
			JsonSchema: &JsonSchema{
				Name:   request.JsonSchema.Name,
				Strict: true,
				Schema: request.JsonSchema.Schema,
			},
		}
	}
	return repacked
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

func (c *Client) repackResponseStream(response *ResponseStream) *payload.Response {
	var message payload.Message
	if len(response.Choices) > 0 {
		message.Role = response.Choices[0].Delta.Role
		message.Content = response.Choices[0].Delta.Content
	}
	return &payload.Response{
		Message: &message,
	}
}

type Request struct {
	Model          string          `json:"model"`
	Messages       []*Message      `json:"messages"`
	Stream         bool            `json:"stream"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
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
	Message *Message `json:"message"`
}

type ResponseStream struct {
	Choices []*ChoiceStream `json:"choices"`
}

type ChoiceStream struct {
	Delta Message `json:"delta"`
}

type ResponseFormat struct {
	Type       string      `json:"type"`
	JsonSchema *JsonSchema `json:"json_schema"`
}

type JsonSchema struct {
	Name   string           `json:"name"`
	Strict bool             `json:"strict"`
	Schema *json.RawMessage `json:"schema"`
}
