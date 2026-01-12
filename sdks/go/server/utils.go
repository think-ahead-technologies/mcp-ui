package mcpuiserver

import (
	"encoding/base64"
)

// encodeBase64 encodes a string to base64
func encodeBase64(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

// buildMetadata builds the metadata map from UI metadata and custom metadata
func buildMetadata(opts *CreateUIResourceOptions) map[string]interface{} {
	if opts.UIMetadata == nil && opts.Metadata == nil && opts.ResourceProps == nil {
		return nil
	}

	meta := make(map[string]interface{})

	// Add prefixed UI metadata
	for k, v := range opts.UIMetadata {
		meta[UIMetadataPrefix+k] = v
	}

	// Add custom metadata (can override)
	for k, v := range opts.Metadata {
		meta[k] = v
	}

	// Merge with resource props metadata
	if opts.ResourceProps != nil {
		if propsMeta, ok := opts.ResourceProps["_meta"].(map[string]interface{}); ok {
			for k, v := range propsMeta {
				meta[k] = v
			}
		}
	}

	if len(meta) == 0 {
		return nil
	}

	return meta
}
