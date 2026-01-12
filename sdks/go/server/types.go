// Package mcpuiserver provides a Go SDK for creating MCP-UI server resources.
// It enables Go-based MCP servers to create UI resources with HTML content,
// external URLs, and Remote DOM components.
package mcpuiserver

import (
	"errors"
	"fmt"
	"strings"
)

// URI scheme and metadata constants
const (
	// URIScheme is the required prefix for all UI resource URIs
	URIScheme = "ui://"

	// UIMetadataPrefix is prepended to UI-specific metadata keys
	UIMetadataPrefix = "mcpui.dev/ui-"

	// MIME type constants
	MimeTypeHTML           = "text/html"
	MimeTypeURIList        = "text/uri-list"
	MimeTypeRemoteDomReact = "application/vnd.mcp-ui.remote-dom+javascript; framework=react"
	MimeTypeRemoteDomWC    = "application/vnd.mcp-ui.remote-dom+javascript; framework=webcomponents"

	// Adapter MIME type constants
	MimeTypeAppsSdkAdapter = "text/html+skybridge"
	MimeTypeMCPAppsAdapter = "text/html"
)

// UIMetadataKey defines standard metadata keys
const (
	UIMetadataKeyPreferredFrameSize = "preferred-frame-size"
	UIMetadataKeyInitialRenderData  = "initial-render-data"
)

// Error definitions
var (
	ErrInvalidURI       = errors.New("URI must start with 'ui://'")
	ErrEmptyHTMLString  = errors.New("htmlString must be provided as a non-empty string when content type is 'rawHtml'")
	ErrEmptyIframeURL   = errors.New("iframeUrl must be provided as a non-empty string when content type is 'externalUrl'")
	ErrEmptyScript      = errors.New("script must be provided as a non-empty string when content type is 'remoteDom'")
	ErrInvalidFramework = errors.New("framework must be 'react' or 'webcomponents'")
	ErrInvalidEncoding  = errors.New("encoding must be 'text' or 'blob'")
	ErrNilContent       = errors.New("content cannot be nil")
)

// InvalidURIError wraps the URI validation error with the actual URI
type InvalidURIError struct {
	URI string
}

func (e *InvalidURIError) Error() string {
	return fmt.Sprintf("URI must start with 'ui://' but got: %s", e.URI)
}

func (e *InvalidURIError) Is(target error) bool {
	return target == ErrInvalidURI
}

// ContentType represents the type of UI content
type ContentType string

const (
	ContentTypeRawHTML     ContentType = "rawHtml"
	ContentTypeExternalURL ContentType = "externalUrl"
	ContentTypeRemoteDOM   ContentType = "remoteDom"
)

// Encoding represents the resource encoding type
type Encoding string

const (
	EncodingText Encoding = "text"
	EncodingBlob Encoding = "blob"
)

// RemoteDOMFramework represents the framework used for Remote DOM rendering
type RemoteDOMFramework string

const (
	FrameworkReact         RemoteDOMFramework = "react"
	FrameworkWebComponents RemoteDOMFramework = "webcomponents"
)

// ResourceContentPayload is the interface for content payloads
type ResourceContentPayload interface {
	contentType() ContentType
	validate() error
}

// RawHTMLPayload represents raw HTML content
type RawHTMLPayload struct {
	Type       ContentType `json:"type"`
	HTMLString string      `json:"htmlString"`
}

func (p *RawHTMLPayload) contentType() ContentType {
	return ContentTypeRawHTML
}

func (p *RawHTMLPayload) validate() error {
	if p.HTMLString == "" {
		return ErrEmptyHTMLString
	}
	return nil
}

// ExternalURLPayload represents external URL content
type ExternalURLPayload struct {
	Type      ContentType `json:"type"`
	IframeURL string      `json:"iframeUrl"`
}

func (p *ExternalURLPayload) contentType() ContentType {
	return ContentTypeExternalURL
}

func (p *ExternalURLPayload) validate() error {
	if p.IframeURL == "" {
		return ErrEmptyIframeURL
	}
	return nil
}

// RemoteDOMPayload represents remote DOM content
type RemoteDOMPayload struct {
	Type      ContentType        `json:"type"`
	Script    string             `json:"script"`
	Framework RemoteDOMFramework `json:"framework"`
}

func (p *RemoteDOMPayload) contentType() ContentType {
	return ContentTypeRemoteDOM
}

func (p *RemoteDOMPayload) validate() error {
	if p.Script == "" {
		return ErrEmptyScript
	}
	if p.Framework != FrameworkReact && p.Framework != FrameworkWebComponents {
		return ErrInvalidFramework
	}
	return nil
}

// ResourceContent represents the actual resource content
type ResourceContent struct {
	URI      string                 `json:"uri"`
	MimeType string                 `json:"mimeType"`
	Text     string                 `json:"text,omitempty"`
	Blob     string                 `json:"blob,omitempty"`
	Meta     map[string]interface{} `json:"_meta,omitempty"`
}

// UIResource represents a UI resource for MCP responses
type UIResource struct {
	Type        string                 `json:"type"` // Always "resource"
	Resource    ResourceContent        `json:"resource"`
	Annotations map[string]interface{} `json:"annotations,omitempty"`
	Meta        map[string]interface{} `json:"_meta,omitempty"`
}

// CreateUIResourceOptions contains all options for creating a UI resource
type CreateUIResourceOptions struct {
	URI                   string
	Content               ResourceContentPayload
	Encoding              Encoding
	UIMetadata            map[string]interface{}
	Metadata              map[string]interface{}
	ResourceProps         map[string]interface{}
	EmbeddedResourceProps map[string]interface{}
	Protocol              *ProtocolConfig // Server-side protocol selection with external adapter scripts
}

// ProtocolType defines the UI protocol to use for a session
type ProtocolType string

const (
	// ProtocolTypeGeneric represents standard MCP-UI protocol (no adapter needed)
	ProtocolTypeGeneric ProtocolType = "generic"
	// ProtocolTypeAppsSDK represents ChatGPT/Apps SDK protocol
	ProtocolTypeAppsSDK ProtocolType = "appssdk"
	// ProtocolTypeMCPApps represents MCP Apps SEP (Streaming Extensible Protocol)
	ProtocolTypeMCPApps ProtocolType = "mcpapps"
)

// ProtocolConfig holds protocol-specific configuration for server-side protocol selection.
// This eliminates the need to inject large adapter runtimes into HTML content,
// instead using external script references.
type ProtocolConfig struct {
	// Type specifies which protocol to use
	Type ProtocolType
	// Version specifies the adapter version (e.g., "v1")
	Version string
	// BaseURL specifies the base URL for external adapter scripts
	BaseURL string
	// Config contains protocol-specific settings
	Config map[string]interface{}
}

// Option is a functional option for CreateUIResourceOptions
type Option func(*CreateUIResourceOptions)

// WithUIMetadata sets UI-specific metadata (will be prefixed with "mcpui.dev/ui-")
func WithUIMetadata(metadata map[string]interface{}) Option {
	return func(o *CreateUIResourceOptions) {
		o.UIMetadata = metadata
	}
}

// WithMetadata sets custom metadata
func WithMetadata(metadata map[string]interface{}) Option {
	return func(o *CreateUIResourceOptions) {
		o.Metadata = metadata
	}
}

// WithResourceProps sets additional resource properties
func WithResourceProps(props map[string]interface{}) Option {
	return func(o *CreateUIResourceOptions) {
		o.ResourceProps = props
	}
}

// WithEmbeddedResourceProps sets embedded resource properties
func WithEmbeddedResourceProps(props map[string]interface{}) Option {
	return func(o *CreateUIResourceOptions) {
		o.EmbeddedResourceProps = props
	}
}

// WithProtocol sets the protocol type for this resource.
// This enables server-side protocol selection using external adapter scripts,
// eliminating the need to inject large adapter runtimes into HTML content.
// Example: WithProtocol(ProtocolTypeAppsSDK)
func WithProtocol(protocol ProtocolType) Option {
	return func(o *CreateUIResourceOptions) {
		if o.Protocol == nil {
			o.Protocol = &ProtocolConfig{}
		}
		o.Protocol.Type = protocol
	}
}

// WithProtocolConfig sets the full protocol configuration for this resource.
// This allows complete control over protocol selection, versioning, and configuration.
// Example: WithProtocolConfig(&ProtocolConfig{Type: ProtocolTypeAppsSDK, Version: "v1"})
func WithProtocolConfig(config *ProtocolConfig) Option {
	return func(o *CreateUIResourceOptions) {
		o.Protocol = config
	}
}

// WithProtocolVersion sets the adapter version for the protocol.
// This determines which version of the external adapter script to load.
// Example: WithProtocolVersion("v2")
func WithProtocolVersion(version string) Option {
	return func(o *CreateUIResourceOptions) {
		if o.Protocol == nil {
			o.Protocol = &ProtocolConfig{}
		}
		o.Protocol.Version = version
	}
}

// WithProtocolBaseURL sets the base URL for external adapter scripts.
// Use this to self-host adapter scripts instead of using the default CDN.
// Example: WithProtocolBaseURL("https://myserver.com/adapters")
func WithProtocolBaseURL(baseURL string) Option {
	return func(o *CreateUIResourceOptions) {
		if o.Protocol == nil {
			o.Protocol = &ProtocolConfig{}
		}
		o.Protocol.BaseURL = baseURL
	}
}

// validateURI validates that a URI starts with the ui:// scheme
func validateURI(uri string) error {
	if !strings.HasPrefix(uri, URIScheme) {
		return &InvalidURIError{URI: uri}
	}
	return nil
}
