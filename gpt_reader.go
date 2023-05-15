package tarot

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

var _ GPTReader = &ChatGPTReader{}

type ChatGPTReader struct {
	chatGPTCli *openai.Client
}

func NewChatGPTReader(chatGPTCli *openai.Client) *ChatGPTReader {
	return &ChatGPTReader{chatGPTCli: chatGPTCli}
}

func (r *ChatGPTReader) Chat(ctx context.Context, systemMsg, userMsg string) (string, error) {
	req := openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo0301,
		N:        1,
		Messages: make([]openai.ChatCompletionMessage, 0),
	}
	if len(systemMsg) != 0 {
		req.Messages = append(req.Messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemMsg,
		})
	}
	req.Messages = append(req.Messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMsg,
	})

	resp, err := r.chatGPTCli.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
