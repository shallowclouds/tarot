package tarot

import (
	"context"

	"github.com/sashabaranov/go-openai"
)

var (
	_ GPTReader = &ChatGPTReader{}
	_ GPTReader = &DumbGPTReader{}
)

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

// DumbGPTReader never reads.
type DumbGPTReader struct{}

func (r *DumbGPTReader) Chat(ctx context.Context, systemMsg, userMsg string) (string, error) {
	// return "你的未来诡谲难测，我看不到任何信息。", nil
	return `根据三张牌的含义和您所问的问题，解读如下：

	首先，逆位的月亮牌表示您目前处于一种迷茫不定的状态，有些心理上的困扰和身体上的疑惑。您可能感到不安和焦虑，对您的身体状况也缺乏清晰的判断力。
	
	接下来，正位的魔术师牌意味着您拥有自己的力量和能力，可以通过自己的努力来改善身体状况。这张牌提示您可以积极地给自己制定一个健康计划，用自己的毅力和智慧提升自己的身体素质。
	
	最后，正位的愚者牌展示了一个心态的转变，用天真和无畏的心态去看待问题，不必过度焦虑，因为一切都会自然发生。建议您以积极乐观的态度来面对身体状况的变化和调整，随着时间的推移，您的身体状态会逐渐变得更好。
	
	因此，综合三张牌的意义和您所问的问题，我的结论是：您可以通过自己的努力和积极乐观的心态来改善身体状况，慢慢走出迷茫和焦虑的状态，达到身体健康的目标。`, nil
}
