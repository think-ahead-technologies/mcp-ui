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

	// Apply adapter if provided (only for RawHTML content)
	if options.Adapter != nil {
		if _, isRawHTML := content.(*RawHTMLPayload); isRawHTML {
			contentString = wrapWithAdapter(contentString, options.Adapter)
			mimeType = options.Adapter.GetMIMEType()
		}
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

// wrapWithAdapter wraps HTML content with an adapter script.
// It injects the adapter script into the <head> tag, creating one if it doesn't exist.
func wrapWithAdapter(htmlContent string, adapter Adapter) string {
	adapterScript := adapter.GetScript()

	// Check if there's a <head> tag
	headStart := findInsensitive(htmlContent, "<head>")
	headEnd := findInsensitive(htmlContent, "</head>")

	if headStart >= 0 && headEnd >= 0 && headEnd > headStart {
		// Inject adapter script after <head>
		headTagEnd := headStart + len("<head>")
		return htmlContent[:headTagEnd] + "\n" + adapterScript + htmlContent[headTagEnd:]
	}

	// Check if there's an <html> tag
	htmlStart := findInsensitive(htmlContent, "<html>")

	if htmlStart >= 0 {
		// Insert <head> with adapter after <html>
		htmlTagEnd := htmlStart + len("<html>")
		headSection := fmt.Sprintf("\n<head>\n%s\n</head>\n", adapterScript)
		return htmlContent[:htmlTagEnd] + headSection + htmlContent[htmlTagEnd:]
	}

	// No <html> or <head> tag, wrap everything
	return fmt.Sprintf("<html>\n<head>\n%s\n</head>\n<body>\n%s\n</body>\n</html>",
		adapterScript, htmlContent)
}

// findInsensitive finds a substring case-insensitively and returns the index.
// Returns -1 if not found.
func findInsensitive(s, substr string) int {
	sLower := toLower(s)
	substrLower := toLower(substr)
	return indexOf(sLower, substrLower)
}

// toLower converts a string to lowercase.
func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// indexOf finds the first occurrence of a substring.
func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
