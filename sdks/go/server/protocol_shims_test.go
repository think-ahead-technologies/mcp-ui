package mcpuiserver

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenericProtocolShim_GenerateScriptTag(t *testing.T) {
	shim := &GenericProtocolShim{}
	script := shim.GenerateScriptTag()
	assert.Equal(t, "", script, "generic protocol should not generate script tag")
}

func TestGenericProtocolShim_GetMIMEType(t *testing.T) {
	shim := &GenericProtocolShim{}
	mimeType := shim.GetMIMEType()
	assert.Equal(t, MimeTypeHTML, mimeType)
}

func TestAppsSdkProtocolShim_GenerateScriptTag(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		version string
		config  map[string]interface{}
		wantURL string
	}{
		{
			name:    "default configuration",
			baseURL: "https://cdn.example.com",
			version: "v1",
			config:  nil,
			wantURL: "https://cdn.example.com/appssdk-v1.js",
		},
		{
			name:    "with custom config",
			baseURL: "https://cdn.example.com",
			version: "v2",
			config: map[string]interface{}{
				"timeout":        5000,
				"intentHandling": "ignore",
			},
			wantURL: "https://cdn.example.com/appssdk-v2.js",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shim := &AppsSdkProtocolShim{
				BaseURL: tt.baseURL,
				Version: tt.version,
				Config:  tt.config,
			}

			script := shim.GenerateScriptTag()

			// Verify script tag structure
			assert.True(t, strings.HasPrefix(script, "<script src=\""))
			assert.True(t, strings.HasSuffix(script, "</script>"))
			assert.Contains(t, script, tt.wantURL)
			assert.Contains(t, script, "data-mcp-config=")

			// Verify config is valid JSON
			if tt.config != nil {
				// Extract config from script tag
				parts := strings.Split(script, "data-mcp-config='")
				assert.Len(t, parts, 2)
				configPart := strings.Split(parts[1], "'")[0]

				var parsedConfig map[string]interface{}
				err := json.Unmarshal([]byte(configPart), &parsedConfig)
				assert.NoError(t, err)
				// Compare keys (JSON unmarshaling converts numbers to float64)
				assert.Len(t, parsedConfig, len(tt.config))
				for k := range tt.config {
					assert.Contains(t, parsedConfig, k)
				}
			}
		})
	}
}

func TestAppsSdkProtocolShim_GetMIMEType(t *testing.T) {
	shim := &AppsSdkProtocolShim{}
	mimeType := shim.GetMIMEType()
	assert.Equal(t, MimeTypeAppsSdkAdapter, mimeType)
	assert.Equal(t, "text/html+skybridge", mimeType)
}

func TestMcpAppsProtocolShim_GenerateScriptTag(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		version string
		config  map[string]interface{}
		wantURL string
	}{
		{
			name:    "default configuration",
			baseURL: "https://cdn.example.com",
			version: "v1",
			config:  nil,
			wantURL: "https://cdn.example.com/mcpapps-v1.js",
		},
		{
			name:    "with custom timeout",
			baseURL: "https://cdn.example.com",
			version: "v2",
			config: map[string]interface{}{
				"timeout": 10000,
			},
			wantURL: "https://cdn.example.com/mcpapps-v2.js",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shim := &McpAppsProtocolShim{
				BaseURL: tt.baseURL,
				Version: tt.version,
				Config:  tt.config,
			}

			script := shim.GenerateScriptTag()

			// Verify script tag structure
			assert.True(t, strings.HasPrefix(script, "<script src=\""))
			assert.True(t, strings.HasSuffix(script, "</script>"))
			assert.Contains(t, script, tt.wantURL)
			assert.Contains(t, script, "data-mcp-config=")

			// Verify config is valid JSON
			if tt.config != nil {
				// Extract config from script tag
				parts := strings.Split(script, "data-mcp-config='")
				assert.Len(t, parts, 2)
				configPart := strings.Split(parts[1], "'")[0]

				var parsedConfig map[string]interface{}
				err := json.Unmarshal([]byte(configPart), &parsedConfig)
				assert.NoError(t, err)
				// Compare keys (JSON unmarshaling converts numbers to float64)
				assert.Len(t, parsedConfig, len(tt.config))
				for k := range tt.config {
					assert.Contains(t, parsedConfig, k)
				}
			}
		})
	}
}

func TestMcpAppsProtocolShim_GetMIMEType(t *testing.T) {
	shim := &McpAppsProtocolShim{}
	mimeType := shim.GetMIMEType()
	assert.Equal(t, MimeTypeMCPAppsAdapter, mimeType)
	assert.Equal(t, "text/html;profile=mcp-app", mimeType)
}

func TestGetProtocolShimGenerator(t *testing.T) {
	tests := []struct {
		name           string
		config         *ProtocolConfig
		expectedType   string
		expectedMIME   string
		expectedScript bool
	}{
		{
			name: "generic protocol",
			config: &ProtocolConfig{
				Type: ProtocolTypeGeneric,
			},
			expectedType:   "GenericProtocolShim",
			expectedMIME:   MimeTypeHTML,
			expectedScript: false,
		},
		{
			name: "apps SDK protocol",
			config: &ProtocolConfig{
				Type: ProtocolTypeAppsSDK,
			},
			expectedType:   "AppsSdkProtocolShim",
			expectedMIME:   MimeTypeAppsSdkAdapter,
			expectedScript: true,
		},
		{
			name: "MCP Apps protocol",
			config: &ProtocolConfig{
				Type: ProtocolTypeMCPApps,
			},
			expectedType:   "McpAppsProtocolShim",
			expectedMIME:   MimeTypeMCPAppsAdapter,
			expectedScript: true,
		},
		{
			name: "with custom base URL",
			config: &ProtocolConfig{
				Type:    ProtocolTypeMCPApps,
				BaseURL: "https://custom.cdn.com",
			},
			expectedType:   "McpAppsProtocolShim",
			expectedMIME:   MimeTypeMCPAppsAdapter,
			expectedScript: true,
		},
		{
			name: "with custom version",
			config: &ProtocolConfig{
				Type:    ProtocolTypeMCPApps,
				Version: "v5",
			},
			expectedType:   "McpAppsProtocolShim",
			expectedMIME:   MimeTypeMCPAppsAdapter,
			expectedScript: true,
		},
		{
			name: "defaults applied",
			config: &ProtocolConfig{
				Type: ProtocolTypeMCPApps,
			},
			expectedType:   "McpAppsProtocolShim",
			expectedMIME:   MimeTypeMCPAppsAdapter,
			expectedScript: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shim := getProtocolShimGenerator(tt.config)
			assert.NotNil(t, shim)

			mimeType := shim.GetMIMEType()
			assert.Equal(t, tt.expectedMIME, mimeType)

			script := shim.GenerateScriptTag()
			if tt.expectedScript {
				assert.NotEmpty(t, script)
				assert.Contains(t, script, "<script")
			} else {
				assert.Empty(t, script)
			}
		})
	}
}

func TestGetProtocolShimGenerator_DefaultValues(t *testing.T) {
	t.Run("default base URL", func(t *testing.T) {
		config := &ProtocolConfig{
			Type: ProtocolTypeMCPApps,
		}
		shim := getProtocolShimGenerator(config)
		mcpAppsShim, ok := shim.(*McpAppsProtocolShim)
		assert.True(t, ok)
		assert.Equal(t, DefaultAdapterBaseURL, mcpAppsShim.BaseURL)
	})

	t.Run("default version", func(t *testing.T) {
		config := &ProtocolConfig{
			Type: ProtocolTypeMCPApps,
		}
		shim := getProtocolShimGenerator(config)
		mcpAppsShim, ok := shim.(*McpAppsProtocolShim)
		assert.True(t, ok)
		assert.Equal(t, DefaultAdapterVersion, mcpAppsShim.Version)
	})
}

func TestProtocolShimGenerator_Interface(t *testing.T) {
	// Verify all shims implement the interface
	var _ ProtocolShimGenerator = &GenericProtocolShim{}
	var _ ProtocolShimGenerator = &AppsSdkProtocolShim{}
	var _ ProtocolShimGenerator = &McpAppsProtocolShim{}
}
