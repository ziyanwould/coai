package cloudflare

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	adaptercommon "chat/adapter/common"
	"chat/globals"
	"chat/utils"
	"fmt"
	"strings"
)

// extractPureBase64 extracts pure base64 string from data URI (removes data:image/...;base64, prefix)
func extractPureBase64(base64Data string) string {
	// Remove data:image/...;base64, prefix if present
	if strings.Contains(base64Data, ";base64,") {
		parts := strings.Split(base64Data, ";base64,")
		if len(parts) == 2 {
			return parts[1]
		}
	}
	return base64Data
}

type ImageProps struct {
	Model     string
	Prompt    string
	Guidance  float32
	Seed      *int
	Height    int
	Width     int
	NumSteps  int
	InputImage string  // Base64 encoded input image for img2img
	MaskImage  string  // Base64 encoded mask image for inpainting
	Strength   float32 // Strength parameter for img2img (0.0-1.0)
	Proxy     globals.ProxyConfig
}

// CreateImageRequest calls Cloudflare Workers AI to generate an image
func (c *ChatInstance) CreateImageRequest(props ImageProps) (string, error) {
	isInpainting := strings.Contains(props.Model, "inpainting")

	// Set default values
	guidance := props.Guidance
	if guidance == 0 {
		guidance = 4.5
	}

	height, width := props.Height, props.Width
	if height == 0 && width == 0 {
		height, width = 1024, 1024
	}

	numSteps := props.NumSteps
	if numSteps == 0 {
		numSteps = 20
	}

	strength := props.Strength
	if strength == 0 && props.InputImage != "" {
		// Check if this model supports strength parameter
		if !strings.Contains(props.Model, "flux") &&
		   !strings.Contains(props.Model, "phoenix") &&
		   !strings.Contains(props.Model, "lucid-origin") {
			strength = 0.75 // Default strength for img2img
		}
	}

	var requestBodyBytes []byte
	var err error

	// Use the standard ImageRequest format for all models, including inpainting
	// Extract pure base64 for inpainting models if needed
	inputImage := props.InputImage
	maskImage := props.MaskImage

	if isInpainting {
		inputImage = extractPureBase64(props.InputImage)
		maskImage = extractPureBase64(props.MaskImage)
		fmt.Printf("[DEBUG] Using inpainting with image_b64 length: %d, mask_b64 length: %d\n", len(inputImage), len(maskImage))
	}

	requestBody := ImageRequest{
		Prompt:   props.Prompt,
		Guidance: guidance,
		Seed:     props.Seed,
		Height:   height,
		Width:    width,
		NumSteps: numSteps,
		ImageB64: inputImage,
		MaskB64:  maskImage,
		Strength: strength,
	}

	requestBodyBytes, err = json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// Make HTTP request with the appropriate request body
	buffer, err := utils.HttpRaw(
		c.GetImageEndpoint(props.Model),
		"POST",
		c.GetHeader(),
		bytes.NewReader(requestBodyBytes),
		[]globals.ProxyConfig{props.Proxy},
	)

	if err != nil {
		return "", fmt.Errorf("cloudflare api request failed: %v", err)
	}

	if len(buffer) == 0 {
		return "", fmt.Errorf("cloudflare error: empty response")
	}

	// Check if response is JSON format (Flux, Leonardo models)
	if len(buffer) > 0 && buffer[0] == '{' {
		// First try to parse as nested structure (Flux, Leonardo models)
		var nestedResp ImageResponseNested
		if err := json.Unmarshal(buffer, &nestedResp); err == nil {
			if nestedResp.Success && nestedResp.Result.Image != "" {
				// Return the base64 data directly from JSON response
				return nestedResp.Result.Image, nil
			}
			if !nestedResp.Success && len(nestedResp.Errors) > 0 {
				return "", fmt.Errorf("cloudflare error: %s", nestedResp.Errors[0].Message)
			}
		}

		// Try to parse as error response
		var errorResp ErrorResponse
		if err := json.Unmarshal(buffer, &errorResp); err == nil {
			if !errorResp.Success && len(errorResp.Errors) > 0 {
				return "", fmt.Errorf("cloudflare error: %s", errorResp.Errors[0].Message)
			}
		}

		// If we get here, it's an unexpected JSON format
		return "", fmt.Errorf("cloudflare error: unexpected JSON response format")
	}

	// Handle binary response (Phoenix model)
	base64Data := base64.StdEncoding.EncodeToString(buffer)
	return base64Data, nil
}

// CreateImage generates an image using Cloudflare Workers AI and returns markdown
func (c *ChatInstance) CreateImage(props *adaptercommon.ChatProps) (string, error) {
	prompt := c.GetLatestPrompt(props)
	if prompt == "" {
		return "", fmt.Errorf("empty prompt")
	}

	// Parse model to check if it's a supported image model
	if !c.IsImageModel(props.Model) {
		return "", fmt.Errorf("model %s is not supported for image generation", props.Model)
	}

	// Extract input image for img2img models and check for inappropriate usage
	var inputImage, maskImage string
	content, images := utils.ExtractImages(prompt, true)
	prompt = strings.TrimSpace(content) // Clean prompt text

	isInpainting := strings.Contains(props.Model, "inpainting")
	isImg2Img := globals.IsCloudflareImg2ImgModel(props.Model)

	// Debug logging for inpainting
	if isInpainting {
		fmt.Printf("[DEBUG] Inpainting model detected: %s\n", props.Model)
		fmt.Printf("[DEBUG] Extracted %d images from prompt\n", len(images))
		for i, img := range images {
			prefix := img
			if len(img) > 50 {
				prefix = img[:50] + "..."
			}
			fmt.Printf("[DEBUG] Image %d: %s\n", i, prefix)
		}
	}

	if len(images) > 0 {
		// User provided images
		if isImg2Img {
			// This is an img2img model, process the input image
			img, err := utils.NewImage(images[0])
			if err != nil {
				return "", fmt.Errorf("failed to process input image: %v", err)
			}
			inputImage = img.ToRawBase64()

			// For inpainting, we need a mask image
			if isInpainting {
				fmt.Printf("[DEBUG] Processing inpainting - found %d images\n", len(images))
				if len(images) > 1 {
					// User provided both image and mask
					fmt.Printf("[DEBUG] Processing mask image from images[1]\n")
					mask, err := utils.NewImage(images[1])
					if err != nil {
						return "", fmt.Errorf("failed to process mask image: %v", err)
					}
					maskImage = mask.ToRawBase64()
					fmt.Printf("[DEBUG] Mask image processed successfully, length: %d\n", len(maskImage))
				} else {
					// Create a simple center mask for inpainting
					// For now, we'll return an error asking for a mask
					fmt.Printf("[DEBUG] Insufficient images for inpainting: only %d found\n", len(images))
					return "", fmt.Errorf("修复模型 %s 需要提供遮罩图片。请上传两张图片：原图和遮罩图（黑白图，白色区域为要修复的部分）", props.Model)
				}
			}
		} else {
			// This is a regular text-to-image model but user provided an image
			return "", fmt.Errorf("模型 %s 不支持图片输入。请使用支持图生图的模型，如: @cf/runwayml/stable-diffusion-v1-5-img2img", props.Model)
		}
	} else if isImg2Img {
		// This is an img2img model but no image was provided
		if isInpainting {
			return "", fmt.Errorf("修复模型 %s 需要输入图片和遮罩。请上传图片附件", props.Model)
		} else {
			return "", fmt.Errorf("图生图模型 %s 需要输入图片。请上传图片附件", props.Model)
		}
	}

	// 根据不同模型调整参数
	guidance := float32(4.5)
	height, width := 1024, 1024
	steps := 20

	// 检查提示词中是否有尺寸要求
	lowerPrompt := strings.ToLower(prompt)
	if strings.Contains(lowerPrompt, "放大") || strings.Contains(lowerPrompt, "大图") ||
	   strings.Contains(lowerPrompt, "larger") || strings.Contains(lowerPrompt, "bigger") {
		height, width = 1536, 1536  // 更大尺寸
	}
	if strings.Contains(lowerPrompt, "宽屏") || strings.Contains(lowerPrompt, "横向") ||
	   strings.Contains(lowerPrompt, "landscape") {
		height, width = 1024, 1536  // 宽屏比例
	}
	if strings.Contains(lowerPrompt, "竖屏") || strings.Contains(lowerPrompt, "纵向") ||
	   strings.Contains(lowerPrompt, "portrait") {
		height, width = 1536, 1024  // 竖屏比例
	}

	// Flux模型使用更高质量的设置
	if strings.Contains(props.Model, "flux") {
		guidance = 3.5
		steps = 28
	}
	// DreamShaper使用LCM优化设置
	if strings.Contains(props.Model, "dreamshaper") || strings.Contains(props.Model, "lcm") {
		guidance = 7.5
		steps = 8
	}

	// Set default strength for img2img models
	strength := float32(0.75)

	// Different models have different strength requirements
	if globals.IsCloudflareImg2ImgModel(props.Model) {
		if strings.Contains(props.Model, "inpainting") {
			strength = 0.85 // Higher strength for inpainting
		}
	}

	// Some models don't support strength parameter or have different requirements
	supportsStrength := true
	if strings.Contains(props.Model, "flux") ||
	   strings.Contains(props.Model, "phoenix") ||
	   strings.Contains(props.Model, "lucid-origin") {
		supportsStrength = false
	}

	// Flux models that do support strength need it to be >= 1
	if strings.Contains(props.Model, "flux") && inputImage != "" {
		strength = 1.0 // Minimum value for Flux models
	}

	// Only set strength for models that support it
	requestStrength := float32(0)
	if supportsStrength && inputImage != "" {
		requestStrength = strength
	}

	base64Data, err := c.CreateImageRequest(ImageProps{
		Model:      props.Model,
		Prompt:     prompt,
		Guidance:   guidance,
		Height:     height,
		Width:      width,
		NumSteps:   steps,
		InputImage: inputImage,
		MaskImage:  maskImage,
		Strength:   requestStrength,
		Proxy:      props.Proxy,
	})

	if err != nil {
		if strings.Contains(err.Error(), "safety") || strings.Contains(err.Error(), "inappropriate") {
			return err.Error(), nil
		}
		return "", err
	}

	if base64Data == "" {
		return "", fmt.Errorf("no image generated")
	}

	// Cloudflare returns binary image data, we need to create a proper data URI
	imageDataURI := fmt.Sprintf("data:image/png;base64,%s", base64Data)
	return utils.GetImageMarkdown(imageDataURI), nil
}

// GetLatestPrompt extracts the latest user prompt from messages
func (c *ChatInstance) GetLatestPrompt(props *adaptercommon.ChatProps) string {
	for i := len(props.Message) - 1; i >= 0; i-- {
		if props.Message[i].Role == globals.User {
			return props.Message[i].Content
		}
	}
	return ""
}

// IsImageModel checks if the model supports image generation
func (c *ChatInstance) IsImageModel(model string) bool {
	return globals.IsCloudflareImageModel(model)
}