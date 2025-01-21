package openai

import (
	adaptercommon "chat/adapter/common"
	"chat/globals"
	"chat/utils"
	"errors"
	"fmt"
	"regexp"
)

func formatMessages(props *adaptercommon.ChatProps) interface{} {
	if globals.IsVisionModel(props.Model) || utils.IsCustomVisionModel(props.Model) {
		fmt.Printf("确认是视觉模型: %s\n", props.Model) // 添加打印模型名称的代码
		return utils.Each[globals.Message, Message](props.Message, func(message globals.Message) Message {
			if message.Role == globals.User {
				content, urls := utils.ExtractImages(message.Content, true)
				images := utils.EachNotNil[string, MessageContent](urls, func(url string) *MessageContent {
					obj, err := utils.NewImage(url)
					props.Buffer.AddImage(obj)
					if err != nil {
						globals.Info(fmt.Sprintf("无法处理图像: %s (来源: %s)", err.Error(), utils.Extract(url, 24, "...")))
					}
					fmt.Printf("正在处理图片URL: %s\n", url) // 打印正在处理的图片 URL
					return &MessageContent{
						Type: "image_url",
						ImageUrl: &ImageUrl{
							Url: url,
						},
					}
				})

				return Message{
					Role: message.Role,
					Content: utils.Prepend(images, MessageContent{
						Type: "text",
						Text: &content,
					}),
					Name:         message.Name,
					FunctionCall: message.FunctionCall,
					ToolCalls:    message.ToolCalls,
					ToolCallId:   message.ToolCallId,
				}
			}

			return Message{
				Role: message.Role,
				Content: MessageContents{
					MessageContent{
						Type: "text",
						Text: &message.Content,
					},
				},
				Name:         message.Name,
				FunctionCall: message.FunctionCall,
				ToolCalls:    message.ToolCalls,
				ToolCallId:   message.ToolCallId,
			}
		})
	}

	return props.Message
}

func processChatResponse(data string) *ChatStreamResponse {
	return utils.UnmarshalForm[ChatStreamResponse](data)
}

func processCompletionResponse(data string) *CompletionResponse {
	return utils.UnmarshalForm[CompletionResponse](data)
}

func processChatErrorResponse(data string) *ChatStreamErrorResponse {
	return utils.UnmarshalForm[ChatStreamErrorResponse](data)
}

func getChoices(form *ChatStreamResponse) *globals.Chunk {
	if len(form.Choices) == 0 {
		return &globals.Chunk{Content: ""}
	}

	choice := form.Choices[0].Delta

	return &globals.Chunk{
		Content:      choice.Content,
		ToolCall:     choice.ToolCalls,
		FunctionCall: choice.FunctionCall,
	}
}

func getCompletionChoices(form *CompletionResponse) string {
	if len(form.Choices) == 0 {
		return ""
	}

	return form.Choices[0].Text
}

func getRobustnessResult(chunk string) string {
	exp := `\"content\":\"(.*?)\"`
	compile, err := regexp.Compile(exp)
	if err != nil {
		return ""
	}

	matches := compile.FindStringSubmatch(chunk)
	if len(matches) > 1 {
		return utils.ProcessRobustnessChar(matches[1])
	} else {
		return ""
	}
}

func (c *ChatInstance) ProcessLine(data string, isCompletionType bool) (*globals.Chunk, error) {
	if isCompletionType {
		// openai legacy support
		if completion := processCompletionResponse(data); completion != nil {
			return &globals.Chunk{
				Content: getCompletionChoices(completion),
			}, nil
		}

		globals.Warn(fmt.Sprintf("openai 错误：无法解析完成响应: %s", data))
		return &globals.Chunk{Content: ""}, errors.New("解析错误：无法解析完成响应")
	}

	if form := processChatResponse(data); form != nil {
		return getChoices(form), nil
	}

	if form := processChatErrorResponse(data); form != nil {
		return &globals.Chunk{Content: ""}, errors.New(fmt.Sprintf("openai 错误: %s (类型: %s)", form.Error.Message, form.Error.Type))
	}

	globals.Warn(fmt.Sprintf("openai 错误：无法解析聊天补全响应: %s", data))
	return &globals.Chunk{Content: ""}, errors.New("解析错误：无法解析聊天补全响应")
}
