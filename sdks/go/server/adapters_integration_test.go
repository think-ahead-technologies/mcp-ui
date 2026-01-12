package mcpuiserver

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/MCP-UI-Org/mcp-ui/sdks/go/server/adapters/appssdk"
	"github.com/MCP-UI-Org/mcp-ui/sdks/go/server/adapters/mcpapps"
	"github.com/stretchr/testify/assert"
)

func TestCreateUIResource_WithAppsSdkAdapter(t *testing.T) {
	adapter, err := appssdk.NewAdapter()
	assert.NoError(t, err)

	resource, err := CreateUIResource(
		"ui://test-widget",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<h1>Hello World</h1>",
		},
		EncodingText,
		WithAdapter(adapter),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)

	// Verify MIME type is overridden by adapter
	assert.Equal(t, "text/html+skybridge", resource.Resource.MimeType)

	// Verify content contains adapter script
	assert.Contains(t, resource.Resource.Text, "<script>")
	assert.Contains(t, resource.Resource.Text, "MCPUIAppsSdkAdapter")
	assert.Contains(t, resource.Resource.Text, "<h1>Hello World</h1>")
}

func TestCreateUIResource_WithMcpAppsAdapter(t *testing.T) {
	adapter, err := mcpapps.NewAdapter()
	assert.NoError(t, err)

	resource, err := CreateUIResource(
		"ui://test-widget",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<h1>Hello World</h1>",
		},
		EncodingText,
		WithAdapter(adapter),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)

	// Verify MIME type
	assert.Equal(t, "text/html", resource.Resource.MimeType)

	// Verify content contains adapter script
	assert.Contains(t, resource.Resource.Text, "<script>")
	assert.Contains(t, resource.Resource.Text, "McpAppsAdapter")
	assert.Contains(t, resource.Resource.Text, "<h1>Hello World</h1>")
}

func TestCreateUIResource_WithAdapter_BlobEncoding(t *testing.T) {
	adapter, err := appssdk.NewAdapter()
	assert.NoError(t, err)

	resource, err := CreateUIResource(
		"ui://test-widget",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<h1>Hello World</h1>",
		},
		EncodingBlob,
		WithAdapter(adapter),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)

	// Verify MIME type
	assert.Equal(t, "text/html+skybridge", resource.Resource.MimeType)

	// Verify blob encoding (content should be base64)
	assert.NotEmpty(t, resource.Resource.Blob)
	assert.Empty(t, resource.Resource.Text)

	// Decode and verify content
	decoded, err := base64.StdEncoding.DecodeString(resource.Resource.Blob)
	assert.NoError(t, err)
	decodedStr := string(decoded)
	assert.Contains(t, decodedStr, "<script>")
	assert.Contains(t, decodedStr, "MCPUIAppsSdkAdapter")
	assert.Contains(t, decodedStr, "<h1>Hello World</h1>")
}

func TestCreateUIResource_WithoutAdapter(t *testing.T) {
	// Verify existing behavior is preserved when no adapter is provided
	resource, err := CreateUIResource(
		"ui://test-widget",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<h1>Hello World</h1>",
		},
		EncodingText,
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)

	// Verify MIME type is standard HTML
	assert.Equal(t, "text/html", resource.Resource.MimeType)

	// Verify content does NOT contain adapter script
	assert.NotContains(t, resource.Resource.Text, "MCPUIAppsSdkAdapter")
	assert.NotContains(t, resource.Resource.Text, "McpAppsAdapter")
	assert.Equal(t, "<h1>Hello World</h1>", resource.Resource.Text)
}

func TestCreateUIResource_Adapter_OnlyWorksWithRawHTML(t *testing.T) {
	adapter, err := appssdk.NewAdapter()
	assert.NoError(t, err)

	t.Run("adapter ignored for ExternalURL", func(t *testing.T) {
		resource, err := CreateUIResource(
			"ui://test-widget",
			&ExternalURLPayload{
				Type:      ContentTypeExternalURL,
				IframeURL: "https://example.com",
			},
			EncodingText,
			WithAdapter(adapter),
		)

		assert.NoError(t, err)
		assert.NotNil(t, resource)

		// Adapter should not be applied
		assert.Equal(t, "text/uri-list", resource.Resource.MimeType)
		assert.Equal(t, "https://example.com", resource.Resource.Text)
		assert.NotContains(t, resource.Resource.Text, "MCPUIAppsSdkAdapter")
	})

	t.Run("adapter ignored for RemoteDOM", func(t *testing.T) {
		resource, err := CreateUIResource(
			"ui://test-widget",
			&RemoteDOMPayload{
				Type:      ContentTypeRemoteDOM,
				Script:    "console.log('test')",
				Framework: FrameworkReact,
			},
			EncodingText,
			WithAdapter(adapter),
		)

		assert.NoError(t, err)
		assert.NotNil(t, resource)

		// Adapter should not be applied
		assert.Equal(t, MimeTypeRemoteDomReact, resource.Resource.MimeType)
		assert.Equal(t, "console.log('test')", resource.Resource.Text)
		assert.NotContains(t, resource.Resource.Text, "MCPUIAppsSdkAdapter")
	})
}

func TestCreateUIResource_Adapter_HTMLWithHead(t *testing.T) {
	adapter, err := appssdk.NewAdapter()
	assert.NoError(t, err)

	htmlWithHead := `<html>
<head>
	<title>Test</title>
</head>
<body>
	<h1>Hello</h1>
</body>
</html>`

	resource, err := CreateUIResource(
		"ui://test-widget",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: htmlWithHead,
		},
		EncodingText,
		WithAdapter(adapter),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)

	// Verify adapter script is injected into existing head
	assert.Contains(t, resource.Resource.Text, "<head>")
	assert.Contains(t, resource.Resource.Text, "<script>")
	assert.Contains(t, resource.Resource.Text, "MCPUIAppsSdkAdapter")
	assert.Contains(t, resource.Resource.Text, "<title>Test</title>")
	assert.Contains(t, resource.Resource.Text, "<h1>Hello</h1>")

	// Script should be after <head> but before other head content
	headIdx := strings.Index(resource.Resource.Text, "<head>")
	scriptIdx := strings.Index(resource.Resource.Text, "<script>")
	titleIdx := strings.Index(resource.Resource.Text, "<title>")
	assert.True(t, headIdx < scriptIdx, "script should be after <head>")
	assert.True(t, scriptIdx < titleIdx, "script should be before <title>")
}

func TestCreateUIResource_Adapter_HTMLWithoutHead(t *testing.T) {
	adapter, err := appssdk.NewAdapter()
	assert.NoError(t, err)

	htmlWithoutHead := `<html>
<body>
	<h1>Hello</h1>
</body>
</html>`

	resource, err := CreateUIResource(
		"ui://test-widget",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: htmlWithoutHead,
		},
		EncodingText,
		WithAdapter(adapter),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)

	// Verify head tag is created and adapter script is injected
	assert.Contains(t, resource.Resource.Text, "<html>")
	assert.Contains(t, resource.Resource.Text, "<head>")
	assert.Contains(t, resource.Resource.Text, "</head>")
	assert.Contains(t, resource.Resource.Text, "<script>")
	assert.Contains(t, resource.Resource.Text, "MCPUIAppsSdkAdapter")
	assert.Contains(t, resource.Resource.Text, "<h1>Hello</h1>")

	// Head should be created right after <html>
	htmlIdx := strings.Index(resource.Resource.Text, "<html>")
	headIdx := strings.Index(resource.Resource.Text, "<head>")
	assert.True(t, htmlIdx < headIdx, "<head> should be after <html>")
}

func TestCreateUIResource_Adapter_HTMLFragmentNoTags(t *testing.T) {
	adapter, err := appssdk.NewAdapter()
	assert.NoError(t, err)

	htmlFragment := `<h1>Hello</h1><p>World</p>`

	resource, err := CreateUIResource(
		"ui://test-widget",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: htmlFragment,
		},
		EncodingText,
		WithAdapter(adapter),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)

	// Verify full HTML structure is created
	assert.Contains(t, resource.Resource.Text, "<html>")
	assert.Contains(t, resource.Resource.Text, "<head>")
	assert.Contains(t, resource.Resource.Text, "</head>")
	assert.Contains(t, resource.Resource.Text, "<body>")
	assert.Contains(t, resource.Resource.Text, "</body>")
	assert.Contains(t, resource.Resource.Text, "</html>")
	assert.Contains(t, resource.Resource.Text, "<script>")
	assert.Contains(t, resource.Resource.Text, "MCPUIAppsSdkAdapter")
	assert.Contains(t, resource.Resource.Text, "<h1>Hello</h1>")
	assert.Contains(t, resource.Resource.Text, "<p>World</p>")
}

func TestCreateUIResource_Adapter_WithCustomConfig(t *testing.T) {
	adapter, err := appssdk.NewAdapter(
		appssdk.WithTimeout(5000),
		appssdk.WithIntentHandling("ignore"),
	)
	assert.NoError(t, err)

	resource, err := CreateUIResource(
		"ui://test-widget",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<h1>Hello</h1>",
		},
		EncodingText,
		WithAdapter(adapter),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)

	// Verify custom config is in the script
	assert.Contains(t, resource.Resource.Text, `"timeout":5000`)
	assert.Contains(t, resource.Resource.Text, `"intentHandling":"ignore"`)
}

func TestCreateUIResource_BothAdapters(t *testing.T) {
	t.Run("Apps SDK adapter", func(t *testing.T) {
		adapter, err := appssdk.NewAdapter()
		assert.NoError(t, err)

		resource, err := CreateUIResource(
			"ui://widget-apps-sdk",
			&RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<h1>Apps SDK Widget</h1>",
			},
			EncodingText,
			WithAdapter(adapter),
		)

		assert.NoError(t, err)
		assert.Equal(t, "text/html+skybridge", resource.Resource.MimeType)
		assert.Contains(t, resource.Resource.Text, "MCPUIAppsSdkAdapter")
	})

	t.Run("MCP Apps adapter", func(t *testing.T) {
		adapter, err := mcpapps.NewAdapter()
		assert.NoError(t, err)

		resource, err := CreateUIResource(
			"ui://widget-mcp-apps",
			&RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<h1>MCP Apps Widget</h1>",
			},
			EncodingText,
			WithAdapter(adapter),
		)

		assert.NoError(t, err)
		assert.Equal(t, "text/html", resource.Resource.MimeType)
		assert.Contains(t, resource.Resource.Text, "McpAppsAdapter")
	})
}
