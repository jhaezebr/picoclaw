// PicoClaw - Ultra-lightweight personal AI agent
// License: MIT
//
// Copyright (c) 2026 PicoClaw contributors

package providers

import (
	"fmt"
	"strings"

	"github.com/sipeed/picoclaw/pkg/config"
)

// ExtractProtocol extracts the protocol prefix and model identifier from a model string.
// If no prefix is specified, it defaults to "openai".
// Examples:
//   - "openai/gpt-4o" -> ("openai", "gpt-4o")
//   - "anthropic/claude-sonnet-4.6" -> ("anthropic", "claude-sonnet-4.6")
//   - "gpt-4o" -> ("openai", "gpt-4o")  // default protocol
func ExtractProtocol(model string) (protocol, modelID string) {
	model = strings.TrimSpace(model)
	protocol, modelID, found := strings.Cut(model, "/")
	if !found {
		return "openai", model
	}
	return protocol, modelID
}

// CreateProviderFromConfig creates a provider based on the ModelConfig.
// It uses the protocol prefix in the Model field to determine which provider to create.
// Supported protocols: openai, anthropic, antigravity, claude-cli, codex-cli, github-copilot
// Returns the provider, the model ID (without protocol prefix), and any error.
func CreateProviderFromConfig(cfg *config.ModelConfig) (LLMProvider, string, error) {
	if cfg == nil {
		return nil, "", fmt.Errorf("config is nil")
	}

	if cfg.Model == "" {
		return nil, "", fmt.Errorf("model is required")
	}

	protocol, modelID := ExtractProtocol(cfg.Model)

	switch protocol {
	case "gemini", "ollama", "vllm", "mistral":
		// All other OpenAI-compatible HTTP providers
		if cfg.APIKey == "" && cfg.APIBase == "" {
			return nil, "", fmt.Errorf("api_key or api_base is required for HTTP-based protocol %q", protocol)
		}
		apiBase := cfg.APIBase
		if apiBase == "" {
			apiBase = getDefaultAPIBase(protocol)
		}
		return NewHTTPProviderWithMaxTokensFieldAndRequestTimeout(
			cfg.APIKey,
			apiBase,
			cfg.Proxy,
			cfg.MaxTokensField,
			cfg.RequestTimeout,
		), modelID, nil

	default:
		return nil, "", fmt.Errorf("unknown protocol %q in model %q", protocol, cfg.Model)
	}
}

// getDefaultAPIBase returns the default API base URL for a given protocol.
func getDefaultAPIBase(protocol string) string {
	switch protocol {
	case "gemini":
		return "https://generativelanguage.googleapis.com/v1beta"
	case "ollama":
		return "http://localhost:11434/v1"
	case "vllm":
		return "http://localhost:8000/v1"
	case "mistral":
		return "https://api.mistral.ai/v1"
	default:
		return ""
	}
}
