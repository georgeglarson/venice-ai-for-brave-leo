package main

// Configuration settings
const (
	AppVersion        = "v1.0.2"
	VeniceLabel       = "Venice.AI"
	VeniceModel       = "llama-3.3-70b"
	VeniceEndpoint    = "https://api.venice.ai/api/v1/chat/completions"
	VeniceContextSize = 4000
)

// Windows API constants and types for MessageBox (Windows only)
// These are only used on Windows
var (
	MB_OK              uint32 = 0x00000000
	MB_OKCANCEL        uint32 = 0x00000001
	MB_YESNOCANCEL     uint32 = 0x00000003
	MB_ICONINFORMATION uint32 = 0x00000040
	MB_ICONQUESTION    uint32 = 0x00000020
	IDOK               int    = 1
	IDCANCEL           int    = 2
	IDYES              int    = 6
	IDNO               int    = 7
)

// CustomModel represents a Leo AI custom model configuration
type CustomModel struct {
	APIKey           string `json:"api_key"`
	ContextSize      int    `json:"context_size"`
	EndpointURL      string `json:"endpoint_url"`
	Key              string `json:"key"`
	Label            string `json:"label"`
	ModelRequestName string `json:"model_request_name"`
}

// BravePreferences represents the structure of Brave's Preferences.json file
type BravePreferences struct {
	Brave struct {
		AIChat struct {
			CustomModels   []CustomModel `json:"custom_models"`
			DefaultModelKey string       `json:"default_model_key"`
		} `json:"ai_chat"`
	} `json:"brave"`
}