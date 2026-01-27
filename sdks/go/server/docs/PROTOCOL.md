# MCP-UI Protocol Reference

This document describes the MCP-UI protocol messages supported by the Go SDK.

## Protocol Version

Current protocol version: `2025-11-21`

Access via: `mcpuiserver.ProtocolVersion`

## Message Types

### UI Action Messages

These messages are sent from widgets to the host to request actions.

#### Tool Call

Request execution of an MCP tool.

```go
toolCall := mcpuiserver.UIActionResultToolCall("myTool", map[string]interface{}{
    "param": "value",
})
```

**Message Type:** `tool`

#### Prompt

Request user input via prompt.

```go
prompt := mcpuiserver.UIActionResultPrompt("Enter API key")
```

**Message Type:** `prompt`

#### Link

Open a link in the host.

```go
link := mcpuiserver.UIActionResultLink("https://example.com")
```

**Message Type:** `link`

#### Intent

Request a host-specific action.

```go
intent := mcpuiserver.UIActionResultIntent("showSettings", map[string]interface{}{
    "tab": "account",
})
```

**Message Type:** `intent`

#### Notification

Display a notification to the user.

```go
notification := mcpuiserver.UIActionResultNotification("Data saved!")
```

**Message Type:** `notify`

### Lifecycle Messages

#### Ready

Indicates the widget has loaded and is ready.

```go
msgID := "msg-123"
readyMsg := mcpuiserver.NewLifecycleReadyMessage(&msgID)
```

**Message Type:** `ui-lifecycle-iframe-ready`

**Structure:**
```go
type MCPUILifecycleReadyMessage struct {
    Type      ProtocolMessageType    // "ui-lifecycle-iframe-ready"
    MessageID *string                // Optional message ID
    Payload   map[string]interface{} // Additional payload data
}
```

#### Size Change

Request a change in the widget's display size.

```go
width := 800
height := 600
msgID := "msg-456"
sizeMsg := mcpuiserver.NewSizeChangeMessage(&width, &height, &msgID)
```

**Message Type:** `ui-size-change`

**Structure:**
```go
type MCPUISizeChangeMessage struct {
    Type      ProtocolMessageType // "ui-size-change"
    MessageID *string             // Optional message ID
    Payload   SizeChangePayload
}

type SizeChangePayload struct {
    Width  *int // Width in pixels (optional)
    Height *int // Height in pixels (optional)
}
```

### Data Messages

#### Request Data

Request data from the host.

```go
requestMsg := mcpuiserver.NewRequestDataMessage(
    "getUserData",
    map[string]interface{}{"userId": "123"},
    "msg-789",
)
```

**Message Type:** `ui-request-data`

**Structure:**
```go
type MCPUIRequestDataMessage struct {
    Type      ProtocolMessageType // "ui-request-data"
    MessageID string              // Required message ID
    Payload   RequestDataPayload
}

type RequestDataPayload struct {
    RequestType string                 // Type of data request
    Params      map[string]interface{} // Request parameters
}
```

#### Request Render Data

Request render data from the host.

```go
renderDataRequestMsg := mcpuiserver.NewRequestRenderDataMessage(&msgID)
```

**Message Type:** `ui-request-render-data`

**Structure:**
```go
type MCPUIRequestRenderDataMessage struct {
    Type      ProtocolMessageType    // "ui-request-render-data"
    MessageID *string                // Optional message ID
    Payload   map[string]interface{} // Additional payload data
}
```

#### Render Data Message

Deliver render data to the widget (host → widget).

```go
renderData := mcpuiserver.RenderData{
    Locale:      "en-US",
    Theme:       "dark",
    DisplayMode: mcpuiserver.DisplayModeInline,
    MaxHeight:   600,
}
renderDataMsg := mcpuiserver.NewRenderDataMessage(renderData, &msgID)
```

**Message Type:** `ui-lifecycle-iframe-render-data`

**Structure:**
```go
type MCPUIRenderDataMessage struct {
    Type      ProtocolMessageType // "ui-lifecycle-iframe-render-data"
    MessageID *string             // Optional message ID
    Payload   RenderDataPayload
}

type RenderDataPayload struct {
    RenderData RenderData
}
```

#### Message Received

Acknowledge receipt of a message.

```go
ackMsg := mcpuiserver.NewMessageReceivedMessage("msg-original", &msgID)
```

**Message Type:** `ui-message-received`

**Structure:**
```go
type MCPUIMessageReceivedMessage struct {
    Type      ProtocolMessageType    // "ui-message-received"
    MessageID *string                // Optional message ID
    Payload   MessageReceivedPayload
}

type MessageReceivedPayload struct {
    MessageID string // ID of the acknowledged message
}
```

#### Message Response

Deliver a response to a request.

```go
responseMsg := mcpuiserver.NewMessageResponseMessage(
    "msg-request",
    map[string]interface{}{"status": "ok"},
    nil,
    &msgID,
)
```

**Message Type:** `ui-message-response`

**Structure:**
```go
type MCPUIMessageResponseMessage struct {
    Type      ProtocolMessageType    // "ui-message-response"
    MessageID *string                // Optional message ID
    Payload   MessageResponsePayload
}

type MessageResponsePayload struct {
    MessageID string      // ID of the request message
    Response  interface{} // Response data (optional)
    Error     interface{} // Error data (optional)
}
```

## Render Data

Render data provides context and initialization information to widgets.

### Structure

```go
type RenderData struct {
    ToolInput   map[string]interface{} // Input parameters from tool invocation
    ToolOutput  interface{}            // Output from tool execution
    WidgetState interface{}            // Persistent widget state
    Locale      string                 // User's locale (e.g., "en-US")
    Theme       string                 // UI theme (e.g., "dark", "light")
    DisplayMode DisplayMode            // Display mode (inline, pip, fullscreen)
    MaxHeight   int                    // Maximum height in pixels
}
```

### Display Modes

- `DisplayModeInline`: Widget is displayed inline in the chat
- `DisplayModePIP`: Widget is displayed in picture-in-picture mode
- `DisplayModeFullscreen`: Widget is displayed in fullscreen mode

### Usage

Pass render data when creating a resource:

```go
renderData := mcpuiserver.RenderData{
    Locale:      "en-US",
    Theme:       "dark",
    DisplayMode: mcpuiserver.DisplayModeInline,
    MaxHeight:   600,
    ToolInput: map[string]interface{}{
        "query": "search term",
    },
}

resource, err := mcpuiserver.CreateUIResource(
    "ui://my-widget",
    &mcpuiserver.RawHTMLPayload{
        Type:       mcpuiserver.ContentTypeRawHTML,
        HTMLString: htmlContent,
    },
    mcpuiserver.EncodingText,
    mcpuiserver.WithUIMetadata(map[string]interface{}{
        mcpuiserver.UIMetadataKeyInitialRenderData: renderData,
    }),
)
```

## MCP Apps Standard

### Resource URI Metadata

MCP Apps hosts use a special metadata key to identify UI resources in tool responses:

```go
mcpuiserver.ResourceURIMetaKey // "ui/resourceUri"
```

When returning a UI resource in a tool response, include the resource URI in the tool's metadata:

```go
// In your MCP tool handler
toolResult := map[string]interface{}{
    "content": []interface{}{
        widgetResource, // Your UI resource
    },
    "_meta": map[string]interface{}{
        mcpuiserver.ResourceURIMetaKey: widgetResource.Resource.URI,
    },
}
```

### MIME Type

MCP Apps resources must use the official MIME type:

```go
"text/html;profile=mcp-app"
```

This is automatically set when using the MCP Apps adapter:

```go
resource, err := mcpuiserver.CreateUIResource(
    "ui://my-widget",
    payload,
    encoding,
    mcpuiserver.WithProtocol(mcpuiserver.ProtocolTypeMCPApps),
)
// resource.Resource.MimeType will be "text/html;profile=mcp-app"
```

## Implementation Examples

### Complete Widget with Render Data

```go
package main

import (
    "encoding/json"
    "log"

    mcpuiserver "github.com/MCP-UI-Org/mcp-ui/sdks/go/server"
)

func createDashboardWidget(userID string) (*mcpuiserver.UIResource, error) {
    // Prepare render data
    renderData := mcpuiserver.RenderData{
        Locale:      "en-US",
        Theme:       "dark",
        DisplayMode: mcpuiserver.DisplayModeInline,
        MaxHeight:   800,
        ToolInput: map[string]interface{}{
            "userId": userID,
        },
    }

    // Create widget HTML
    html := `<!DOCTYPE html>
<html>
<head>
    <title>Dashboard</title>
</head>
<body>
    <div id="app">Loading...</div>
    <script>
        // Widget initialization code
        window.addEventListener('message', function(event) {
            if (event.data.type === 'ui-lifecycle-iframe-render-data') {
                const renderData = event.data.payload.renderData;
                console.log('Received render data:', renderData);
                initializeDashboard(renderData);
            }
        });

        // Signal ready
        window.parent.postMessage({
            type: 'ui-lifecycle-iframe-ready'
        }, '*');
    </script>
</body>
</html>`

    // Create resource with MCP Apps adapter
    return mcpuiserver.CreateUIResource(
        "ui://dashboard",
        &mcpuiserver.RawHTMLPayload{
            Type:       mcpuiserver.ContentTypeRawHTML,
            HTMLString: html,
        },
        mcpuiserver.EncodingText,
        mcpuiserver.WithProtocol(mcpuiserver.ProtocolTypeMCPApps),
        mcpuiserver.WithUIMetadata(map[string]interface{}{
            mcpuiserver.UIMetadataKeyInitialRenderData: renderData,
            mcpuiserver.UIMetadataKeyPreferredFrameSize: []string{"1200px", "800px"},
        }),
    )
}

// Usage in MCP tool handler
func handleDashboardTool(userID string) (map[string]interface{}, error) {
    // Create the widget
    widget, err := createDashboardWidget(userID)
    if err != nil {
        return nil, err
    }

    // Return tool result with resource URI in metadata
    return map[string]interface{}{
        "content": []interface{}{
            map[string]interface{}{
                "type": "text",
                "text": "Here's your dashboard:",
            },
            widget,
        },
        "_meta": map[string]interface{}{
            mcpuiserver.ResourceURIMetaKey: widget.Resource.URI,
        },
    }, nil
}

func main() {
    result, err := handleDashboardTool("user-123")
    if err != nil {
        log.Fatal(err)
    }

    jsonBytes, _ := json.MarshalIndent(result, "", "  ")
    log.Println(string(jsonBytes))
}
```

### Widget-Side Protocol Implementation

Example HTML/JavaScript for a widget that communicates with the host:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Interactive Widget</title>
    <style>
        body { font-family: Arial, sans-serif; padding: 20px; }
        button { padding: 10px 20px; margin: 5px; }
    </style>
</head>
<body>
    <div id="app">
        <h1>Interactive Widget</h1>
        <div id="data"></div>
        <button onclick="callTool()">Call Tool</button>
        <button onclick="sendPrompt()">Send Prompt</button>
        <button onclick="resize()">Resize</button>
    </div>

    <script>
        // Listen for messages from host
        window.addEventListener('message', function(event) {
            if (event.data.type === 'ui-lifecycle-iframe-render-data') {
                const renderData = event.data.payload.renderData;
                displayData(renderData);
            } else if (event.data.type === 'ui-message-response') {
                console.log('Received response:', event.data.payload);
            }
        });

        // Display render data
        function displayData(renderData) {
            const dataDiv = document.getElementById('data');
            dataDiv.innerHTML = '<pre>' + JSON.stringify(renderData, null, 2) + '</pre>';
        }

        // Call a tool
        function callTool() {
            window.parent.postMessage({
                type: 'tool',
                messageId: 'msg-' + Date.now(),
                payload: {
                    toolName: 'fetchData',
                    params: { query: 'test' }
                }
            }, '*');
        }

        // Send a prompt
        function sendPrompt() {
            window.parent.postMessage({
                type: 'prompt',
                messageId: 'msg-' + Date.now(),
                payload: { prompt: 'Enter your name' }
            }, '*');
        }

        // Request size change
        function resize() {
            window.parent.postMessage({
                type: 'ui-size-change',
                messageId: 'msg-' + Date.now(),
                payload: { width: 1000, height: 800 }
            }, '*');
        }

        // Signal ready on load
        window.parent.postMessage({
            type: 'ui-lifecycle-iframe-ready',
            messageId: 'msg-ready'
        }, '*');

        // Request render data
        window.parent.postMessage({
            type: 'ui-request-render-data',
            messageId: 'msg-render-data'
        }, '*');
    </script>
</body>
</html>
```

## Protocol Message Flow

### Initialization Sequence

1. **Widget loads** - HTML/JavaScript loads in iframe
2. **Widget signals ready** - Sends `ui-lifecycle-iframe-ready` message
3. **Widget requests render data** - Sends `ui-request-render-data` message
4. **Host delivers render data** - Sends `ui-lifecycle-iframe-render-data` message
5. **Widget initializes** - Uses render data to set up UI

### Tool Call Sequence

1. **Widget calls tool** - Sends `tool` message with tool name and params
2. **Host acknowledges** - Sends `ui-message-received` message
3. **Host executes tool** - Calls the MCP tool handler
4. **Host responds** - Sends `ui-message-response` message with result or error

### Size Change Sequence

1. **Widget requests resize** - Sends `ui-size-change` message with dimensions
2. **Host acknowledges** - Sends `ui-message-received` message
3. **Host resizes container** - Updates iframe/container size

## Message ID Conventions

- Message IDs are optional for most message types
- Use format `msg-{timestamp}` or `msg-{unique-identifier}`
- Include message IDs for request/response correlation
- Hosts use message IDs to track pending requests and responses

## Error Handling

When errors occur, the host sends a `ui-message-response` with an error payload:

```go
errorMsg := mcpuiserver.NewMessageResponseMessage(
    "msg-request-123",
    nil,
    map[string]interface{}{
        "code":    "ERR_TOOL_NOT_FOUND",
        "message": "Tool 'invalidTool' not found",
    },
    &msgID,
)
```

Widget-side error handling:

```javascript
window.addEventListener('message', function(event) {
    if (event.data.type === 'ui-message-response') {
        const payload = event.data.payload;
        if (payload.error) {
            console.error('Error:', payload.error);
            // Handle error
        } else {
            console.log('Response:', payload.response);
            // Handle success
        }
    }
});
```

## Best Practices

1. **Always signal ready** - Send `ui-lifecycle-iframe-ready` when widget loads
2. **Request render data** - Get initial context with `ui-request-render-data`
3. **Use message IDs** - Include message IDs for request/response tracking
4. **Handle errors** - Check for error payload in responses
5. **Validate payloads** - Validate incoming data before using
6. **Set display mode** - Specify appropriate display mode in render data
7. **Include resource URI** - Add `ResourceURIMetaKey` to tool metadata
8. **Use standard MIME type** - Use `text/html;profile=mcp-app` for MCP Apps

## See Also

- [Main README](../README.md) - SDK overview and examples
- [MCP Apps Specification](https://github.com/modelcontextprotocol/specification/blob/main/extensions/apps.md) - Official MCP Apps spec
- [Examples](../examples/) - Complete working examples
