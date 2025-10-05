package services

import (
	"context"
	"encoding/json"
	"fmt"

	"sportsagent/internal/clients"
	"sportsagent/internal/tools"

	"github.com/openai/openai-go/v3"
)

type AgentService struct {
	client      *openai.Client
	rotoreader  *clients.RotoReaderClient
	oddstracker *clients.OddsTrackerClient
	tools       []openai.ChatCompletionToolUnionParam
}

func NewAgentService() *AgentService {
	client := openai.NewClient()

	return &AgentService{
		client:      &client,
		rotoreader:  clients.NewRotoReaderClient(),
		oddstracker: clients.NewOddsTrackerClient(),
		tools:       tools.GetTools(),
	}
}

func (s *AgentService) ProcessQuery(ctx context.Context, query string) (string, error) {
	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(query),
	}

	response, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModelGPT4o,
		Messages: messages,
		Tools:    s.tools,
	})
	if err != nil {
		return "", err
	}

	choice := response.Choices[0]

	if choice.FinishReason == "tool_calls" {
		for _, toolCall := range choice.Message.ToolCalls {
			result := s.executeToolCall(ctx, toolCall)

			messages = append(messages, choice.Message.ToParam())
			messages = append(messages, openai.ToolMessage(result, toolCall.ID))
		}

		finalResponse, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
			Model:    openai.ChatModelGPT4o,
			Messages: messages,
		})
		if err != nil {
			return "", err
		}

		return finalResponse.Choices[0].Message.Content, nil
	}

	return choice.Message.Content, nil
}

func (s *AgentService) executeToolCall(ctx context.Context, toolCall openai.ChatCompletionMessageToolCallUnion) string {
	switch toolCall.Type {
	case "function":
		functionToolCall := toolCall.AsFunction()

		var args map[string]interface{}
		json.Unmarshal([]byte(functionToolCall.Function.Arguments), &args)

		switch functionToolCall.Function.Name {
		case "get_roto_data":
			data, err := s.rotoreader.GetData(ctx, args)
			if err != nil {
				return fmt.Sprintf("error: %v", err)
			}
			return data
		case "get_odds_data":
			data, err := s.oddstracker.GetData(ctx, args)
			if err != nil {
				return fmt.Sprintf("error: %v", err)
			}
			return data
		default:
			return "unknown function"
		}
	default:
		return "unsupported tool type"
	}
}
