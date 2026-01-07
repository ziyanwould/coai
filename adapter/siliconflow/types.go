package siliconflow

import (
	adaptercommon "chat/adapter/common"
	"chat/globals"
)

type ChatInstance struct {
	Endpoint string
	Token    string
	Model    string
}

// Image generation request for SiliconFlow API
type ImageRequest struct {
	Model              string      `json:"model"`
	Prompt             string      `json:"prompt"`
	NegativePrompt     string      `json:"negative_prompt,omitempty"`
	ImageSize          string      `json:"image_size,omitempty"`
	BatchSize          int         `json:"batch_size,omitempty"`
	Seed               *int64      `json:"seed,omitempty"`
	NumInferenceSteps  int         `json:"num_inference_steps,omitempty"`
	GuidanceScale      float32     `json:"guidance_scale,omitempty"`
	Cfg                float32     `json:"cfg,omitempty"` // Qwen model specific parameter
	Image              string      `json:"image,omitempty"`     // Base64 encoded image for editing
	User               interface{} `json:"user,omitempty"`      // User identifier for tracking
	UserIP             string      `json:"user_ip,omitempty"`   // User IP for tracking
}

// Image generation response from SiliconFlow API
type ImageResponse struct {
	Images []struct {
		URL string `json:"url"`
	} `json:"images"`
	Timings struct {
		Inference float64 `json:"inference"`
	} `json:"timings"`
	Seed     int64 `json:"seed"`
	SharedId string `json:"shared_id"`
	Data     []struct {
		URL string `json:"url"`
	} `json:"data"`
	Created int64 `json:"created"`
}

// Error response structure
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

func (c *ChatInstance) GetEndpoint() string {
	return c.Endpoint
}

func (c *ChatInstance) GetToken() string {
	return c.Token
}

func (c *ChatInstance) GetModel() string {
	return c.Model
}

func (c *ChatInstance) GetImageEndpoint() string {
	return c.GetEndpoint() + "/images/generations"
}

func (c *ChatInstance) GetHeader() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + c.GetToken(),
		"Content-Type":  "application/json",
	}
}

func NewChatInstanceFromConfig(config globals.ChannelConfig) adaptercommon.Factory {
	return &ChatInstance{
		Endpoint: config.GetEndpoint(),
		Token:    config.GetRandomSecret(),
		Model:    config.GetModelReflect(""),
	}
}