package core

import (
	"encoding/json"
	"time"
)

type Session struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Provider  string    `json:"provider"`
	Model     string    `json:"model"`
	Workspace string    `json:"workspace"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Message struct {
	ID        string    `json:"id"`
	SessionID string    `json:"session_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Name      string    `json:"name,omitempty"`
	ToolCallID string   `json:"tool_call_id,omitempty"`
	ToolCalls json.RawMessage `json:"tool_calls,omitempty"`
	Provider  string    `json:"provider,omitempty"`
	Model     string    `json:"model,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ProviderProfile struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	BaseURL    string            `json:"base_url"`
	APIKey     string            `json:"api_key,omitempty"`
	ModelsPath string            `json:"models_path,omitempty"`
	ChatPath   string            `json:"chat_path,omitempty"`
	Headers    map[string]string `json:"headers,omitempty"`
}

type Model struct {
	ID      string `json:"id"`
	OwnedBy string `json:"owned_by,omitempty"`
}

type Skill struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	Enabled     bool      `json:"enabled"`
	InstalledAt time.Time `json:"installed_at"`
}

type TaskStatus struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
	ToolCalls int    `json:"tool_calls"`
	Error     string `json:"error,omitempty"`
}
