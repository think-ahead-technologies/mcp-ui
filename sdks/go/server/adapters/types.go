// Package adapters provides adapter implementations for MCP-UI widgets
// to work in different host environments (Apps SDK, MCP Apps, etc.)
package adapters

// AdapterType identifies which adapter implementation to use.
type AdapterType string

const (
	// AdapterTypeAppsSDK represents the Apps SDK adapter (for ChatGPT, etc.)
	AdapterTypeAppsSDK AdapterType = "appssdk"

	// AdapterTypeMCPApps represents the MCP Apps adapter (for MCP Apps SEP)
	AdapterTypeMCPApps AdapterType = "mcpapps"
)

// Adapter is the interface that all adapter implementations must satisfy.
// Adapters translate MCP-UI protocol messages to host-specific protocols
// by injecting JavaScript runtime code into the widget HTML.
type Adapter interface {
	// GetScript returns the complete <script> tag containing the adapter
	// runtime code with injected configuration.
	GetScript() string

	// GetMIMEType returns the MIME type that should be used for resources
	// using this adapter.
	GetMIMEType() string

	// GetType returns the type identifier for this adapter.
	GetType() string
}
