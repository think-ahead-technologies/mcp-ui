package appssdk

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
		wantIntent  string
		wantOrigin  string
		wantErr     bool
		errType     error
	}{
		{
			name:        "default configuration",
			opts:        []Option{},
			wantTimeout: 30000,
			wantIntent:  "prompt",
			wantOrigin:  "",
			wantErr:     false,
		},
		{
			name:        "custom timeout",
			opts:        []Option{WithTimeout(5000)},
			wantTimeout: 5000,
			wantIntent:  "prompt",
			wantOrigin:  "",
			wantErr:     false,
		},
		{
			name:        "custom intent handling - ignore",
			opts:        []Option{WithIntentHandling("ignore")},
			wantTimeout: 30000,
			wantIntent:  "ignore",
			wantOrigin:  "",
			wantErr:     false,
		},
		{
			name:        "custom host origin",
			opts:        []Option{WithHostOrigin("https://example.com")},
			wantTimeout: 30000,
			wantIntent:  "prompt",
			wantOrigin:  "https://example.com",
			wantErr:     false,
		},
		{
			name: "combined custom config",
			opts: []Option{
				WithTimeout(10000),
				WithIntentHandling("ignore"),
				WithHostOrigin("https://test.com"),
			},
			wantTimeout: 10000,
			wantIntent:  "ignore",
			wantOrigin:  "https://test.com",
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
		{
			name:    "invalid intent handling",
			opts:    []Option{WithIntentHandling("invalid")},
			wantErr: true,
			errType: ErrInvalidIntentHandling,
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
				assert.Equal(t, tt.wantIntent, adapter.config.IntentHandling)
				assert.Equal(t, tt.wantOrigin, adapter.config.HostOrigin)
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
				"MCPUIAppsSdkAdapter",
				"initWithConfig",
			},
			configContains: map[string]interface{}{
				"timeout":        30000,
				"intentHandling": "prompt",
			},
		},
		{
			name: "custom timeout in script",
			opts: []Option{WithTimeout(5000)},
			wantInScript: []string{
				"<script>",
				"</script>",
			},
			configContains: map[string]interface{}{
				"timeout":        5000,
				"intentHandling": "prompt",
			},
		},
		{
			name: "custom intent handling in script",
			opts: []Option{WithIntentHandling("ignore")},
			wantInScript: []string{
				"<script>",
				"</script>",
			},
			configContains: map[string]interface{}{
				"timeout":        30000,
				"intentHandling": "ignore",
			},
		},
		{
			name: "custom host origin in script",
			opts: []Option{WithHostOrigin("https://example.com")},
			wantInScript: []string{
				"<script>",
				"</script>",
			},
			configContains: map[string]interface{}{
				"timeout":        30000,
				"intentHandling": "prompt",
				"hostOrigin":     "https://example.com",
			},
		},
		{
			name: "all custom config in script",
			opts: []Option{
				WithTimeout(10000),
				WithIntentHandling("ignore"),
				WithHostOrigin("https://test.com"),
			},
			wantInScript: []string{
				"<script>",
				"</script>",
			},
			configContains: map[string]interface{}{
				"timeout":        10000,
				"intentHandling": "ignore",
				"hostOrigin":     "https://test.com",
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
				// (may be part of a larger JSON object)
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
	assert.Equal(t, "text/html+skybridge", mimeType)
}

func TestAdapter_GetType(t *testing.T) {
	adapter, err := NewAdapter()
	assert.NoError(t, err)

	adapterType := adapter.GetType()
	assert.Equal(t, string(adapters.AdapterTypeAppsSDK), adapterType)
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
				Timeout:        30000,
				IntentHandling: "prompt",
				HostOrigin:     "",
			},
			wantErr: false,
		},
		{
			name: "valid config with ignore",
			config: Config{
				Timeout:        5000,
				IntentHandling: "ignore",
				HostOrigin:     "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "invalid timeout - zero",
			config: Config{
				Timeout:        0,
				IntentHandling: "prompt",
			},
			wantErr: true,
			errType: ErrInvalidTimeout,
		},
		{
			name: "invalid timeout - negative",
			config: Config{
				Timeout:        -1000,
				IntentHandling: "prompt",
			},
			wantErr: true,
			errType: ErrInvalidTimeout,
		},
		{
			name: "invalid intent handling",
			config: Config{
				Timeout:        30000,
				IntentHandling: "invalid",
			},
			wantErr: true,
			errType: ErrInvalidIntentHandling,
		},
		{
			name: "invalid intent handling - empty",
			config: Config{
				Timeout:        30000,
				IntentHandling: "",
			},
			wantErr: true,
			errType: ErrInvalidIntentHandling,
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

	// Verify script contains the adapter class
	assert.Contains(t, script, "MCPUIAppsSdkAdapter")
	assert.Contains(t, script, "function initAdapter")
	assert.Contains(t, script, "function uninstallAdapter")
}

func TestAdapter_SpecialCharactersInConfig(t *testing.T) {
	tests := []struct {
		name       string
		hostOrigin string
	}{
		{
			name:       "URL with query parameters",
			hostOrigin: "https://example.com?param=value&other=test",
		},
		{
			name:       "URL with special characters",
			hostOrigin: "https://example.com/path/to/resource",
		},
		{
			name:       "URL with fragment",
			hostOrigin: "https://example.com#section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adapter, err := NewAdapter(WithHostOrigin(tt.hostOrigin))
			assert.NoError(t, err)

			script := adapter.GetScript()

			// Verify script is valid (contains opening and closing tags)
			assert.True(t, strings.HasPrefix(script, "<script>"))
			assert.True(t, strings.HasSuffix(script, "</script>"))

			// Verify hostOrigin is properly escaped in JSON
			assert.Contains(t, script, `"hostOrigin"`)
		})
	}
}

func TestAdapter_EmptyOptionalFields(t *testing.T) {
	adapter, err := NewAdapter()
	assert.NoError(t, err)

	script := adapter.GetScript()

	// When hostOrigin is empty (default), it should not be in the config JSON
	// We can verify this by checking the serializable config
	config := adapter.serializableConfig()
	_, hasHostOrigin := config["hostOrigin"]
	assert.False(t, hasHostOrigin, "empty hostOrigin should not be in config")

	// But script should still be valid
	assert.Contains(t, script, "<script>")
	assert.Contains(t, script, "</script>")
	assert.Contains(t, script, "MCPUIAppsSdkAdapter")
}
