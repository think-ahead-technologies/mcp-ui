package mcpuiserver

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocolVersion(t *testing.T) {
	// Ensure protocol version is set correctly
	assert.Equal(t, "2025-11-21", ProtocolVersion)
}

func TestResourceURIMetaKey(t *testing.T) {
	// Ensure resource URI meta key matches spec
	assert.Equal(t, "ui/resourceUri", ResourceURIMetaKey)
}

func TestProtocolMessageTypeConstants(t *testing.T) {
	// Verify message type constants match TypeScript SDK
	assert.Equal(t, ProtocolMessageType("tool"), MessageTypeToolCall)
	assert.Equal(t, ProtocolMessageType("prompt"), MessageTypePrompt)
	assert.Equal(t, ProtocolMessageType("link"), MessageTypeLink)
	assert.Equal(t, ProtocolMessageType("intent"), MessageTypeIntent)
	assert.Equal(t, ProtocolMessageType("notify"), MessageTypeNotify)
	assert.Equal(t, ProtocolMessageType("ui-lifecycle-iframe-ready"), MessageTypeLifecycleReady)
	assert.Equal(t, ProtocolMessageType("ui-size-change"), MessageTypeSizeChange)
	assert.Equal(t, ProtocolMessageType("ui-request-data"), MessageTypeRequestData)
	assert.Equal(t, ProtocolMessageType("ui-request-render-data"), MessageTypeRequestRenderData)
	assert.Equal(t, ProtocolMessageType("ui-lifecycle-iframe-render-data"), MessageTypeLifecycleRenderData)
	assert.Equal(t, ProtocolMessageType("ui-message-received"), MessageTypeMessageReceived)
	assert.Equal(t, ProtocolMessageType("ui-message-response"), MessageTypeMessageResponse)
}

func TestDisplayModeConstants(t *testing.T) {
	// Verify display mode constants
	assert.Equal(t, DisplayMode("inline"), DisplayModeInline)
	assert.Equal(t, DisplayMode("pip"), DisplayModePIP)
	assert.Equal(t, DisplayMode("fullscreen"), DisplayModeFullscreen)
}

func TestRenderDataSerialization(t *testing.T) {
	renderData := RenderData{
		Locale:      "en-US",
		Theme:       "dark",
		DisplayMode: DisplayModeInline,
		MaxHeight:   600,
		ToolInput: map[string]interface{}{
			"query": "test",
		},
	}

	jsonBytes, err := json.Marshal(renderData)
	assert.NoError(t, err)

	var decoded RenderData
	err = json.Unmarshal(jsonBytes, &decoded)
	assert.NoError(t, err)

	assert.Equal(t, "en-US", decoded.Locale)
	assert.Equal(t, "dark", decoded.Theme)
	assert.Equal(t, DisplayModeInline, decoded.DisplayMode)
	assert.Equal(t, 600, decoded.MaxHeight)
	assert.Equal(t, "test", decoded.ToolInput["query"])
}

func TestNewLifecycleReadyMessage(t *testing.T) {
	msgID := "msg-123"
	msg := NewLifecycleReadyMessage(&msgID)

	assert.Equal(t, MessageTypeLifecycleReady, msg.Type)
	assert.Equal(t, &msgID, msg.MessageID)
	assert.NotNil(t, msg.Payload)
}

func TestNewLifecycleReadyMessageNoID(t *testing.T) {
	msg := NewLifecycleReadyMessage(nil)

	assert.Equal(t, MessageTypeLifecycleReady, msg.Type)
	assert.Nil(t, msg.MessageID)
	assert.NotNil(t, msg.Payload)
}

func TestNewSizeChangeMessage(t *testing.T) {
	width := 800
	height := 600
	msgID := "msg-456"

	msg := NewSizeChangeMessage(&width, &height, &msgID)

	assert.Equal(t, MessageTypeSizeChange, msg.Type)
	assert.Equal(t, &msgID, msg.MessageID)
	assert.Equal(t, &width, msg.Payload.Width)
	assert.Equal(t, &height, msg.Payload.Height)
}

func TestNewSizeChangeMessageNilDimensions(t *testing.T) {
	msgID := "msg-456"

	msg := NewSizeChangeMessage(nil, nil, &msgID)

	assert.Equal(t, MessageTypeSizeChange, msg.Type)
	assert.Nil(t, msg.Payload.Width)
	assert.Nil(t, msg.Payload.Height)
}

func TestNewRequestDataMessage(t *testing.T) {
	msgID := "msg-789"
	params := map[string]interface{}{
		"userId": "123",
	}

	msg := NewRequestDataMessage("getUserData", params, msgID)

	assert.Equal(t, MessageTypeRequestData, msg.Type)
	assert.Equal(t, msgID, msg.MessageID)
	assert.Equal(t, "getUserData", msg.Payload.RequestType)
	assert.Equal(t, "123", msg.Payload.Params["userId"])
}

func TestNewRequestRenderDataMessage(t *testing.T) {
	msgID := "msg-101"
	msg := NewRequestRenderDataMessage(&msgID)

	assert.Equal(t, MessageTypeRequestRenderData, msg.Type)
	assert.Equal(t, &msgID, msg.MessageID)
	assert.NotNil(t, msg.Payload)
}

func TestNewRenderDataMessage(t *testing.T) {
	renderData := RenderData{
		Locale:      "en-US",
		Theme:       "dark",
		DisplayMode: DisplayModeInline,
		MaxHeight:   600,
	}
	msgID := "msg-202"

	msg := NewRenderDataMessage(renderData, &msgID)

	assert.Equal(t, MessageTypeLifecycleRenderData, msg.Type)
	assert.Equal(t, &msgID, msg.MessageID)
	assert.Equal(t, "en-US", msg.Payload.RenderData.Locale)
	assert.Equal(t, "dark", msg.Payload.RenderData.Theme)
	assert.Equal(t, DisplayModeInline, msg.Payload.RenderData.DisplayMode)
	assert.Equal(t, 600, msg.Payload.RenderData.MaxHeight)
}

func TestNewMessageReceivedMessage(t *testing.T) {
	acknowledgedID := "msg-original"
	msgID := "msg-ack"

	msg := NewMessageReceivedMessage(acknowledgedID, &msgID)

	assert.Equal(t, MessageTypeMessageReceived, msg.Type)
	assert.Equal(t, &msgID, msg.MessageID)
	assert.Equal(t, acknowledgedID, msg.Payload.MessageID)
}

func TestNewMessageResponseMessage(t *testing.T) {
	requestID := "msg-request"
	msgID := "msg-response"
	response := map[string]interface{}{"status": "ok"}

	msg := NewMessageResponseMessage(requestID, response, nil, &msgID)

	assert.Equal(t, MessageTypeMessageResponse, msg.Type)
	assert.Equal(t, &msgID, msg.MessageID)
	assert.Equal(t, requestID, msg.Payload.MessageID)
	assert.Equal(t, response, msg.Payload.Response)
	assert.Nil(t, msg.Payload.Error)
}

func TestNewMessageResponseMessageWithError(t *testing.T) {
	requestID := "msg-request"
	msgID := "msg-response"
	err := map[string]interface{}{"code": "ERR001", "message": "Something failed"}

	msg := NewMessageResponseMessage(requestID, nil, err, &msgID)

	assert.Equal(t, MessageTypeMessageResponse, msg.Type)
	assert.Equal(t, &msgID, msg.MessageID)
	assert.Equal(t, requestID, msg.Payload.MessageID)
	assert.Nil(t, msg.Payload.Response)
	assert.Equal(t, err, msg.Payload.Error)
}

func TestProtocolMessageSerialization(t *testing.T) {
	// Test that all message types can be serialized to JSON
	msgID := "test-123"
	width := 800

	messages := []interface{}{
		NewLifecycleReadyMessage(&msgID),
		NewSizeChangeMessage(&width, &width, &msgID),
		NewRequestDataMessage("test", nil, msgID),
		NewRequestRenderDataMessage(&msgID),
		NewRenderDataMessage(RenderData{}, &msgID),
		NewMessageReceivedMessage("orig", &msgID),
		NewMessageResponseMessage("req", nil, nil, &msgID),
	}

	for _, msg := range messages {
		jsonBytes, err := json.Marshal(msg)
		assert.NoError(t, err)
		assert.NotEmpty(t, jsonBytes)

		// Verify it can be parsed back
		var result map[string]interface{}
		err = json.Unmarshal(jsonBytes, &result)
		assert.NoError(t, err)
		assert.Contains(t, result, "type")
	}
}
