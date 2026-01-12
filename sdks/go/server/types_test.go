package mcpuiserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateURI(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{
			name:    "valid URI",
			uri:     "ui://test",
			wantErr: false,
		},
		{
			name:    "valid URI with path",
			uri:     "ui://component/widget/123",
			wantErr: false,
		},
		{
			name:    "invalid prefix http",
			uri:     "http://test",
			wantErr: true,
		},
		{
			name:    "invalid prefix https",
			uri:     "https://test",
			wantErr: true,
		},
		{
			name:    "invalid prefix other",
			uri:     "invalid://test",
			wantErr: true,
		},
		{
			name:    "empty URI",
			uri:     "",
			wantErr: true,
		},
		{
			name:    "just ui prefix",
			uri:     "ui://",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURI(tt.uri)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRawHTMLPayloadValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload RawHTMLPayload
		wantErr bool
		errType error
	}{
		{
			name: "valid raw HTML",
			payload: RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<p>Test</p>",
			},
			wantErr: false,
		},
		{
			name: "valid raw HTML with complex content",
			payload: RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<html><body><h1>Test</h1><script>console.log('test');</script></body></html>",
			},
			wantErr: false,
		},
		{
			name: "empty HTML string",
			payload: RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "",
			},
			wantErr: true,
			errType: ErrEmptyHTMLString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.validate()
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

func TestExternalURLPayloadValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload ExternalURLPayload
		wantErr bool
		errType error
	}{
		{
			name: "valid external URL",
			payload: ExternalURLPayload{
				Type:      ContentTypeExternalURL,
				IframeURL: "https://example.com",
			},
			wantErr: false,
		},
		{
			name: "valid external URL with path",
			payload: ExternalURLPayload{
				Type:      ContentTypeExternalURL,
				IframeURL: "https://example.com/widget?id=123",
			},
			wantErr: false,
		},
		{
			name: "empty iframe URL",
			payload: ExternalURLPayload{
				Type:      ContentTypeExternalURL,
				IframeURL: "",
			},
			wantErr: true,
			errType: ErrEmptyIframeURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.validate()
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

func TestRemoteDOMPayloadValidation(t *testing.T) {
	tests := []struct {
		name    string
		payload RemoteDOMPayload
		wantErr bool
		errType error
	}{
		{
			name: "valid remote DOM with React",
			payload: RemoteDOMPayload{
				Type:      ContentTypeRemoteDOM,
				Script:    "console.log('test');",
				Framework: FrameworkReact,
			},
			wantErr: false,
		},
		{
			name: "valid remote DOM with WebComponents",
			payload: RemoteDOMPayload{
				Type:      ContentTypeRemoteDOM,
				Script:    "console.log('test');",
				Framework: FrameworkWebComponents,
			},
			wantErr: false,
		},
		{
			name: "empty script",
			payload: RemoteDOMPayload{
				Type:      ContentTypeRemoteDOM,
				Script:    "",
				Framework: FrameworkReact,
			},
			wantErr: true,
			errType: ErrEmptyScript,
		},
		{
			name: "invalid framework",
			payload: RemoteDOMPayload{
				Type:      ContentTypeRemoteDOM,
				Script:    "console.log('test');",
				Framework: "invalid",
			},
			wantErr: true,
			errType: ErrInvalidFramework,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.payload.validate()
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

func TestContentTypeMethod(t *testing.T) {
	tests := []struct {
		name     string
		payload  ResourceContentPayload
		wantType ContentType
	}{
		{
			name: "raw HTML content type",
			payload: &RawHTMLPayload{
				Type:       ContentTypeRawHTML,
				HTMLString: "<p>Test</p>",
			},
			wantType: ContentTypeRawHTML,
		},
		{
			name: "external URL content type",
			payload: &ExternalURLPayload{
				Type:      ContentTypeExternalURL,
				IframeURL: "https://example.com",
			},
			wantType: ContentTypeExternalURL,
		},
		{
			name: "remote DOM content type",
			payload: &RemoteDOMPayload{
				Type:      ContentTypeRemoteDOM,
				Script:    "console.log('test');",
				Framework: FrameworkReact,
			},
			wantType: ContentTypeRemoteDOM,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.payload.contentType()
			assert.Equal(t, tt.wantType, got)
		})
	}
}
