// Package mcpuiserver provides protocol message types and constructors for MCP-UI communication.
package mcpuiserver

// MCPUILifecycleReadyMessage indicates the widget is ready
type MCPUILifecycleReadyMessage struct {
	Type      ProtocolMessageType    `json:"type"`
	MessageID *string                `json:"messageId,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

// MCPUISizeChangeMessage requests a size change
type MCPUISizeChangeMessage struct {
	Type      ProtocolMessageType `json:"type"`
	MessageID *string             `json:"messageId,omitempty"`
	Payload   SizeChangePayload   `json:"payload"`
}

// SizeChangePayload contains width and height dimensions
type SizeChangePayload struct {
	Width  *int `json:"width,omitempty"`
	Height *int `json:"height,omitempty"`
}

// MCPUIRequestDataMessage requests data from the host
type MCPUIRequestDataMessage struct {
	Type      ProtocolMessageType `json:"type"`
	MessageID string              `json:"messageId"`
	Payload   RequestDataPayload  `json:"payload"`
}

// RequestDataPayload contains the request parameters
type RequestDataPayload struct {
	RequestType string                 `json:"requestType"`
	Params      map[string]interface{} `json:"params,omitempty"`
}

// MCPUIRequestRenderDataMessage requests render data
type MCPUIRequestRenderDataMessage struct {
	Type      ProtocolMessageType    `json:"type"`
	MessageID *string                `json:"messageId,omitempty"`
	Payload   map[string]interface{} `json:"payload,omitempty"`
}

// MCPUIRenderDataMessage delivers render data to the widget
type MCPUIRenderDataMessage struct {
	Type      ProtocolMessageType `json:"type"`
	MessageID *string             `json:"messageId,omitempty"`
	Payload   RenderDataPayload   `json:"payload"`
}

// RenderDataPayload wraps the RenderData structure
type RenderDataPayload struct {
	RenderData RenderData `json:"renderData"`
}

// MCPUIMessageReceivedMessage acknowledges message receipt
type MCPUIMessageReceivedMessage struct {
	Type      ProtocolMessageType    `json:"type"`
	MessageID *string                `json:"messageId,omitempty"`
	Payload   MessageReceivedPayload `json:"payload"`
}

// MessageReceivedPayload contains the acknowledged message ID
type MessageReceivedPayload struct {
	MessageID string `json:"messageId"`
}

// MCPUIMessageResponseMessage delivers a response to a request
type MCPUIMessageResponseMessage struct {
	Type      ProtocolMessageType    `json:"type"`
	MessageID *string                `json:"messageId,omitempty"`
	Payload   MessageResponsePayload `json:"payload"`
}

// MessageResponsePayload contains the response or error
type MessageResponsePayload struct {
	MessageID string      `json:"messageId"`
	Response  interface{} `json:"response,omitempty"`
	Error     interface{} `json:"error,omitempty"`
}

// NewLifecycleReadyMessage creates a lifecycle ready message
func NewLifecycleReadyMessage(messageID *string) *MCPUILifecycleReadyMessage {
	return &MCPUILifecycleReadyMessage{
		Type:      MessageTypeLifecycleReady,
		MessageID: messageID,
		Payload:   make(map[string]interface{}),
	}
}

// NewSizeChangeMessage creates a size change message
func NewSizeChangeMessage(width, height *int, messageID *string) *MCPUISizeChangeMessage {
	return &MCPUISizeChangeMessage{
		Type:      MessageTypeSizeChange,
		MessageID: messageID,
		Payload: SizeChangePayload{
			Width:  width,
			Height: height,
		},
	}
}

// NewRequestDataMessage creates a request data message
func NewRequestDataMessage(requestType string, params map[string]interface{}, messageID string) *MCPUIRequestDataMessage {
	return &MCPUIRequestDataMessage{
		Type:      MessageTypeRequestData,
		MessageID: messageID,
		Payload: RequestDataPayload{
			RequestType: requestType,
			Params:      params,
		},
	}
}

// NewRequestRenderDataMessage creates a request render data message
func NewRequestRenderDataMessage(messageID *string) *MCPUIRequestRenderDataMessage {
	return &MCPUIRequestRenderDataMessage{
		Type:      MessageTypeRequestRenderData,
		MessageID: messageID,
		Payload:   make(map[string]interface{}),
	}
}

// NewRenderDataMessage creates a render data message
func NewRenderDataMessage(renderData RenderData, messageID *string) *MCPUIRenderDataMessage {
	return &MCPUIRenderDataMessage{
		Type:      MessageTypeLifecycleRenderData,
		MessageID: messageID,
		Payload: RenderDataPayload{
			RenderData: renderData,
		},
	}
}

// NewMessageReceivedMessage creates a message received acknowledgment
func NewMessageReceivedMessage(acknowledgedMessageID string, messageID *string) *MCPUIMessageReceivedMessage {
	return &MCPUIMessageReceivedMessage{
		Type:      MessageTypeMessageReceived,
		MessageID: messageID,
		Payload: MessageReceivedPayload{
			MessageID: acknowledgedMessageID,
		},
	}
}

// NewMessageResponseMessage creates a message response
func NewMessageResponseMessage(requestMessageID string, response interface{}, err interface{}, messageID *string) *MCPUIMessageResponseMessage {
	return &MCPUIMessageResponseMessage{
		Type:      MessageTypeMessageResponse,
		MessageID: messageID,
		Payload: MessageResponsePayload{
			MessageID: requestMessageID,
			Response:  response,
			Error:     err,
		},
	}
}
