package mcpuiserver

// UIActionResult is an interface for all UI action result types
type UIActionResult interface {
	actionType() string
}

// Tool Call Action Result

// UIActionResultToolCallType represents a tool call action
type UIActionResultToolCallType struct {
	Type      string          `json:"type"`
	Payload   ToolCallPayload `json:"payload"`
	MessageID *string         `json:"messageId,omitempty"`
}

// ToolCallPayload contains the tool call parameters
type ToolCallPayload struct {
	ToolName string                 `json:"toolName"`
	Params   map[string]interface{} `json:"params"`
}

func (r UIActionResultToolCallType) actionType() string {
	return r.Type
}

// UIActionResultToolCall creates a tool call action result.
//
// Example:
//
//	result := UIActionResultToolCall("fetchData", map[string]interface{}{
//	    "query": "user stats",
//	})
func UIActionResultToolCall(toolName string, params map[string]interface{}) UIActionResultToolCallType {
	return UIActionResultToolCallType{
		Type: "tool",
		Payload: ToolCallPayload{
			ToolName: toolName,
			Params:   params,
		},
	}
}

// Prompt Action Result

// UIActionResultPromptType represents a prompt action
type UIActionResultPromptType struct {
	Type      string        `json:"type"`
	Payload   PromptPayload `json:"payload"`
	MessageID *string       `json:"messageId,omitempty"`
}

// PromptPayload contains the prompt text
type PromptPayload struct {
	Prompt string `json:"prompt"`
}

func (r UIActionResultPromptType) actionType() string {
	return r.Type
}

// UIActionResultPrompt creates a prompt action result.
//
// Example:
//
//	result := UIActionResultPrompt("Enter your API key")
func UIActionResultPrompt(prompt string) UIActionResultPromptType {
	return UIActionResultPromptType{
		Type: "prompt",
		Payload: PromptPayload{
			Prompt: prompt,
		},
	}
}

// Link Action Result

// UIActionResultLinkType represents a link action
type UIActionResultLinkType struct {
	Type      string      `json:"type"`
	Payload   LinkPayload `json:"payload"`
	MessageID *string     `json:"messageId,omitempty"`
}

// LinkPayload contains the URL
type LinkPayload struct {
	URL string `json:"url"`
}

func (r UIActionResultLinkType) actionType() string {
	return r.Type
}

// UIActionResultLink creates a link action result.
//
// Example:
//
//	result := UIActionResultLink("https://docs.example.com")
func UIActionResultLink(url string) UIActionResultLinkType {
	return UIActionResultLinkType{
		Type: "link",
		Payload: LinkPayload{
			URL: url,
		},
	}
}

// Intent Action Result

// UIActionResultIntentType represents an intent action
type UIActionResultIntentType struct {
	Type      string        `json:"type"`
	Payload   IntentPayload `json:"payload"`
	MessageID *string       `json:"messageId,omitempty"`
}

// IntentPayload contains the intent and parameters
type IntentPayload struct {
	Intent string                 `json:"intent"`
	Params map[string]interface{} `json:"params"`
}

func (r UIActionResultIntentType) actionType() string {
	return r.Type
}

// UIActionResultIntent creates an intent action result.
//
// Example:
//
//	result := UIActionResultIntent("showSettings", map[string]interface{}{
//	    "tab": "account",
//	})
func UIActionResultIntent(intent string, params map[string]interface{}) UIActionResultIntentType {
	return UIActionResultIntentType{
		Type: "intent",
		Payload: IntentPayload{
			Intent: intent,
			Params: params,
		},
	}
}

// Notification Action Result

// UIActionResultNotificationType represents a notification action
type UIActionResultNotificationType struct {
	Type      string              `json:"type"`
	Payload   NotificationPayload `json:"payload"`
	MessageID *string             `json:"messageId,omitempty"`
}

// NotificationPayload contains the notification message
type NotificationPayload struct {
	Message string `json:"message"`
}

func (r UIActionResultNotificationType) actionType() string {
	return r.Type
}

// UIActionResultNotification creates a notification action result.
//
// Example:
//
//	result := UIActionResultNotification("Data saved successfully!")
func UIActionResultNotification(message string) UIActionResultNotificationType {
	return UIActionResultNotificationType{
		Type: "notify",
		Payload: NotificationPayload{
			Message: message,
		},
	}
}
