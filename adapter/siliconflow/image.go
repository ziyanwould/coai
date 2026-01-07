package siliconflow

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

type ImageProps struct {
	Model             string
	Prompt            string
	NegativePrompt    string
	ImageSize         string
	BatchSize         int
	Seed              *int64
	NumInferenceSteps int
	GuidanceScale     float32
	Cfg               float32 // Qwen model specific
	InputImage        string  // Base64 encoded input image for editing
	User              interface{}
	UserIP            string
	Proxy             globals.ProxyConfig
}

// CreateImageRequest calls SiliconFlow API to generate an image
func (c *ChatInstance) CreateImageRequest(props ImageProps) (string, error) {
	// Set default values
	imageSize := props.ImageSize
	if imageSize == "" {
		imageSize = "1024x1024"
	}

	batchSize := props.BatchSize
	if batchSize == 0 {
		batchSize = 1
	}

	numSteps := props.NumInferenceSteps
	if numSteps == 0 {
		numSteps = 20
	}

	guidanceScale := props.GuidanceScale
	if guidanceScale == 0 {
		guidanceScale = 7.5
	}

	// Build request body
	requestBody := ImageRequest{
		Model:             props.Model,
		Prompt:            props.Prompt,
		NegativePrompt:    props.NegativePrompt,
		ImageSize:         imageSize,
		BatchSize:         batchSize,
		Seed:              props.Seed,
		NumInferenceSteps: numSteps,
		GuidanceScale:     guidanceScale,
		User:              props.User,
		UserIP:            props.UserIP,
	}

	// Add input image for image editing models (needs full data URI format)
	if props.InputImage != "" {
		// For SiliconFlow, we need the full data URI format
		if strings.HasPrefix(props.InputImage, "data:") {
			requestBody.Image = props.InputImage
		} else {
			// Convert raw base64 to data URI
			requestBody.Image = fmt.Sprintf("data:image/png;base64,%s", props.InputImage)
		}
	}

	// Add Qwen-specific parameter if it's a Qwen model
	if strings.Contains(strings.ToLower(props.Model), "qwen") && props.Cfg > 0 {
		requestBody.Cfg = props.Cfg
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// Make HTTP request
	buffer, err := utils.HttpRaw(
		c.GetImageEndpoint(),
		"POST",
		c.GetHeader(),
		bytes.NewReader(requestBodyBytes),
		[]globals.ProxyConfig{props.Proxy},
	)

	if err != nil {
		return "", fmt.Errorf("siliconflow api request failed: %v", err)
	}

	if len(buffer) == 0 {
		return "", fmt.Errorf("siliconflow error: empty response")
	}

	// Parse response
	var response ImageResponse
	if err := json.Unmarshal(buffer, &response); err != nil {
		// Try to parse as error response
		var errorResp ErrorResponse
		if parseErr := json.Unmarshal(buffer, &errorResp); parseErr == nil {
			return "", fmt.Errorf("siliconflow error: %s", errorResp.Error.Message)
		}
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	if len(response.Images) == 0 {
		return "", fmt.Errorf("no image generated")
	}

	// Download the image from URL and convert to base64
	imageURL := response.Images[0].URL
	imageData, err := c.downloadImage(imageURL, props.Proxy)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %v", err)
	}

	return imageData, nil
}

// downloadImage downloads image from URL and returns base64 data
func (c *ChatInstance) downloadImage(url string, proxy globals.ProxyConfig) (string, error) {
	// Use utils.HttpRaw to download the image with proxy support
	imageBytes, err := utils.HttpRaw(url, "GET", map[string]string{}, nil, []globals.ProxyConfig{proxy})
	if err != nil {
		return "", err
	}

	// Convert to base64 data URI
	base64Data := base64.StdEncoding.EncodeToString(imageBytes)

	// For SiliconFlow, assume PNG format (most common)
	return fmt.Sprintf("data:image/png;base64,%s", base64Data), nil
}

// CreateImage generates an image using SiliconFlow API and returns markdown
func (c *ChatInstance) CreateImage(props *adaptercommon.ChatProps) (string, error) {
	prompt := c.GetLatestPrompt(props)
	if prompt == "" {
		return "", fmt.Errorf("empty prompt")
	}

	// Parse model to check if it's a supported image model
	if !c.IsImageModel(props.Model) {
		return "", fmt.Errorf("model %s is not supported for image generation", props.Model)
	}

	// Extract images from prompt
	content, images := utils.ExtractImages(prompt, true)
	prompt = strings.TrimSpace(content)

	isImg2Img := globals.IsSiliconFlowImg2ImgModel(props.Model)

	// Handle image input based on model type
	var inputImage string
	if len(images) > 0 {
		if isImg2Img {
			// This is an image editing model, process the input image
			fmt.Printf("[DEBUG] Processing input image for model: %s\n", props.Model)
			fmt.Printf("[DEBUG] Input image data length: %d characters\n", len(images[0]))
			previewLen := 100
			if len(images[0]) < previewLen {
				previewLen = len(images[0])
			}
			fmt.Printf("[DEBUG] Input image preview: %s...\n", images[0][:previewLen])

			img, err := utils.NewImage(images[0])
			if err != nil {
				// Provide more detailed error information
				fmt.Printf("[DEBUG] Image processing failed: %v\n", err)
				return "", fmt.Errorf("failed to process input image: %v\n\n请检查：\n1. 图片格式是否支持 (PNG, JPG, WEBP)\n2. 图片文件是否完整\n3. 图片大小是否合理", err)
			}
			// For SiliconFlow, we need the full data URI format
			inputImage = img.ToBase64()
			fmt.Printf("[DEBUG] Processed image data URI length: %d\n", len(inputImage))
		} else {
			// This is a text-to-image model but user provided an image
			return "", fmt.Errorf("模型 %s 不支持图片输入，仅支持文生图功能", props.Model)
		}
	} else if isImg2Img {
		// This is an image editing model but no image was provided
		return "", fmt.Errorf("图像编辑模型 %s 需要输入图片。请上传图片附件", props.Model)
	}

	// Set model-specific parameters
	var cfg float32
	var guidanceScale float32 = 7.5
	var steps int = 20
	var imageSize string = "1024x1024"

	// Qwen model specific settings
	if strings.Contains(strings.ToLower(props.Model), "qwen") {
		if strings.Contains(strings.ToLower(props.Model), "edit") {
			// Qwen-Image-Edit uses cfg parameter
			cfg = 4.0 // Default CFG for Qwen-Image-Edit
			if strings.Contains(strings.ToLower(prompt), "高质量") || strings.Contains(strings.ToLower(prompt), "细节") {
				cfg = 6.0 // Higher CFG for better quality
			}
		} else {
			// Regular Qwen-Image model
			cfg = 7.5 // Default CFG for Qwen models
			if strings.Contains(strings.ToLower(prompt), "高质量") || strings.Contains(strings.ToLower(prompt), "细节") {
				cfg = 9.0 // Higher CFG for better quality
			}
		}
	}

	// Kolors model specific settings
	if strings.Contains(strings.ToLower(props.Model), "kolors") {
		guidanceScale = 5.0 // Different guidance scale for Kolors
		steps = 25
	}

	// Check for size hints in prompt
	lowerPrompt := strings.ToLower(prompt)
	if strings.Contains(lowerPrompt, "宽屏") || strings.Contains(lowerPrompt, "横向") ||
	   strings.Contains(lowerPrompt, "landscape") {
		imageSize = "1280x960"
	}
	if strings.Contains(lowerPrompt, "竖屏") || strings.Contains(lowerPrompt, "纵向") ||
	   strings.Contains(lowerPrompt, "portrait") {
		imageSize = "960x1280"
	}

	imageData, err := c.CreateImageRequest(ImageProps{
		Model:             props.Model,
		Prompt:            prompt,
		ImageSize:         imageSize,
		BatchSize:         1,
		NumInferenceSteps: steps,
		GuidanceScale:     guidanceScale,
		Cfg:               cfg,
		InputImage:        inputImage,
		User:              props.User,
		UserIP:            props.Ip,
		Proxy:             props.Proxy,
	})

	if err != nil {
		if strings.Contains(err.Error(), "safety") || strings.Contains(err.Error(), "inappropriate") {
			return err.Error(), nil
		}
		return "", err
	}

	if imageData == "" {
		return "", fmt.Errorf("no image generated")
	}

	return utils.GetImageMarkdown(imageData), nil
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
	return globals.IsSiliconFlowImageModel(model)
}