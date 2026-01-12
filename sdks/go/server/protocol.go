package mcpuiserver

// ParseProtocolFromInitialize extracts protocol preference from MCP initialize request metadata.
// The client can declare its preferred protocol by including "mcp-ui-protocol" in the metadata
// field of the initialize request params.
//
// Supported protocol values:
//   - "appssdk" - ChatGPT/Apps SDK protocol
//   - "mcpapps" - MCP Apps SEP (Streaming Extensible Protocol)
//   - "generic" - Standard MCP-UI protocol (no adapter)
//
// If no protocol is specified or an unsupported value is provided, returns ProtocolTypeGeneric.
//
// Example usage:
//
//	func handleInitialize(req InitializeRequest) InitializeResponse {
//	    protocol := mcpuiserver.ParseProtocolFromInitialize(req.Params)
//	    // Store protocol in session context
//	    sessions[sessionID] = &SessionContext{Protocol: protocol}
//	    return InitializeResponse{...}
//	}
func ParseProtocolFromInitialize(initializeParams map[string]interface{}) ProtocolType {
	if metadata, ok := initializeParams["metadata"].(map[string]interface{}); ok {
		if protocol, ok := metadata["mcp-ui-protocol"].(string); ok {
			switch protocol {
			case "appssdk":
				return ProtocolTypeAppsSDK
			case "mcpapps":
				return ProtocolTypeMCPApps
			case "generic":
				return ProtocolTypeGeneric
			}
		}
	}
	return ProtocolTypeGeneric // Default fallback
}

// ParseProtocolConfig extracts full protocol configuration from MCP initialize request metadata.
// This includes the protocol type and any additional protocol-specific configuration
// provided by the client in the "mcp-ui-protocol-config" metadata field.
//
// Example client metadata:
//
//	{
//	  "metadata": {
//	    "mcp-ui-protocol": "appssdk",
//	    "mcp-ui-protocol-config": {
//	      "timeout": 10000,
//	      "intentHandling": "prompt"
//	    }
//	  }
//	}
//
// Example usage:
//
//	func handleInitialize(req InitializeRequest) InitializeResponse {
//	    protocolConfig := mcpuiserver.ParseProtocolConfig(req.Params)
//	    // Store protocol config in session context
//	    sessions[sessionID] = &SessionContext{ProtocolConfig: protocolConfig}
//	    return InitializeResponse{...}
//	}
func ParseProtocolConfig(initializeParams map[string]interface{}) *ProtocolConfig {
	protocol := ParseProtocolFromInitialize(initializeParams)
	config := &ProtocolConfig{
		Type: protocol,
	}

	// Extract optional protocol-specific config
	if metadata, ok := initializeParams["metadata"].(map[string]interface{}); ok {
		if protocolConfig, ok := metadata["mcp-ui-protocol-config"].(map[string]interface{}); ok {
			config.Config = protocolConfig
		}
	}

	return config
}
