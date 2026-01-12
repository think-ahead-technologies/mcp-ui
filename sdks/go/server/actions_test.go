package mcpuiserver

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUIActionResultToolCall(t *testing.T) {
	tests := []struct {
		name     string
		toolName string
		params   map[string]interface{}
		want     UIActionResultToolCallType
	}{
		{
			name:     "simple tool call",
			toolName: "testTool",
			params: map[string]interface{}{
				"param1": "value1",
			},
			want: UIActionResultToolCallType{
				Type: "tool",
				Payload: ToolCallPayload{
					ToolName: "testTool",
					Params: map[string]interface{}{
						"param1": "value1",
					},
				},
			},
		},
		{
			name:     "tool call with complex params",
			toolName: "complexTool",
			params: map[string]interface{}{
				"str":  "value",
				"num":  42,
				"bool": true,
				"obj": map[string]interface{}{
					"nested": "data",
				},
			},
			want: UIActionResultToolCallType{
				Type: "tool",
				Payload: ToolCallPayload{
					ToolName: "complexTool",
					Params: map[string]interface{}{
						"str":  "value",
						"num":  42,
						"bool": true,
						"obj": map[string]interface{}{
							"nested": "data",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UIActionResultToolCall(tt.toolName, tt.params)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUIActionResultPrompt(t *testing.T) {
	tests := []struct {
		name   string
		prompt string
		want   UIActionResultPromptType
	}{
		{
			name:   "simple prompt",
			prompt: "Enter your name",
			want: UIActionResultPromptType{
				Type: "prompt",
				Payload: PromptPayload{
					Prompt: "Enter your name",
				},
			},
		},
		{
			name:   "complex prompt",
			prompt: "Please provide the following information:\n1. Name\n2. Email",
			want: UIActionResultPromptType{
				Type: "prompt",
				Payload: PromptPayload{
					Prompt: "Please provide the following information:\n1. Name\n2. Email",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UIActionResultPrompt(tt.prompt)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUIActionResultLink(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want UIActionResultLinkType
	}{
		{
			name: "simple URL",
			url:  "https://example.com",
			want: UIActionResultLinkType{
				Type: "link",
				Payload: LinkPayload{
					URL: "https://example.com",
				},
			},
		},
		{
			name: "URL with path and query",
			url:  "https://docs.example.com/guide?section=api",
			want: UIActionResultLinkType{
				Type: "link",
				Payload: LinkPayload{
					URL: "https://docs.example.com/guide?section=api",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UIActionResultLink(tt.url)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUIActionResultIntent(t *testing.T) {
	tests := []struct {
		name   string
		intent string
		params map[string]interface{}
		want   UIActionResultIntentType
	}{
		{
			name:   "simple intent",
			intent: "showSettings",
			params: map[string]interface{}{
				"tab": "account",
			},
			want: UIActionResultIntentType{
				Type: "intent",
				Payload: IntentPayload{
					Intent: "showSettings",
					Params: map[string]interface{}{
						"tab": "account",
					},
				},
			},
		},
		{
			name:   "complex intent",
			intent: "navigate",
			params: map[string]interface{}{
				"screen": "profile",
				"userId": 123,
				"options": map[string]interface{}{
					"animated": true,
				},
			},
			want: UIActionResultIntentType{
				Type: "intent",
				Payload: IntentPayload{
					Intent: "navigate",
					Params: map[string]interface{}{
						"screen": "profile",
						"userId": 123,
						"options": map[string]interface{}{
							"animated": true,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UIActionResultIntent(tt.intent, tt.params)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUIActionResultNotification(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    UIActionResultNotificationType
	}{
		{
			name:    "simple notification",
			message: "Data saved successfully!",
			want: UIActionResultNotificationType{
				Type: "notify",
				Payload: NotificationPayload{
					Message: "Data saved successfully!",
				},
			},
		},
		{
			name:    "error notification",
			message: "Failed to process request",
			want: UIActionResultNotificationType{
				Type: "notify",
				Payload: NotificationPayload{
					Message: "Failed to process request",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UIActionResultNotification(tt.message)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUIActionResultJSON(t *testing.T) {
	tests := []struct {
		name     string
		result   interface{}
		wantJSON string
	}{
		{
			name: "tool call serialization",
			result: UIActionResultToolCall("testTool", map[string]interface{}{
				"param1": "value1",
			}),
			wantJSON: `{
				"type": "tool",
				"payload": {
					"toolName": "testTool",
					"params": {
						"param1": "value1"
					}
				}
			}`,
		},
		{
			name:   "prompt serialization",
			result: UIActionResultPrompt("Enter your name"),
			wantJSON: `{
				"type": "prompt",
				"payload": {
					"prompt": "Enter your name"
				}
			}`,
		},
		{
			name:   "link serialization",
			result: UIActionResultLink("https://example.com"),
			wantJSON: `{
				"type": "link",
				"payload": {
					"url": "https://example.com"
				}
			}`,
		},
		{
			name: "intent serialization",
			result: UIActionResultIntent("showSettings", map[string]interface{}{
				"tab": "account",
			}),
			wantJSON: `{
				"type": "intent",
				"payload": {
					"intent": "showSettings",
					"params": {
						"tab": "account"
					}
				}
			}`,
		},
		{
			name:   "notification serialization",
			result: UIActionResultNotification("Data saved!"),
			wantJSON: `{
				"type": "notify",
				"payload": {
					"message": "Data saved!"
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.result)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.wantJSON, string(jsonBytes))
		})
	}
}

func TestUIActionResultWithMessageID(t *testing.T) {
	messageID := "msg-123"
	result := UIActionResultToolCall("testTool", map[string]interface{}{
		"param": "value",
	})
	result.MessageID = &messageID

	jsonBytes, err := json.Marshal(result)
	assert.NoError(t, err)

	wantJSON := `{
		"type": "tool",
		"payload": {
			"toolName": "testTool",
			"params": {
				"param": "value"
			}
		},
		"messageId": "msg-123"
	}`

	assert.JSONEq(t, wantJSON, string(jsonBytes))
}
