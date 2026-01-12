package mcpuiserver

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUIResource(t *testing.T) {
	tests := []struct {
		name        string
		uri         string
		content     ResourceContentPayload
		encoding    Encoding
		opts        []Option
		want        *UIResource
		wantErr     bool
		errContains string
	}{
		{
			name: "text-based raw HTML",
			uri:  "ui://test-html",
			content: &RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<p>Test</p>",
			},
			encoding: EncodingText,
			want: &UIResource{
				Type: "resource",
				Resource: ResourceContent{
					URI:      "ui://test-html",
					MimeType: MimeTypeHTML,
					Text:     "<p>Test</p>",
				},
			},
			wantErr: false,
		},
		{
			name: "blob-based raw HTML",
			uri:  "ui://test-html-blob",
			content: &RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<h1>Blob</h1>",
			},
			encoding: EncodingBlob,
			want: &UIResource{
				Type: "resource",
				Resource: ResourceContent{
					URI:      "ui://test-html-blob",
					MimeType: MimeTypeHTML,
					Blob:     "PGgxPkJsb2I8L2gxPg==", // base64 of "<h1>Blob</h1>"
				},
			},
			wantErr: false,
		},
		{
			name: "text-based external URL",
			uri:  "ui://test-url",
			content: &ExternalURLPayload{
				Type:      ContentTypeExternalURL,
				IframeURL: "https://example.com",
			},
			encoding: EncodingText,
			want: &UIResource{
				Type: "resource",
				Resource: ResourceContent{
					URI:      "ui://test-url",
					MimeType: MimeTypeURIList,
					Text:     "https://example.com",
				},
			},
			wantErr: false,
		},
		{
			name: "blob-based external URL",
			uri:  "ui://test-url-blob",
			content: &ExternalURLPayload{
				Type:      ContentTypeExternalURL,
				IframeURL: "https://example.com/widget",
			},
			encoding: EncodingBlob,
			want: &UIResource{
				Type: "resource",
				Resource: ResourceContent{
					URI:      "ui://test-url-blob",
					MimeType: MimeTypeURIList,
					Blob:     "aHR0cHM6Ly9leGFtcGxlLmNvbS93aWRnZXQ=", // base64 of "https://example.com/widget"
				},
			},
			wantErr: false,
		},
		{
			name: "text-based remote DOM with React",
			uri:  "ui://test-remote-react",
			content: &RemoteDOMPayload{
				Type:      ContentTypeRemoteDOM,
				Script:    "console.log('React');",
				Framework: FrameworkReact,
			},
			encoding: EncodingText,
			want: &UIResource{
				Type: "resource",
				Resource: ResourceContent{
					URI:      "ui://test-remote-react",
					MimeType: MimeTypeRemoteDomReact,
					Text:     "console.log('React');",
				},
			},
			wantErr: false,
		},
		{
			name: "blob-based remote DOM with WebComponents",
			uri:  "ui://test-remote-wc",
			content: &RemoteDOMPayload{
				Type:      ContentTypeRemoteDOM,
				Script:    "console.log('WebComponents');",
				Framework: FrameworkWebComponents,
			},
			encoding: EncodingBlob,
			want: &UIResource{
				Type: "resource",
				Resource: ResourceContent{
					URI:      "ui://test-remote-wc",
					MimeType: MimeTypeRemoteDomWC,
					Blob:     "Y29uc29sZS5sb2coJ1dlYkNvbXBvbmVudHMnKTs=", // base64 of "console.log('WebComponents');"
				},
			},
			wantErr: false,
		},
		{
			name: "with UI metadata",
			uri:  "ui://test-metadata",
			content: &RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<p>Metadata</p>",
			},
			encoding: EncodingText,
			opts: []Option{
				WithUIMetadata(map[string]interface{}{
					UIMetadataKeyPreferredFrameSize: []string{"800px", "600px"},
				}),
			},
			want: &UIResource{
				Type: "resource",
				Resource: ResourceContent{
					URI:      "ui://test-metadata",
					MimeType: MimeTypeHTML,
					Text:     "<p>Metadata</p>",
					Meta: map[string]interface{}{
						UIMetadataPrefix + UIMetadataKeyPreferredFrameSize: []string{"800px", "600px"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with custom metadata",
			uri:  "ui://test-custom-metadata",
			content: &RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<p>Custom</p>",
			},
			encoding: EncodingText,
			opts: []Option{
				WithMetadata(map[string]interface{}{
					"custom.author": "TestAuthor",
				}),
			},
			want: &UIResource{
				Type: "resource",
				Resource: ResourceContent{
					URI:      "ui://test-custom-metadata",
					MimeType: MimeTypeHTML,
					Text:     "<p>Custom</p>",
					Meta: map[string]interface{}{
						"custom.author": "TestAuthor",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "with UI and custom metadata",
			uri:  "ui://test-both-metadata",
			content: &RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<p>Both</p>",
			},
			encoding: EncodingText,
			opts: []Option{
				WithUIMetadata(map[string]interface{}{
					UIMetadataKeyInitialRenderData: map[string]interface{}{
						"userId": "123",
					},
				}),
				WithMetadata(map[string]interface{}{
					"custom.version": "1.0",
				}),
			},
			want: &UIResource{
				Type: "resource",
				Resource: ResourceContent{
					URI:      "ui://test-both-metadata",
					MimeType: MimeTypeHTML,
					Text:     "<p>Both</p>",
					Meta: map[string]interface{}{
						UIMetadataPrefix + UIMetadataKeyInitialRenderData: map[string]interface{}{
							"userId": "123",
						},
						"custom.version": "1.0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid URI prefix",
			uri:  "invalid://test",
			content: &RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<p>Test</p>",
			},
			encoding:    EncodingText,
			wantErr:     true,
			errContains: "URI must start with 'ui://'",
		},
		{
			name: "empty HTML string",
			uri:  "ui://test",
			content: &RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "",
			},
			encoding:    EncodingText,
			wantErr:     true,
			errContains: "htmlString must be provided",
		},
		{
			name:        "nil content",
			uri:         "ui://test",
			content:     nil,
			encoding:    EncodingText,
			wantErr:     true,
			errContains: "content cannot be nil",
		},
		{
			name: "empty iframe URL",
			uri:  "ui://test",
			content: &ExternalURLPayload{
				Type:      ContentTypeExternalURL,
				IframeURL: "",
			},
			encoding:    EncodingText,
			wantErr:     true,
			errContains: "iframeUrl must be provided",
		},
		{
			name: "empty script",
			uri:  "ui://test",
			content: &RemoteDOMPayload{
				Type:      ContentTypeRemoteDOM,
				Script:    "",
				Framework: FrameworkReact,
			},
			encoding:    EncodingText,
			wantErr:     true,
			errContains: "script must be provided",
		},
		{
			name: "invalid framework",
			uri:  "ui://test",
			content: &RemoteDOMPayload{
				Type:      ContentTypeRemoteDOM,
				Script:    "console.log('test');",
				Framework: "invalid",
			},
			encoding:    EncodingText,
			wantErr:     true,
			errContains: "framework must be",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateUIResource(tt.uri, tt.content, tt.encoding, tt.opts...)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCreateUIResourceJSON(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		content  ResourceContentPayload
		encoding Encoding
		wantJSON string
	}{
		{
			name: "text HTML JSON serialization",
			uri:  "ui://test",
			content: &RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<p>Test</p>",
			},
			encoding: EncodingText,
			wantJSON: `{
				"type": "resource",
				"resource": {
					"uri": "ui://test",
					"mimeType": "text/html",
					"text": "<p>Test</p>"
				}
			}`,
		},
		{
			name: "blob HTML JSON serialization",
			uri:  "ui://test",
			content: &RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<p>Test</p>",
			},
			encoding: EncodingBlob,
			wantJSON: `{
				"type": "resource",
				"resource": {
					"uri": "ui://test",
					"mimeType": "text/html",
					"blob": "PHA+VGVzdDwvcD4="
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateUIResource(tt.uri, tt.content, tt.encoding)
			assert.NoError(t, err)

			gotJSON, err := json.Marshal(got)
			assert.NoError(t, err)

			assert.JSONEq(t, tt.wantJSON, string(gotJSON))
		})
	}
}

func TestCreateUIResourceWithEmbeddedProps(t *testing.T) {
	uri := "ui://test"
	content := &RawHTMLPayload{
		Type:       ContentTypeRawHTML,
		HTMLString: "<p>Test</p>",
	}

	got, err := CreateUIResource(
		uri,
		content,
		EncodingText,
		WithEmbeddedResourceProps(map[string]interface{}{
			"annotations": map[string]interface{}{
				"priority": 1,
			},
			"_meta": map[string]interface{}{
				"custom": "value",
			},
		}),
	)

	assert.NoError(t, err)
	assert.Equal(t, "resource", got.Type)
	assert.Equal(t, map[string]interface{}{"priority": 1}, got.Annotations)
	assert.Equal(t, map[string]interface{}{"custom": "value"}, got.Meta)
}
