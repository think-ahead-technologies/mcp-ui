package mcpapps

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/MCP-UI-Org/mcp-ui/sdks/go/server/adapters"
	"github.com/stretchr/testify/assert"
)

func TestNewAdapter(t *testing.T) {
	tests := []struct {
		name        string
		opts        []Option
		wantTimeout int
		wantErr     bool
		errType     error
	}{
		{
			name:        "default configuration",
			opts:        []Option{},
			wantTimeout: 30000,
			wantErr:     false,
		},
		{
			name:        "custom timeout",
			opts:        []Option{WithTimeout(5000)},
			wantTimeout: 5000,
			wantErr:     false,
		},
		{
			name:        "custom timeout - 10 seconds",
			opts:        []Option{WithTimeout(10000)},
			wantTimeout: 10000,
			wantErr:     false,
		},
		{
			name:    "invalid timeout - zero",
			opts:    []Option{WithTimeout(0)},
			wantErr: true,
			errType: ErrInvalidTimeout,
		},
		{
			name:    "invalid timeout - negative",
			opts:    []Option{WithTimeout(-1000)},
			wantErr: true,
			errType: ErrInvalidTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, err := NewAdapter(tt.opts...)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
				assert.Nil(t, adapter)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, adapter)
				assert.Equal(t, tt.wantTimeout, adapter.config.Timeout)
			}
		})
	}
}

func TestAdapter_GetScript(t *testing.T) {
	tests := []struct {
		name           string
		opts           []Option
		wantInScript   []string
		configContains map[string]interface{}
	}{
		{
			name: "default config in script",
			opts: []Option{},
			wantInScript: []string{
				"<script>",
				"</script>",
				"const config =",
				"McpAppsAdapter",
				"initAdapter",
				"LATEST_PROTOCOL_VERSION",
				"METHODS",
			},
			configContains: map[string]interface{}{
				"timeout": 30000,
			},
		},
		{
			name: "custom timeout in script",
			opts: []Option{WithTimeout(5000)},
			wantInScript: []string{
				"<script>",
				"</script>",
				"McpAppsAdapter",
			},
			configContains: map[string]interface{}{
				"timeout": 5000,
			},
		},
		{
			name: "custom timeout - 10 seconds",
			opts: []Option{WithTimeout(10000)},
			wantInScript: []string{
				"<script>",
				"</script>",
			},
			configContains: map[string]interface{}{
				"timeout": 10000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, err := NewAdapter(tt.opts...)
			assert.NoError(t, err)

			script := adapter.GetScript()

			// Check that script contains expected strings
			for _, want := range tt.wantInScript {
				assert.Contains(t, script, want, "script should contain %q", want)
			}

			// Verify config JSON is in the script
			for key, expectedValue := range tt.configContains {
				configJSON, err := json.Marshal(map[string]interface{}{key: expectedValue})
				assert.NoError(t, err)

				// Check that the key-value pair appears in the script
				configStr := string(configJSON)
				configStr = configStr[1 : len(configStr)-1] // Remove outer braces
				assert.Contains(t, script, configStr, "script should contain config %s=%v", key, expectedValue)
			}

			// Verify script structure
			assert.True(t, strings.HasPrefix(script, "<script>"), "script should start with <script>")
			assert.True(t, strings.HasSuffix(script, "</script>"), "script should end with </script>")
		})
	}
}

func TestAdapter_GetMIMEType(t *testing.T) {
	adapter, err := NewAdapter()
	assert.NoError(t, err)

	mimeType := adapter.GetMIMEType()
	assert.Equal(t, "text/html", mimeType)
}

func TestAdapter_GetType(t *testing.T) {
	adapter, err := NewAdapter()
	assert.NoError(t, err)

	adapterType := adapter.GetType()
	assert.Equal(t, string(adapters.AdapterTypeMCPApps), adapterType)
}

func TestAdapter_ImplementsInterface(t *testing.T) {
	adapter, err := NewAdapter()
	assert.NoError(t, err)

	// Verify adapter implements the Adapter interface
	var _ adapters.Adapter = adapter
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errType error
	}{
		{
			name: "valid config",
			config: Config{
				Timeout: 30000,
			},
			wantErr: false,
		},
		{
			name: "valid config - custom timeout",
			config: Config{
				Timeout: 5000,
			},
			wantErr: false,
		},
		{
			name: "invalid timeout - zero",
			config: Config{
				Timeout: 0,
			},
			wantErr: true,
			errType: ErrInvalidTimeout,
		},
		{
			name: "invalid timeout - negative",
			config: Config{
				Timeout: -1000,
			},
			wantErr: true,
			errType: ErrInvalidTimeout,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAdapter_ScriptContainsAdapterClass(t *testing.T) {
	adapter, err := NewAdapter()
	assert.NoError(t, err)

	script := adapter.GetScript()

	// Verify script contains the adapter class and protocol constants
	assert.Contains(t, script, "McpAppsAdapter")
	assert.Contains(t, script, "function initAdapter")
	assert.Contains(t, script, "function uninstallAdapter")
	assert.Contains(t, script, "LATEST_PROTOCOL_VERSION")
	assert.Contains(t, script, `"2025-11-21"`, "should contain protocol version")
}

func TestAdapter_ScriptContainsProtocolMethods(t *testing.T) {
	adapter, err := NewAdapter()
	assert.NoError(t, err)

	script := adapter.GetScript()

	// Verify script contains JSON-RPC protocol methods
	protocolMethods := []string{
		"ui/initialize",
		"ui/notifications/initialized",
		"ui/notifications/tool-input",
		"ui/notifications/tool-result",
		"tools/call",
		"ui/open-link",
		"ui/message",
	}

	for _, method := range protocolMethods {
		assert.Contains(t, script, method, "script should contain protocol method %q", method)
	}
}

func TestAdapter_SerializableConfig(t *testing.T) {
	tests := []struct {
		name    string
		opts    []Option
		wantCfg map[string]interface{}
	}{
		{
			name: "default timeout",
			opts: []Option{},
			wantCfg: map[string]interface{}{
				"timeout": 30000,
			},
		},
		{
			name: "custom timeout",
			opts: []Option{WithTimeout(5000)},
			wantCfg: map[string]interface{}{
				"timeout": 5000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, err := NewAdapter(tt.opts...)
			assert.NoError(t, err)

			config := adapter.serializableConfig()
			assert.Equal(t, tt.wantCfg, config)
		})
	}
}

func TestAdapter_ScriptInitialization(t *testing.T) {
	adapter, err := NewAdapter()
	assert.NoError(t, err)

	script := adapter.GetScript()

	// Verify script contains initialization code
	assert.Contains(t, script, "window.McpAppsAdapter")
	assert.Contains(t, script, "initWithConfig")
	assert.Contains(t, script, ".initWithConfig()")
}

func TestAdapter_JSONRPCConstants(t *testing.T) {
	adapter, err := NewAdapter()
	assert.NoError(t, err)

	script := adapter.GetScript()

	// Verify JSON-RPC constants are present
	assert.Contains(t, script, `jsonrpc: "2.0"`)
	assert.Contains(t, script, "METHODS")
}
