# MCP-UI Server SDK for Go

Go SDK for creating MCP-UI server resources. This SDK enables Go-based MCP servers to create UI resources with HTML content, external URLs, and Remote DOM components.

## Installation

```bash
go get github.com/MCP-UI-Org/mcp-ui/sdks/go/server
```

## Quick Start

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"

    mcpuiserver "github.com/MCP-UI-Org/mcp-ui/sdks/go/server"
)

func main() {
    // Create a simple HTML UI resource
    resource, err := mcpuiserver.CreateUIResource(
        "ui://greeting",
        &mcpuiserver.RawHTMLPayload{
            Type:       mcpuiserver.ContentTypeRawHTML,
            HTMLString: "<h1>Hello, World!</h1>",
        },
        mcpuiserver.EncodingText,
    )
    if err != nil {
        log.Fatal(err)
    }

    // Serialize to JSON for MCP response
    jsonBytes, _ := json.MarshalIndent(resource, "", "  ")
    fmt.Println(string(jsonBytes))
}
```

## Features

- **Three content types:**
  - Raw HTML (inline HTML strings)
  - External URL (iframe URLs)
  - Remote DOM (JavaScript components with React/WebComponents)
- **Two encoding options:** text or blob (base64)
- **Metadata handling:** UI-specific metadata with automatic prefixing
- **Helper functions:** UI action results (tool calls, prompts, links, intents, notifications)
- **Strong typing:** Type-safe API with validation
- **No dependencies:** Uses only Go standard library
- **MCP Apps Standard:** Full support for MCP Apps SEP protocol (version 2025-11-21)

## MCP Apps Standard Support

This SDK fully supports the MCP Apps Standard (protocol version 2025-11-21). Key features:

- **Standard MIME Type**: Uses `text/html;profile=mcp-app` for MCP Apps adapter
- **Protocol Message Types**: Complete type definitions for MCP-UI protocol messages
- **Render Data Support**: Type-safe render data structures for widget initialization
- **Protocol Version Tracking**: Built-in protocol version constants

### Protocol Constants

```go
// Protocol version
mcpuiserver.ProtocolVersion // "2025-11-21"

// Resource URI metadata key (for MCP Apps hosts)
mcpuiserver.ResourceURIMetaKey // "ui/resourceUri"

// MIME types
mcpuiserver.MimeTypeMCPAppsAdapter // "text/html;profile=mcp-app"
```

### Protocol Message Types

The SDK provides complete type definitions for MCP-UI protocol messages:

```go
// Message type constants
mcpuiserver.MessageTypeToolCall
mcpuiserver.MessageTypePrompt
mcpuiserver.MessageTypeLifecycleReady
mcpuiserver.MessageTypeSizeChange
mcpuiserver.MessageTypeRequestData
mcpuiserver.MessageTypeRequestRenderData
// ... and more

// Message structures with helper constructors
msgID := "msg-123"
readyMsg := mcpuiserver.NewLifecycleReadyMessage(&msgID)
```

### Render Data

Pass initialization data to widgets using the `RenderData` type:

```go
renderData := mcpuiserver.RenderData{
    Locale:      "en-US",
    Theme:       "dark",
    DisplayMode: mcpuiserver.DisplayModeInline,
    MaxHeight:   600,
    ToolInput: map[string]interface{}{
        "query": "user data",
    },
}

resource, err := mcpuiserver.CreateUIResource(
    "ui://my-widget",
    &mcpuiserver.RawHTMLPayload{
        Type:       mcpuiserver.ContentTypeRawHTML,
        HTMLString: "<h1>My Widget</h1>",
    },
    mcpuiserver.EncodingText,
    mcpuiserver.WithUIMetadata(map[string]interface{}{
        mcpuiserver.UIMetadataKeyInitialRenderData: renderData,
    }),
)
```

### Display Modes

Specify how widgets should be displayed:

```go
mcpuiserver.DisplayModeInline     // Widget displayed inline in chat
mcpuiserver.DisplayModePIP        // Picture-in-picture mode
mcpuiserver.DisplayModeFullscreen // Fullscreen mode
```

For more details on the protocol, see [docs/PROTOCOL.md](docs/PROTOCOL.md).

## Usage Examples

### Creating UI Resources

#### Raw HTML Resource

```go
resource, err := mcpuiserver.CreateUIResource(
    "ui://my-widget",
    &mcpuiserver.RawHTMLPayload{
        Type:       mcpuiserver.ContentTypeRawHTML,
        HTMLString: "<div><h1>My Widget</h1><p>Content here</p></div>",
    },
    mcpuiserver.EncodingText,
)
```

#### External URL Resource

```go
resource, err := mcpuiserver.CreateUIResource(
    "ui://dashboard",
    &mcpuiserver.ExternalURLPayload{
        Type:      mcpuiserver.ContentTypeExternalURL,
        IframeURL: "https://example.com/dashboard",
    },
    mcpuiserver.EncodingText,
)
```

#### Remote DOM Resource (React)

```go
resource, err := mcpuiserver.CreateUIResource(
    "ui://interactive-component",
    &mcpuiserver.RemoteDOMPayload{
        Type:      mcpuiserver.ContentTypeRemoteDOM,
        Script:    "console.log('React component');",
        Framework: mcpuiserver.FrameworkReact,
    },
    mcpuiserver.EncodingText,
)
```

#### Remote DOM Resource (WebComponents)

```go
resource, err := mcpuiserver.CreateUIResource(
    "ui://web-component",
    &mcpuiserver.RemoteDOMPayload{
        Type:      mcpuiserver.ContentTypeRemoteDOM,
        Script:    "console.log('Web component');",
        Framework: mcpuiserver.FrameworkWebComponents,
    },
    mcpuiserver.EncodingBlob, // Use base64 encoding
)
```

### Using Metadata

#### UI-Specific Metadata

UI metadata is automatically prefixed with `mcpui.dev/ui-`:

```go
resource, err := mcpuiserver.CreateUIResource(
    "ui://sized-widget",
    &mcpuiserver.RawHTMLPayload{
        Type:       mcpuiserver.ContentTypeRawHTML,
        HTMLString: "<h1>Sized Widget</h1>",
    },
    mcpuiserver.EncodingText,
    mcpuiserver.WithUIMetadata(map[string]interface{}{
        mcpuiserver.UIMetadataKeyPreferredFrameSize: []string{"800px", "600px"},
        mcpuiserver.UIMetadataKeyInitialRenderData: map[string]interface{}{
            "userId": "123",
            "theme":  "dark",
        },
    }),
)
```

#### Custom Metadata

```go
resource, err := mcpuiserver.CreateUIResource(
    "ui://my-widget",
    &mcpuiserver.RawHTMLPayload{
        Type:       mcpuiserver.ContentTypeRawHTML,
        HTMLString: "<h1>My Widget</h1>",
    },
    mcpuiserver.EncodingText,
    mcpuiserver.WithMetadata(map[string]interface{}{
        "custom.author":  "MyServer",
        "custom.version": "1.0.0",
    }),
)
```

#### Combining Metadata Types

```go
resource, err := mcpuiserver.CreateUIResource(
    "ui://full-widget",
    &mcpuiserver.RawHTMLPayload{
        Type:       mcpuiserver.ContentTypeRawHTML,
        HTMLString: "<h1>Full Widget</h1>",
    },
    mcpuiserver.EncodingText,
    mcpuiserver.WithUIMetadata(map[string]interface{}{
        mcpuiserver.UIMetadataKeyPreferredFrameSize: []string{"1200px", "800px"},
    }),
    mcpuiserver.WithMetadata(map[string]interface{}{
        "custom.category": "dashboard",
    }),
)
```

### UI Action Results

UI action results allow widgets to communicate actions back to the host.

#### Tool Call

```go
toolCall := mcpuiserver.UIActionResultToolCall("fetchData", map[string]interface{}{
    "query": "user stats",
    "limit": 100,
})
```

#### Prompt

```go
prompt := mcpuiserver.UIActionResultPrompt("Enter your API key")
```

#### Link

```go
link := mcpuiserver.UIActionResultLink("https://docs.example.com")
```

#### Intent

```go
intent := mcpuiserver.UIActionResultIntent("showSettings", map[string]interface{}{
    "tab": "account",
})
```

#### Notification

```go
notification := mcpuiserver.UIActionResultNotification("Data saved successfully!")
```

#### With Message ID

```go
toolCall := mcpuiserver.UIActionResultToolCall("fetchData", map[string]interface{}{
    "query": "user stats",
})

messageID := "msg-123"
toolCall.MessageID = &messageID
```

## API Reference

### Core Function

#### `CreateUIResource`

```go
func CreateUIResource(
    uri string,
    content ResourceContentPayload,
    encoding Encoding,
    opts ...Option,
) (*UIResource, error)
```

Creates a UI resource for inclusion in MCP tool results.

**Parameters:**
- `uri` - Resource identifier starting with `ui://`
- `content` - Content payload (RawHTMLPayload, ExternalURLPayload, or RemoteDOMPayload)
- `encoding` - Encoding type (EncodingText or EncodingBlob)
- `opts` - Optional functional options for metadata and properties

**Returns:**
- `*UIResource` - The created UI resource
- `error` - Validation or processing errors

### Content Payloads

#### `RawHTMLPayload`

```go
type RawHTMLPayload struct {
    Type       ContentType // ContentTypeRawHTML
    HTMLString string      // HTML content
}
```

#### `ExternalURLPayload`

```go
type ExternalURLPayload struct {
    Type      ContentType // ContentTypeExternalURL
    IframeURL string      // URL to display in iframe
}
```

#### `RemoteDOMPayload`

```go
type RemoteDOMPayload struct {
    Type      ContentType        // ContentTypeRemoteDOM
    Script    string             // JavaScript code
    Framework RemoteDOMFramework // FrameworkReact or FrameworkWebComponents
}
```

### Functional Options

#### `WithUIMetadata`

```go
func WithUIMetadata(metadata map[string]interface{}) Option
```

Sets UI-specific metadata (will be prefixed with `mcpui.dev/ui-`).

#### `WithMetadata`

```go
func WithMetadata(metadata map[string]interface{}) Option
```

Sets custom metadata (no prefix).

#### `WithResourceProps`

```go
func WithResourceProps(props map[string]interface{}) Option
```

Sets additional resource properties.

#### `WithEmbeddedResourceProps`

```go
func WithEmbeddedResourceProps(props map[string]interface{}) Option
```

Sets embedded resource properties (annotations, _meta).

### UI Action Result Constructors

- `UIActionResultToolCall(toolName string, params map[string]interface{}) UIActionResultToolCallType`
- `UIActionResultPrompt(prompt string) UIActionResultPromptType`
- `UIActionResultLink(url string) UIActionResultLinkType`
- `UIActionResultIntent(intent string, params map[string]interface{}) UIActionResultIntentType`
- `UIActionResultNotification(message string) UIActionResultNotificationType`

### Constants

#### Content Types

- `ContentTypeRawHTML` - Raw HTML content
- `ContentTypeExternalURL` - External URL
- `ContentTypeRemoteDOM` - Remote DOM component

#### Encoding Types

- `EncodingText` - Text encoding
- `EncodingBlob` - Base64 encoding

#### Frameworks

- `FrameworkReact` - React framework
- `FrameworkWebComponents` - Web Components framework

#### MIME Types

- `MimeTypeHTML` - `text/html`
- `MimeTypeURIList` - `text/uri-list`
- `MimeTypeRemoteDomReact` - `application/vnd.mcp-ui.remote-dom+javascript; framework=react`
- `MimeTypeRemoteDomWC` - `application/vnd.mcp-ui.remote-dom+javascript; framework=webcomponents`
- `MimeTypeMCPAppsAdapter` - `text/html;profile=mcp-app` (MCP Apps standard)
- `MimeTypeAppsSdkAdapter` - `text/html+skybridge` (Apps SDK)

#### Protocol Constants

- `ProtocolVersion` - `2025-11-21` (MCP Apps SEP protocol version)
- `ResourceURIMetaKey` - `ui/resourceUri` (metadata key for resource URI)

#### Message Type Constants

- `MessageTypeToolCall` - `tool`
- `MessageTypePrompt` - `prompt`
- `MessageTypeLink` - `link`
- `MessageTypeIntent` - `intent`
- `MessageTypeNotify` - `notify`
- `MessageTypeLifecycleReady` - `ui-lifecycle-iframe-ready`
- `MessageTypeSizeChange` - `ui-size-change`
- `MessageTypeRequestData` - `ui-request-data`
- `MessageTypeRequestRenderData` - `ui-request-render-data`
- `MessageTypeLifecycleRenderData` - `ui-lifecycle-iframe-render-data`
- `MessageTypeMessageReceived` - `ui-message-received`
- `MessageTypeMessageResponse` - `ui-message-response`

#### Display Mode Constants

- `DisplayModeInline` - `inline`
- `DisplayModePIP` - `pip`
- `DisplayModeFullscreen` - `fullscreen`

#### Metadata Keys

- `UIMetadataKeyPreferredFrameSize` - `preferred-frame-size`
- `UIMetadataKeyInitialRenderData` - `initial-render-data`

#### Prefixes

- `UIMetadataPrefix` - `mcpui.dev/ui-`
- `URIScheme` - `ui://`

### Error Types

- `ErrInvalidURI` - URI doesn't start with `ui://`
- `ErrEmptyHTMLString` - HTML string is empty
- `ErrEmptyIframeURL` - Iframe URL is empty
- `ErrEmptyScript` - Script is empty
- `ErrInvalidFramework` - Framework is not 'react' or 'webcomponents'
- `ErrInvalidEncoding` - Encoding is not 'text' or 'blob'
- `ErrNilContent` - Content is nil

## Error Handling

All errors are strongly typed and can be checked using `errors.Is`:

```go
resource, err := mcpuiserver.CreateUIResource(
    "http://invalid",
    &mcpuiserver.RawHTMLPayload{
        Type:       mcpuiserver.ContentTypeRawHTML,
        HTMLString: "<p>Test</p>",
    },
    mcpuiserver.EncodingText,
)

if errors.Is(err, mcpuiserver.ErrInvalidURI) {
    fmt.Println("Invalid URI scheme")
}
```

## Testing

Run tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

Run specific test:

```bash
go test -v -run TestCreateUIResource
```

## Migration from v1.x

### MIME Type Change (Breaking)

The MCP Apps adapter now uses the official MIME type `text/html;profile=mcp-app` instead of plain `text/html`.

**Before:**
```go
// Implicit MIME type was "text/html"
adapter.GetMIMEType() // Returns "text/html"
```

**After:**
```go
// Now returns the official MCP Apps MIME type
adapter.GetMIMEType() // Returns "text/html;profile=mcp-app"
```

**Impact:** This is a breaking change for the MIME type value, but should be transparent for most users. MCP Apps hosts expect the `profile=mcp-app` variant and will work correctly with this change.

### New Constants and Types

The SDK now exports several new constants and types for protocol alignment:

- `ProtocolVersion`: Protocol version constant ("2025-11-21")
- `ResourceURIMetaKey`: Metadata key for resource URI ("ui/resourceUri")
- `RenderData`: Type-safe render data structure
- `DisplayMode`: Display mode constants (inline, pip, fullscreen)
- Protocol message types and structures

These additions are backward compatible - existing code will continue to work.

## License

Apache-2.0

## Links

- **Homepage:** https://mcpui.dev
- **Repository:** https://github.com/MCP-UI-Org/mcp-ui
- **Documentation:** https://mcpui.dev/guide/introduction
- **Go Package Documentation:** https://pkg.go.dev/github.com/MCP-UI-Org/mcp-ui/sdks/go/server

## Contributing

Contributions are welcome! Please see the main repository for contribution guidelines.

## Support

For issues and questions:
- Open an issue on GitHub: https://github.com/MCP-UI-Org/mcp-ui/issues
- Visit the documentation: https://mcpui.dev
