package mcpuiserver

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUIResource_WithProtocol(t *testing.T) {
	tests := []struct {
		name         string
		protocol     ProtocolType
		expectedMIME string
		hasScript    bool
	}{
		{
			name:         "generic protocol",
			protocol:     ProtocolTypeGeneric,
			expectedMIME: MimeTypeHTML,
			hasScript:    false,
		},
		{
			name:         "Apps SDK protocol",
			protocol:     ProtocolTypeAppsSDK,
			expectedMIME: MimeTypeAppsSdkAdapter,
			hasScript:    true,
		},
		{
			name:         "MCP Apps protocol",
			protocol:     ProtocolTypeMCPApps,
			expectedMIME: MimeTypeMCPAppsAdapter,
			hasScript:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource, err := CreateUIResource(
				"ui://test",
				&RawHTMLPayload{
					Type:       ContentTypeRawHTML,
					HTMLString: "<h1>Test</h1>",
				},
				EncodingText,
				WithProtocol(tt.protocol),
			)

			assert.NoError(t, err)
			assert.NotNil(t, resource)
			assert.Equal(t, tt.expectedMIME, resource.Resource.MimeType)

			if tt.hasScript {
				assert.Contains(t, resource.Resource.Text, "<script")
			}
		})
	}
}

func TestCreateUIResource_WithProtocolConfig(t *testing.T) {
	config := &ProtocolConfig{
		Type:    ProtocolTypeMCPApps,
		Version: "v2",
		BaseURL: "https://custom.cdn.com",
		Config: map[string]interface{}{
			"timeout": 5000,
		},
	}

	resource, err := CreateUIResource(
		"ui://test",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<h1>Test</h1>",
		},
		EncodingText,
		WithProtocolConfig(config),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)
	assert.Equal(t, MimeTypeMCPAppsAdapter, resource.Resource.MimeType)
	assert.Contains(t, resource.Resource.Text, "https://custom.cdn.com/mcpapps-v2.js")
}

func TestCreateUIResource_WithProtocolVersion(t *testing.T) {
	resource, err := CreateUIResource(
		"ui://test",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<h1>Test</h1>",
		},
		EncodingText,
		WithProtocol(ProtocolTypeMCPApps),
		WithProtocolVersion("v3"),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)
	assert.Contains(t, resource.Resource.Text, "mcpapps-v3.js")
}

func TestCreateUIResource_WithProtocolBaseURL(t *testing.T) {
	resource, err := CreateUIResource(
		"ui://test",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<h1>Test</h1>",
		},
		EncodingText,
		WithProtocol(ProtocolTypeMCPApps),
		WithProtocolBaseURL("https://my-cdn.example.com"),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)
	assert.Contains(t, resource.Resource.Text, "https://my-cdn.example.com/mcpapps-v1.js")
}

func TestCreateUIResource_WithRenderData(t *testing.T) {
	renderData := RenderData{
		Locale:      "en-US",
		Theme:       "dark",
		DisplayMode: DisplayModeInline,
		MaxHeight:   600,
		ToolInput: map[string]interface{}{
			"query": "test",
		},
	}

	resource, err := CreateUIResource(
		"ui://test",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<h1>Test</h1>",
		},
		EncodingText,
		WithUIMetadata(map[string]interface{}{
			UIMetadataKeyInitialRenderData: renderData,
		}),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)
	assert.NotNil(t, resource.Resource.Meta)

	// Verify render data is in metadata with correct prefix
	metaKey := UIMetadataPrefix + UIMetadataKeyInitialRenderData
	assert.Contains(t, resource.Resource.Meta, metaKey)

	// Verify render data structure
	storedData, ok := resource.Resource.Meta[metaKey].(RenderData)
	assert.True(t, ok)
	assert.Equal(t, "en-US", storedData.Locale)
	assert.Equal(t, "dark", storedData.Theme)
	assert.Equal(t, DisplayModeInline, storedData.DisplayMode)
	assert.Equal(t, 600, storedData.MaxHeight)
}

func TestCreateUIResource_MCPAppsWithRenderData(t *testing.T) {
	renderData := RenderData{
		Locale:      "fr-FR",
		Theme:       "light",
		DisplayMode: DisplayModePIP,
		MaxHeight:   800,
		ToolInput: map[string]interface{}{
			"userId": "12345",
		},
		ToolOutput: map[string]interface{}{
			"result": "success",
		},
	}

	resource, err := CreateUIResource(
		"ui://widget",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<div>Widget</div>",
		},
		EncodingText,
		WithProtocol(ProtocolTypeMCPApps),
		WithUIMetadata(map[string]interface{}{
			UIMetadataKeyInitialRenderData:  renderData,
			UIMetadataKeyPreferredFrameSize: []string{"1200px", "800px"},
		}),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)

	// Verify MIME type
	assert.Equal(t, "text/html;profile=mcp-app", resource.Resource.MimeType)

	// Verify script injection
	assert.Contains(t, resource.Resource.Text, "<script")
	assert.Contains(t, resource.Resource.Text, "mcpapps")

	// Verify render data in metadata
	metaKey := UIMetadataPrefix + UIMetadataKeyInitialRenderData
	assert.Contains(t, resource.Resource.Meta, metaKey)

	// Verify preferred frame size
	frameSizeKey := UIMetadataPrefix + UIMetadataKeyPreferredFrameSize
	assert.Contains(t, resource.Resource.Meta, frameSizeKey)
}

func TestCreateUIResource_WithResourceProps(t *testing.T) {
	// ResourceProps with _meta key adds to resource metadata
	resource, err := CreateUIResource(
		"ui://test",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<h1>Test</h1>",
		},
		EncodingText,
		WithResourceProps(map[string]interface{}{
			"_meta": map[string]interface{}{
				"custom.property": "value",
			},
		}),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)
	assert.Contains(t, resource.Resource.Meta, "custom.property")
	assert.Equal(t, "value", resource.Resource.Meta["custom.property"])
}

func TestCreateUIResource_CompleteIntegration(t *testing.T) {
	// Test combining all features: protocol, render data, metadata, props
	renderData := RenderData{
		Locale:      "ja-JP",
		Theme:       "dark",
		DisplayMode: DisplayModeFullscreen,
		MaxHeight:   1000,
		ToolInput: map[string]interface{}{
			"mode": "advanced",
		},
	}

	resource, err := CreateUIResource(
		"ui://complete-widget",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<div id='app'>Complete Widget</div>",
		},
		EncodingText,
		WithProtocol(ProtocolTypeMCPApps),
		WithProtocolVersion("v2"),
		WithUIMetadata(map[string]interface{}{
			UIMetadataKeyInitialRenderData:  renderData,
			UIMetadataKeyPreferredFrameSize: []string{"1600px", "1000px"},
		}),
		WithMetadata(map[string]interface{}{
			"author":  "Test Suite",
			"version": "1.0.0",
		}),
		WithResourceProps(map[string]interface{}{
			"_meta": map[string]interface{}{
				"custom.category": "dashboard",
			},
		}),
		WithEmbeddedResourceProps(map[string]interface{}{
			"annotations": map[string]interface{}{
				"priority": "high",
			},
		}),
	)

	assert.NoError(t, err)
	assert.NotNil(t, resource)

	// Verify resource structure
	assert.Equal(t, "resource", resource.Type)
	assert.Equal(t, "ui://complete-widget", resource.Resource.URI)
	assert.Equal(t, "text/html;profile=mcp-app", resource.Resource.MimeType)

	// Verify script injection
	assert.Contains(t, resource.Resource.Text, "mcpapps-v2.js")

	// Verify metadata
	assert.Contains(t, resource.Resource.Meta, UIMetadataPrefix+UIMetadataKeyInitialRenderData)
	assert.Contains(t, resource.Resource.Meta, UIMetadataPrefix+UIMetadataKeyPreferredFrameSize)
	assert.Contains(t, resource.Resource.Meta, "author")
	assert.Contains(t, resource.Resource.Meta, "version")
	assert.Contains(t, resource.Resource.Meta, "custom.category")

	// Verify annotations
	assert.Contains(t, resource.Annotations, "priority")

	// Verify serialization
	jsonBytes, err := json.Marshal(resource)
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonBytes)

	// Verify deserialization
	var decoded UIResource
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, resource.Resource.URI, decoded.Resource.URI)
}

func TestInvalidURIError_Is(t *testing.T) {
	err := &InvalidURIError{URI: "http://invalid"}

	assert.True(t, err.Is(ErrInvalidURI))
	assert.False(t, err.Is(ErrEmptyHTMLString))
	assert.Contains(t, err.Error(), "http://invalid")
}

func TestMIMETypeConstants(t *testing.T) {
	// Verify MIME type constant values
	assert.Equal(t, "text/html", MimeTypeHTML)
	assert.Equal(t, "text/uri-list", MimeTypeURIList)
	assert.Equal(t, "text/html+skybridge", MimeTypeAppsSdkAdapter)
	assert.Equal(t, "text/html;profile=mcp-app", MimeTypeMCPAppsAdapter)
}

func TestProtocolTypeConstants(t *testing.T) {
	// Verify protocol type constant values
	assert.Equal(t, ProtocolType("generic"), ProtocolTypeGeneric)
	assert.Equal(t, ProtocolType("appssdk"), ProtocolTypeAppsSDK)
	assert.Equal(t, ProtocolType("mcpapps"), ProtocolTypeMCPApps)
}

func TestHTMLScriptInjection(t *testing.T) {
	tests := []struct {
		name          string
		html          string
		protocol      ProtocolType
		shouldContain []string
	}{
		{
			name:     "inject into head",
			html:     "<html><head><title>Test</title></head><body>Content</body></html>",
			protocol: ProtocolTypeMCPApps,
			shouldContain: []string{
				"<head>",
				"<script",
				"mcpapps",
				"<title>Test</title>",
			},
		},
		{
			name:     "inject into body if no head",
			html:     "<html><body>Content</body></html>",
			protocol: ProtocolTypeMCPApps,
			shouldContain: []string{
				"<body>",
				"<script",
				"mcpapps",
			},
		},
		{
			name:     "wrap minimal HTML",
			html:     "<h1>Title</h1>",
			protocol: ProtocolTypeMCPApps,
			shouldContain: []string{
				"<html>",
				"<head>",
				"<script",
				"mcpapps",
				"<body>",
				"<h1>Title</h1>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource, err := CreateUIResource(
				"ui://test",
				&RawHTMLPayload{
					Type:       ContentTypeRawHTML,
					HTMLString: tt.html,
				},
				EncodingText,
				WithProtocol(tt.protocol),
			)

			assert.NoError(t, err)
			assert.NotNil(t, resource)

			for _, expected := range tt.shouldContain {
				assert.Contains(t, resource.Resource.Text, expected)
			}
		})
	}
}

func TestRenderDataEdgeCases(t *testing.T) {
	tests := []struct {
		name       string
		renderData RenderData
		wantValid  bool
	}{
		{
			name: "all fields populated",
			renderData: RenderData{
				ToolInput:   map[string]interface{}{"key": "value"},
				ToolOutput:  "result",
				WidgetState: map[string]interface{}{"state": "active"},
				Locale:      "en-US",
				Theme:       "dark",
				DisplayMode: DisplayModeInline,
				MaxHeight:   600,
			},
			wantValid: true,
		},
		{
			name:       "empty render data",
			renderData: RenderData{},
			wantValid:  true,
		},
		{
			name: "only locale",
			renderData: RenderData{
				Locale: "fr-FR",
			},
			wantValid: true,
		},
		{
			name: "zero max height",
			renderData: RenderData{
				MaxHeight: 0,
			},
			wantValid: true,
		},
		{
			name: "nil tool input",
			renderData: RenderData{
				ToolInput: nil,
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource, err := CreateUIResource(
				"ui://test",
				&RawHTMLPayload{
					Type:       ContentTypeRawHTML,
					HTMLString: "<h1>Test</h1>",
				},
				EncodingText,
				WithUIMetadata(map[string]interface{}{
					UIMetadataKeyInitialRenderData: tt.renderData,
				}),
			)

			if tt.wantValid {
				assert.NoError(t, err)
				assert.NotNil(t, resource)

				// Verify serialization works
				jsonBytes, err := json.Marshal(resource)
				assert.NoError(t, err)
				assert.NotEmpty(t, jsonBytes)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestProtocolConfigCombinations(t *testing.T) {
	// Test various combinations of protocol config options
	t.Run("protocol with version", func(t *testing.T) {
		resource, err := CreateUIResource(
			"ui://test",
			&RawHTMLPayload{Type: ContentTypeRawHTML, HTMLString: "<div/>"},
			EncodingText,
			WithProtocol(ProtocolTypeMCPApps),
			WithProtocolVersion("v10"),
		)
		assert.NoError(t, err)
		assert.Contains(t, resource.Resource.Text, "mcpapps-v10.js")
	})

	t.Run("protocol with base URL", func(t *testing.T) {
		resource, err := CreateUIResource(
			"ui://test",
			&RawHTMLPayload{Type: ContentTypeRawHTML, HTMLString: "<div/>"},
			EncodingText,
			WithProtocol(ProtocolTypeAppsSDK),
			WithProtocolBaseURL("https://my-server.com/scripts"),
		)
		assert.NoError(t, err)
		assert.Contains(t, resource.Resource.Text, "https://my-server.com/scripts/appssdk")
	})

	t.Run("protocol with version and base URL", func(t *testing.T) {
		resource, err := CreateUIResource(
			"ui://test",
			&RawHTMLPayload{Type: ContentTypeRawHTML, HTMLString: "<div/>"},
			EncodingText,
			WithProtocol(ProtocolTypeMCPApps),
			WithProtocolVersion("v5"),
			WithProtocolBaseURL("https://cdn.custom.net"),
		)
		assert.NoError(t, err)
		assert.Contains(t, resource.Resource.Text, "https://cdn.custom.net/mcpapps-v5.js")
	})
}

func TestDefaultAdapterConstants(t *testing.T) {
	// Verify default adapter constants
	assert.Equal(t, "https://cdn.mcp-ui.dev/adapters", DefaultAdapterBaseURL)
	assert.Equal(t, "v1", DefaultAdapterVersion)
}

func TestResourceURIMetaKeyUsage(t *testing.T) {
	// Simulate tool response pattern
	resource, err := CreateUIResource(
		"ui://my-widget",
		&RawHTMLPayload{
			Type:       ContentTypeRawHTML,
			HTMLString: "<h1>Widget</h1>",
		},
		EncodingText,
		WithProtocol(ProtocolTypeMCPApps),
	)

	assert.NoError(t, err)

	// Construct tool response
	toolResponse := map[string]interface{}{
		"content": []interface{}{
			resource,
		},
		"_meta": map[string]interface{}{
			ResourceURIMetaKey: resource.Resource.URI,
		},
	}

	// Verify structure
	assert.Contains(t, toolResponse, "_meta")
	meta := toolResponse["_meta"].(map[string]interface{})
	assert.Equal(t, "ui://my-widget", meta[ResourceURIMetaKey])
	assert.Equal(t, "ui/resourceUri", ResourceURIMetaKey)
}
