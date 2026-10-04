package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Dien45/Mobile-agent/internal/core"
)

type ToolDefinition struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}
type ToolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}
type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}
type ChatMessage struct {
	Role       string     `json:"role"`
	Content    any        `json:"content,omitempty"`
	Name       string     `json:"name,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
}
type ChatResult struct { Message ChatMessage; PromptTokens, CompletionTokens int }

type Client struct { HTTP *http.Client }
func New() *Client { return &Client{HTTP: &http.Client{Timeout: 5 * time.Minute}} }

func endpoint(base, custom, fallback string) (string, error) {
	if custom == "" { custom = fallback }
	base = strings.TrimRight(base, "/") + "/"
	custom = strings.TrimLeft(custom, "/")
	u, err := url.Parse(base); if err != nil { return "", err }
	r, err := url.Parse(custom); if err != nil { return "", err }
	return u.ResolveReference(r).String(), nil
}

func (c *Client) request(ctx context.Context, p core.ProviderProfile, method, endpoint string, body any) ([]byte, error) {
	var reader io.Reader
	if body != nil { b, err := json.Marshal(body); if err != nil { return nil, err }; reader = bytes.NewReader(b) }
	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader); if err != nil { return nil, err }
	req.Header.Set("Accept", "application/json"); if body != nil { req.Header.Set("Content-Type", "application/json") }
	if p.APIKey != "" { req.Header.Set("Authorization", "Bearer "+p.APIKey) }
	for k, v := range p.Headers { req.Header.Set(k, v) }
	resp, err := c.HTTP.Do(req); if err != nil { return nil, err }
	defer resp.Body.Close(); data, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20)); if err != nil { return nil, err }
	if resp.StatusCode < 200 || resp.StatusCode >= 300 { return nil, fmt.Errorf("provider returned %s: %s", resp.Status, strings.TrimSpace(string(data))) }
	return data, nil
}

func (c *Client) Models(ctx context.Context, p core.ProviderProfile) ([]core.Model, error) {
	ep, err := endpoint(p.BaseURL, p.ModelsPath, "models"); if err != nil { return nil, err }
	b, err := c.request(ctx, p, http.MethodGet, ep, nil); if err != nil { return nil, err }
	var raw struct { Data []struct { ID, OwnedBy string } `json:"data"` }
	if err := json.Unmarshal(b, &raw); err != nil { return nil, err }
	out := make([]core.Model, 0, len(raw.Data)); for _, m := range raw.Data { out = append(out, core.Model{ID:m.ID, OwnedBy:m.OwnedBy}) }
	return out, nil
}

func (c *Client) Chat(ctx context.Context, p core.ProviderProfile, model string, messages []ChatMessage, tools []ToolDefinition) (ChatResult, error) {
	ep, err := endpoint(p.BaseURL, p.ChatPath, "chat/completions"); if err != nil { return ChatResult{}, err }
	payload := map[string]any{"model": model, "messages": messages, "stream": false}
	if len(tools) > 0 { payload["tools"] = tools; payload["tool_choice"] = "auto" }
	b, err := c.request(ctx, p, http.MethodPost, ep, payload); if err != nil { return ChatResult{}, err }
	var raw struct {
		Choices []struct { Message ChatMessage `json:"message"` } `json:"choices"`
		Usage struct {
			Prompt     int `json:"prompt_tokens"`
			Completion int `json:"completion_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(b, &raw); err != nil { return ChatResult{}, fmt.Errorf("decode chat response: %w", err) }
	if len(raw.Choices) == 0 { return ChatResult{}, fmt.Errorf("provider returned no choices") }
	return ChatResult{Message: raw.Choices[0].Message, PromptTokens: raw.Usage.Prompt, CompletionTokens: raw.Usage.Completion}, nil
}
