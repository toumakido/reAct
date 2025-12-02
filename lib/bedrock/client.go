package bedrock

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/toumakido/reAct/lib/types"
)

type Client struct {
	client  *bedrockruntime.Client
	modelID string
}

type invokeRequest struct {
	AnthropicVersion string          `json:"anthropic_version"`
	MaxTokens        int             `json:"max_tokens"`
	System           string          `json:"system,omitempty"`
	Messages         []types.Message `json:"messages"`
}

type invokeResponse struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Usage struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

type InvokeResult struct {
	Text         string
	InputTokens  int
	OutputTokens int
}

// NewClient creates a new Bedrock client
func NewClient(ctx context.Context) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	return &Client{
		client:  client,
		modelID: "global.anthropic.claude-haiku-4-5-20251001-v1:0",
	}, nil
}

// InvokeModel sends messages to Claude and returns the response
func (c *Client) InvokeModel(ctx context.Context, systemPrompt string, messages []types.Message) (*InvokeResult, error) {
	request := invokeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        4096,
		System:           systemPrompt,
		Messages:         messages,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	output, err := c.client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(c.modelID),
		ContentType: aws.String("application/json"),
		Body:        requestBody,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to invoke model: %w", err)
	}

	var response invokeResponse
	if err := json.Unmarshal(output.Body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Content) == 0 {
		return nil, fmt.Errorf("no content in response")
	}

	return &InvokeResult{
		Text:         response.Content[0].Text,
		InputTokens:  response.Usage.InputTokens,
		OutputTokens: response.Usage.OutputTokens,
	}, nil
}
