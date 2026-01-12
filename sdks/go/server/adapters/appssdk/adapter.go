// Package appssdk provides the Apps SDK adapter for MCP-UI widgets.
// It enables widgets to run in ChatGPT and other Apps SDK environments.
package appssdk

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MCP-UI-Org/mcp-ui/sdks/go/server/adapters"
)

// Configuration errors
var (
	ErrInvalidTimeout        = errors.New("timeout must be positive")
	ErrInvalidIntentHandling = errors.New("intentHandling must be 'prompt' or 'ignore'")
)

// Config holds the configuration for the Apps SDK adapter.
type Config struct {
	// Timeout in milliseconds for async operations (default: 30000)
	Timeout int

	// IntentHandling specifies how to handle intent messages:
	// - "prompt": Convert intent to prompt message (default)
	// - "ignore": Ignore intent messages
	IntentHandling string

	// HostOrigin for MessageEvents (default: empty, uses window.location.origin)
	HostOrigin string
}

// Validate validates the adapter configuration.
func (c *Config) Validate() error {
	if c.Timeout <= 0 {
		return ErrInvalidTimeout
	}
	if c.IntentHandling != "prompt" && c.IntentHandling != "ignore" {
		return ErrInvalidIntentHandling
	}
	return nil
}

// Option is a functional option for configuring the Apps SDK adapter.
type Option func(*Config)

// WithTimeout sets the timeout in milliseconds for async operations.
func WithTimeout(milliseconds int) Option {
	return func(c *Config) {
		c.Timeout = milliseconds
	}
}

// WithIntentHandling sets how to handle intent messages.
// Valid values are "prompt" (convert to prompt) or "ignore" (ignore intents).
func WithIntentHandling(handling string) Option {
	return func(c *Config) {
		c.IntentHandling = handling
	}
}

// WithHostOrigin sets the host origin for MessageEvents.
func WithHostOrigin(origin string) Option {
	return func(c *Config) {
		c.HostOrigin = origin
	}
}

// Adapter implements the Apps SDK adapter for MCP-UI widgets.
type Adapter struct {
	config *Config
}

// NewAdapter creates a new Apps SDK adapter with the provided options.
// Default configuration:
//   - Timeout: 30000ms
//   - IntentHandling: "prompt"
//   - HostOrigin: "" (uses window.location.origin)
func NewAdapter(opts ...Option) (*Adapter, error) {
	config := &Config{
		Timeout:        30000,
		IntentHandling: "prompt",
		HostOrigin:     "",
	}

	for _, opt := range opts {
		opt(config)
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &Adapter{config: config}, nil
}

// GetScript returns the complete <script> tag containing the adapter runtime
// with the injected configuration.
func (a *Adapter) GetScript() string {
	configJSON := a.serializableConfig()
	configStr, _ := json.Marshal(configJSON)

	return fmt.Sprintf(
		"<script>\nconst config = %s;\n%s\nwindow.MCPUIAppsSdkAdapter.initWithConfig();\n</script>",
		configStr,
		adapterRuntimeScript,
	)
}

// GetMIMEType returns the MIME type for Apps SDK adapter resources.
func (a *Adapter) GetMIMEType() string {
	// Using the constant from the parent package
	return "text/html+skybridge"
}

// GetType returns the adapter type identifier.
func (a *Adapter) GetType() string {
	return string(adapters.AdapterTypeAppsSDK)
}

// serializableConfig returns a config map with only serializable fields.
// The logger field is not serialized since it can't be JSON-marshaled.
func (a *Adapter) serializableConfig() map[string]interface{} {
	config := map[string]interface{}{
		"timeout":        a.config.Timeout,
		"intentHandling": a.config.IntentHandling,
	}

	if a.config.HostOrigin != "" {
		config["hostOrigin"] = a.config.HostOrigin
	}

	return config
}
