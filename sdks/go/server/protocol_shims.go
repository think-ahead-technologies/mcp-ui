package mcpuiserver

import (
	"encoding/json"
	"fmt"
)

const (
	// DefaultAdapterBaseURL is the default CDN URL for external adapter scripts
	DefaultAdapterBaseURL = "https://cdn.mcp-ui.dev/adapters"
	// DefaultAdapterVersion is the default version of adapter scripts to use
	DefaultAdapterVersion = "v1"
)

// ProtocolShimGenerator generates script tags for external protocol adapters.
// Each protocol type has its own generator that creates the appropriate
// script reference for loading the external adapter.
type ProtocolShimGenerator interface {
	// GenerateScriptTag returns the HTML script tag to inject into the resource
	GenerateScriptTag() string
	// GetMIMEType returns the MIME type for this protocol
	GetMIMEType() string
}

// GenericProtocolShim generates script tags for the generic MCP-UI protocol.
// No external adapter is needed for the generic protocol.
type GenericProtocolShim struct{}

// GenerateScriptTag returns an empty string as no script is needed for generic protocol
func (g *GenericProtocolShim) GenerateScriptTag() string {
	return "" // No script needed for generic protocol
}

// GetMIMEType returns the standard HTML MIME type
func (g *GenericProtocolShim) GetMIMEType() string {
	return MimeTypeHTML
}

// AppsSdkProtocolShim generates script tags for the ChatGPT/Apps SDK adapter.
// This adapter enables widgets to run in ChatGPT and other Apps SDK environments.
type AppsSdkProtocolShim struct {
	BaseURL string
	Version string
	Config  map[string]interface{}
}

// GenerateScriptTag returns a script tag that loads the Apps SDK adapter from an external URL
func (a *AppsSdkProtocolShim) GenerateScriptTag() string {
	scriptURL := fmt.Sprintf("%s/appssdk-%s.js", a.BaseURL, a.Version)

	// Serialize config to JSON
	configJSON := "{}"
	if a.Config != nil && len(a.Config) > 0 {
		if jsonBytes, err := json.Marshal(a.Config); err == nil {
			configJSON = string(jsonBytes)
		}
	}

	// Generate script tag with configuration in data attribute
	return fmt.Sprintf(`<script src="%s" data-mcp-config='%s'></script>`, scriptURL, configJSON)
}

// GetMIMEType returns the Apps SDK specific MIME type
func (a *AppsSdkProtocolShim) GetMIMEType() string {
	return MimeTypeAppsSdkAdapter // "text/html+skybridge"
}

// McpAppsProtocolShim generates script tags for the MCP Apps SEP adapter.
// This adapter enables widgets to run in MCP Apps environments using the
// Streaming Extensible Protocol (SEP).
type McpAppsProtocolShim struct {
	BaseURL string
	Version string
	Config  map[string]interface{}
}

// GenerateScriptTag returns a script tag that loads the MCP Apps adapter from an external URL
func (m *McpAppsProtocolShim) GenerateScriptTag() string {
	scriptURL := fmt.Sprintf("%s/mcpapps-%s.js", m.BaseURL, m.Version)

	// Serialize config to JSON
	configJSON := "{}"
	if m.Config != nil && len(m.Config) > 0 {
		if jsonBytes, err := json.Marshal(m.Config); err == nil {
			configJSON = string(jsonBytes)
		}
	}

	// Generate script tag with configuration in data attribute
	return fmt.Sprintf(`<script src="%s" data-mcp-config='%s'></script>`, scriptURL, configJSON)
}

// GetMIMEType returns the standard HTML MIME type for MCP Apps
func (m *McpAppsProtocolShim) GetMIMEType() string {
	return MimeTypeMCPAppsAdapter // "text/html"
}

// getProtocolShimGenerator creates the appropriate shim generator based on protocol configuration.
// It handles default values for BaseURL and Version if not specified in the config.
func getProtocolShimGenerator(config *ProtocolConfig) ProtocolShimGenerator {
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = DefaultAdapterBaseURL
	}

	version := config.Version
	if version == "" {
		version = DefaultAdapterVersion
	}

	switch config.Type {
	case ProtocolTypeAppsSDK:
		return &AppsSdkProtocolShim{
			BaseURL: baseURL,
			Version: version,
			Config:  config.Config,
		}
	case ProtocolTypeMCPApps:
		return &McpAppsProtocolShim{
			BaseURL: baseURL,
			Version: version,
			Config:  config.Config,
		}
	default:
		return &GenericProtocolShim{}
	}
}
