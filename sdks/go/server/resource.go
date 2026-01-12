package mcpuiserver

import (
	"fmt"
)

// CreateUIResource creates a UIResource for inclusion in MCP tool results.
//
// The function validates the URI (must start with "ui://"), processes the content
// based on its type, and applies the specified encoding.
//
// Parameters:
//   - uri: Resource identifier starting with "ui://"
//   - content: Content payload (RawHTMLPayload, ExternalURLPayload, or RemoteDOMPayload)
//   - encoding: Encoding type (EncodingText or EncodingBlob)
//   - opts: Optional functional options for metadata and properties
//
// Returns:
//   - *UIResource: The created UI resource
//   - error: Validation or processing errors
//
// Example:
//
//	resource, err := CreateUIResource(
//	    "ui://greeting",
//	    &RawHTMLPayload{
//	        Type: ContentTypeRawHTML,
//	        HTMLString: "<h1>Hello, World!</h1>",
//	    },
//	    EncodingText,
//	)
func CreateUIResource(uri string, content ResourceContentPayload, encoding Encoding, opts ...Option) (*UIResource, error) {
	// Validate URI
	if err := validateURI(uri); err != nil {
		return nil, err
	}

	// Validate content
	if content == nil {
		return nil, ErrNilContent
	}
	if err := content.validate(); err != nil {
		return nil, err
	}

	// Validate encoding
	if encoding != EncodingText && encoding != EncodingBlob {
		return nil, ErrInvalidEncoding
	}

	// Apply options
	options := &CreateUIResourceOptions{
		URI:      uri,
		Content:  content,
		Encoding: encoding,
	}
	for _, opt := range opts {
		opt(options)
	}

	// Determine content string and MIME type
	var contentString string
	var mimeType string

	switch c := content.(type) {
	case *RawHTMLPayload:
		contentString = c.HTMLString
		mimeType = MimeTypeHTML
	case *ExternalURLPayload:
		contentString = c.IframeURL
		mimeType = MimeTypeURIList
	case *RemoteDOMPayload:
		contentString = c.Script
		if c.Framework == FrameworkReact {
			mimeType = MimeTypeRemoteDomReact
		} else {
			mimeType = MimeTypeRemoteDomWC
		}
	default:
		return nil, fmt.Errorf("unsupported content type: %T", content)
	}

	// Build resource content
	resourceContent := ResourceContent{
		URI:      uri,
		MimeType: mimeType,
	}

	// Apply encoding
	switch encoding {
	case EncodingText:
		resourceContent.Text = contentString
	case EncodingBlob:
		resourceContent.Blob = encodeBase64(contentString)
	}

	// Add metadata
	resourceContent.Meta = buildMetadata(options)

	// Build UI resource
	resource := &UIResource{
		Type:     "resource",
		Resource: resourceContent,
	}

	// Add embedded resource props
	if options.EmbeddedResourceProps != nil {
		if annotations, ok := options.EmbeddedResourceProps["annotations"]; ok {
			if annotationsMap, ok := annotations.(map[string]interface{}); ok {
				resource.Annotations = annotationsMap
			}
		}
		if meta, ok := options.EmbeddedResourceProps["_meta"]; ok {
			if metaMap, ok := meta.(map[string]interface{}); ok {
				resource.Meta = metaMap
			}
		}
	}

	return resource, nil
}
