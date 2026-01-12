// Package mcpapps provides the MCP Apps adapter for MCP-UI widgets.
// It enables widgets to run in MCP Apps SEP (Streaming Extensible Protocol) environments.
package mcpapps

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MCP-UI-Org/mcp-ui/sdks/go/server/adapters"
)

// Configuration errors
var (
	ErrInvalidTimeout = errors.New("timeout must be positive")
)

// Config holds the configuration for the MCP Apps adapter.
type Config struct {
	// Timeout in milliseconds for async operations (default: 30000)
	Timeout int
}

// Validate validates the adapter configuration.
func (c *Config) Validate() error {
	if c.Timeout <= 0 {
		return ErrInvalidTimeout
	}
	return nil
}

// Option is a functional option for configuring the MCP Apps adapter.
type Option func(*Config)

// WithTimeout sets the timeout in milliseconds for async operations.
func WithTimeout(milliseconds int) Option {
	return func(c *Config) {
		c.Timeout = milliseconds
	}
}

// Adapter implements the MCP Apps adapter for MCP-UI widgets.
type Adapter struct {
	config *Config
}

// NewAdapter creates a new MCP Apps adapter with the provided options.
// Default configuration:
//   - Timeout: 30000ms
func NewAdapter(opts ...Option) (*Adapter, error) {
	config := &Config{
		Timeout: 30000,
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
		"<script>\nconst config = %s;\n%s\nwindow.McpAppsAdapter = { init: initAdapter, initWithConfig: () => initAdapter(config), uninstall: uninstallAdapter };\nwindow.McpAppsAdapter.initWithConfig();\n</script>",
		configStr,
		adapterRuntimeScript,
	)
}

// GetMIMEType returns the MIME type for MCP Apps adapter resources.
func (a *Adapter) GetMIMEType() string {
	return "text/html"
}

// GetType returns the adapter type identifier.
func (a *Adapter) GetType() string {
	return string(adapters.AdapterTypeMCPApps)
}

// serializableConfig returns a config map with only serializable fields.
func (a *Adapter) serializableConfig() map[string]interface{} {
	return map[string]interface{}{
		"timeout": a.config.Timeout,
	}
}
