package main

import (
	"encoding/json"
	"fmt"
	"log"

	mcpuiserver "github.com/MCP-UI-Org/mcp-ui/sdks/go/server"
)

func main() {
	fmt.Println("MCP-UI Server SDK for Go - Examples")
	fmt.Println("====================================")
	fmt.Println()

	// Example 1: Simple HTML Resource
	example1()

	// Example 2: External URL Resource
	example2()

	// Example 3: Remote DOM Resource with React
	example3()

	// Example 4: Resource with Metadata
	example4()

	// Example 5: Blob-encoded Resource
	example5()

	// Example 6: UI Action Results
	example6()
}

func example1() {
	fmt.Println("Example 1: Simple HTML Resource")
	fmt.Println("--------------------------------")

	resource, err := mcpuiserver.CreateUIResource(
		"ui://greeting",
		&mcpuiserver.RawHTMLPayload{
			Type:       mcpuiserver.ContentTypeRawHTML,
			HTMLString: "<h1>Hello, World!</h1><p>This is a simple HTML widget.</p>",
		},
		mcpuiserver.EncodingText,
	)

	if err != nil {
		log.Fatal(err)
	}

	printJSON(resource)
}

func example2() {
	fmt.Println("Example 2: External URL Resource")
	fmt.Println("---------------------------------")

	resource, err := mcpuiserver.CreateUIResource(
		"ui://dashboard",
		&mcpuiserver.ExternalURLPayload{
			Type:      mcpuiserver.ContentTypeExternalURL,
			IframeURL: "https://example.com/dashboard",
		},
		mcpuiserver.EncodingText,
	)

	if err != nil {
		log.Fatal(err)
	}

	printJSON(resource)
}

func example3() {
	fmt.Println("Example 3: Remote DOM Resource with React")
	fmt.Println("------------------------------------------")

	scriptContent := `
import React from 'react';
import { createRoot } from 'react-dom/client';

const App = () => {
  return (
    <div>
      <h1>React Component</h1>
      <p>This is a React-based UI widget.</p>
    </div>
  );
};

const root = createRoot(document.getElementById('root'));
root.render(<App />);
`

	resource, err := mcpuiserver.CreateUIResource(
		"ui://react-widget",
		&mcpuiserver.RemoteDOMPayload{
			Type:      mcpuiserver.ContentTypeRemoteDOM,
			Script:    scriptContent,
			Framework: mcpuiserver.FrameworkReact,
		},
		mcpuiserver.EncodingText,
	)

	if err != nil {
		log.Fatal(err)
	}

	printJSON(resource)
}

func example4() {
	fmt.Println("Example 4: Resource with Metadata")
	fmt.Println("----------------------------------")

	resource, err := mcpuiserver.CreateUIResource(
		"ui://sized-widget",
		&mcpuiserver.RawHTMLPayload{
			Type:       mcpuiserver.ContentTypeRawHTML,
			HTMLString: "<h1>Sized Widget</h1><p>This widget has specific size preferences.</p>",
		},
		mcpuiserver.EncodingText,
		mcpuiserver.WithUIMetadata(map[string]interface{}{
			mcpuiserver.UIMetadataKeyPreferredFrameSize: []string{"800px", "600px"},
			mcpuiserver.UIMetadataKeyInitialRenderData: map[string]interface{}{
				"userId": "12345",
				"theme":  "dark",
			},
		}),
		mcpuiserver.WithMetadata(map[string]interface{}{
			"custom.author":  "Example Server",
			"custom.version": "1.0.0",
		}),
	)

	if err != nil {
		log.Fatal(err)
	}

	printJSON(resource)
}

func example5() {
	fmt.Println("Example 5: Blob-encoded Resource")
	fmt.Println("---------------------------------")

	largeHTML := `
<!DOCTYPE html>
<html>
<head>
    <title>Large Widget</title>
    <style>
        body { font-family: Arial, sans-serif; padding: 20px; }
        h1 { color: #333; }
    </style>
</head>
<body>
    <h1>Large HTML Widget</h1>
    <p>This content is base64-encoded for efficient transmission.</p>
</body>
</html>
`

	resource, err := mcpuiserver.CreateUIResource(
		"ui://large-widget",
		&mcpuiserver.RawHTMLPayload{
			Type:       mcpuiserver.ContentTypeRawHTML,
			HTMLString: largeHTML,
		},
		mcpuiserver.EncodingBlob, // Base64 encoding
	)

	if err != nil {
		log.Fatal(err)
	}

	printJSON(resource)
}

func example6() {
	fmt.Println("Example 6: UI Action Results")
	fmt.Println("-----------------------------")

	// Tool call
	toolCall := mcpuiserver.UIActionResultToolCall("fetchData", map[string]interface{}{
		"query": "user stats",
		"limit": 100,
	})
	fmt.Println("\nTool Call:")
	printJSON(toolCall)

	// Prompt
	prompt := mcpuiserver.UIActionResultPrompt("Enter your API key")
	fmt.Println("\nPrompt:")
	printJSON(prompt)

	// Link
	link := mcpuiserver.UIActionResultLink("https://docs.example.com")
	fmt.Println("\nLink:")
	printJSON(link)

	// Intent
	intent := mcpuiserver.UIActionResultIntent("showSettings", map[string]interface{}{
		"tab": "account",
	})
	fmt.Println("\nIntent:")
	printJSON(intent)

	// Notification
	notification := mcpuiserver.UIActionResultNotification("Data saved successfully!")
	fmt.Println("\nNotification:")
	printJSON(notification)
}

func printJSON(v interface{}) {
	jsonBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonBytes))
	fmt.Println()
}
