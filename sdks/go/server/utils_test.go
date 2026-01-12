package mcpuiserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeBase64(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple ASCII",
			input: "hello world",
			want:  "aGVsbG8gd29ybGQ=",
		},
		{
			name:  "UTF-8 characters",
			input: "你好,世界",
			want:  "5L2g5aW9LOS4lueVjA==",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "special characters",
			input: "`~!@#$%^&*()_+-=[]{}\\|;':\",./<>?",
			want:  "YH4hQCMkJV4mKigpXystPVtde31cfDsnOiIsLi88Pj8=",
		},
		{
			name:  "HTML content",
			input: "<p>Test</p>",
			want:  "PHA+VGVzdDwvcD4=",
		},
		{
			name:  "JavaScript content",
			input: "console.log('test');",
			want:  "Y29uc29sZS5sb2coJ3Rlc3QnKTs=",
		},
		{
			name:  "URL",
			input: "https://example.com/widget",
			want:  "aHR0cHM6Ly9leGFtcGxlLmNvbS93aWRnZXQ=",
		},
		{
			name:  "newlines",
			input: "line1\nline2\nline3",
			want:  "bGluZTEKbGluZTIKbGluZTM=",
		},
		{
			name:  "emoji",
			input: "Hello 👋 World 🌍",
			want:  "SGVsbG8g8J+RiyBXb3JsZCDwn4yN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeBase64(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBuildMetadata(t *testing.T) {
	tests := []struct {
		name string
		opts *CreateUIResourceOptions
		want map[string]interface{}
	}{
		{
			name: "UI metadata with prefix",
			opts: &CreateUIResourceOptions{
				UIMetadata: map[string]interface{}{
					UIMetadataKeyPreferredFrameSize: []string{"800px", "600px"},
				},
			},
			want: map[string]interface{}{
				UIMetadataPrefix + UIMetadataKeyPreferredFrameSize: []string{"800px", "600px"},
			},
		},
		{
			name: "custom metadata without prefix",
			opts: &CreateUIResourceOptions{
				Metadata: map[string]interface{}{
					"custom.author": "TestAuthor",
				},
			},
			want: map[string]interface{}{
				"custom.author": "TestAuthor",
			},
		},
		{
			name: "UI and custom metadata",
			opts: &CreateUIResourceOptions{
				UIMetadata: map[string]interface{}{
					UIMetadataKeyInitialRenderData: map[string]interface{}{
						"userId": "123",
					},
				},
				Metadata: map[string]interface{}{
					"custom.version": "1.0",
				},
			},
			want: map[string]interface{}{
				UIMetadataPrefix + UIMetadataKeyInitialRenderData: map[string]interface{}{
					"userId": "123",
				},
				"custom.version": "1.0",
			},
		},
		{
			name: "metadata override order",
			opts: &CreateUIResourceOptions{
				UIMetadata: map[string]interface{}{
					"test": "ui-value",
				},
				Metadata: map[string]interface{}{
					"test": "custom-value",
				},
			},
			want: map[string]interface{}{
				UIMetadataPrefix + "test": "ui-value",
				"test":                    "custom-value",
			},
		},
		{
			name: "resource props metadata",
			opts: &CreateUIResourceOptions{
				ResourceProps: map[string]interface{}{
					"_meta": map[string]interface{}{
						"prop": "value",
					},
				},
			},
			want: map[string]interface{}{
				"prop": "value",
			},
		},
		{
			name: "all metadata types merged",
			opts: &CreateUIResourceOptions{
				UIMetadata: map[string]interface{}{
					"ui-key": "ui-value",
				},
				Metadata: map[string]interface{}{
					"custom-key": "custom-value",
				},
				ResourceProps: map[string]interface{}{
					"_meta": map[string]interface{}{
						"prop-key": "prop-value",
					},
				},
			},
			want: map[string]interface{}{
				UIMetadataPrefix + "ui-key": "ui-value",
				"custom-key":                "custom-value",
				"prop-key":                  "prop-value",
			},
		},
		{
			name: "no metadata",
			opts: &CreateUIResourceOptions{},
			want: nil,
		},
		{
			name: "empty metadata maps",
			opts: &CreateUIResourceOptions{
				UIMetadata: map[string]interface{}{},
				Metadata:   map[string]interface{}{},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildMetadata(tt.opts)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBuildMetadataOverridePrecedence(t *testing.T) {
	// Test that resource props metadata can override both UI and custom metadata
	opts := &CreateUIResourceOptions{
		UIMetadata: map[string]interface{}{
			"shared": "from-ui",
		},
		Metadata: map[string]interface{}{
			"shared": "from-custom",
		},
		ResourceProps: map[string]interface{}{
			"_meta": map[string]interface{}{
				"shared":                    "from-props",
				UIMetadataPrefix + "shared": "from-props-ui",
			},
		},
	}

	got := buildMetadata(opts)

	// Resource props should override previous values
	assert.Equal(t, "from-props", got["shared"])
	assert.Equal(t, "from-props-ui", got[UIMetadataPrefix+"shared"])
}
