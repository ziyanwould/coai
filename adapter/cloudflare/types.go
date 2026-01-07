package cloudflare

import (
	adaptercommon "chat/adapter/common"
	"chat/globals"
)

type ChatInstance struct {
	Endpoint string
	Token    string
	Model    string
}

// Image generation request for Cloudflare Workers AI
type ImageRequest struct {
	Prompt    string      `json:"prompt"`
	Guidance  float32     `json:"guidance,omitempty"`
	Seed      *int        `json:"seed,omitempty"`
	Height    int         `json:"height,omitempty"`
	Width     int         `json:"width,omitempty"`
	NumSteps  int         `json:"num_steps,omitempty"`
	ImageB64  string      `json:"image_b64,omitempty"` // Base64 encoded input image for img2img
	MaskB64   string      `json:"mask_b64,omitempty"`  // Base64 encoded mask image for inpainting
	Strength  float32     `json:"strength,omitempty"`  // Strength parameter for img2img (0.0-1.0)
	User      interface{} `json:"user,omitempty"`      // User identifier for tracking
	Userip    string      `json:"user_ip,omitempty"`   // User IP for tracking
}

// Inpainting request for Cloudflare Workers AI (special parameter names)
type InpaintingRequest struct {
	Prompt    string      `json:"prompt"`
	Guidance  float32     `json:"guidance,omitempty"`
	Seed      *int        `json:"seed,omitempty"`
	Height    int         `json:"height,omitempty"`
	Width     int         `json:"width,omitempty"`
	NumSteps  int         `json:"num_steps,omitempty"`
	ImageB64  string      `json:"image_b64"`       // Input image as base64 string (without data: prefix)
	MaskImage string      `json:"mask_image"`      // Mask image as base64 string (without data: prefix)
	User      interface{} `json:"user,omitempty"`  // User identifier for tracking
	Userip    string      `json:"user_ip,omitempty"` // User IP for tracking
}

// Image generation response from Cloudflare Workers AI
type ImageResponse struct {
	Success bool   `json:"success"`
	Result  string `json:"result"` // Base64 encoded image
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// Image generation response with nested result structure (for Flux and Leonardo models)
type ImageResponseNested struct {
	Success bool `json:"success"`
	Result  struct {
		Image string `json:"image"` // Base64 encoded image
	} `json:"result"`
	Errors []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// Error response structure
type ErrorResponse struct {
	Success bool `json:"success"`
	Errors  []struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"errors"`
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

func (c *ChatInstance) GetImageEndpoint(model string) string {
	return c.GetEndpoint() + "/" + model
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