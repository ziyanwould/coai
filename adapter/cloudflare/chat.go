package cloudflare

import (
	adaptercommon "chat/adapter/common"
	"chat/globals"
	"fmt"
)

// CreateStreamChatRequest handles both text chat and image generation requests
func (c *ChatInstance) CreateStreamChatRequest(props *adaptercommon.ChatProps, hook globals.Hook) error {
	// Check if this is an image generation request
	if c.IsImageModel(props.Model) {
		return c.handleImageGeneration(props, hook)
	}

	// For now, Cloudflare adapter only supports image generation
	// You can extend this to support text chat models if needed
	return fmt.Errorf("text chat is not supported by Cloudflare adapter yet, only image generation models are supported")
}

// handleImageGeneration processes image generation requests
func (c *ChatInstance) handleImageGeneration(props *adaptercommon.ChatProps, hook globals.Hook) error {
	// Send start event
	if hook != nil {
		if err := hook(&globals.Chunk{Content: "[start]"}); err != nil {
			return err
		}
	}

	// Generate image
	result, err := c.CreateImage(props)
	if err != nil {
		// Send error event
		if hook != nil {
			if hookErr := hook(&globals.Chunk{Content: fmt.Sprintf("Error: %v", err)}); hookErr != nil {
				return hookErr
			}
		}
		return err
	}

	// Send image result
	if hook != nil {
		if err := hook(&globals.Chunk{Content: result}); err != nil {
			return err
		}
	}

	// Send end event
	if hook != nil {
		if err := hook(&globals.Chunk{Content: "[end]"}); err != nil {
			return err
		}
	}

	return nil
}

// GetSupportedModels returns the list of supported Cloudflare AI models
func GetSupportedModels() []string {
	return globals.CloudflareImageModels
}