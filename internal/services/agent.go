package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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
	log.Printf("AgentService: processing query (len=%d)", len(query))

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(query),
	}

	response, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    openai.ChatModelGPT4o,
		Messages: messages,
		Tools:    s.tools,
	})
	if err != nil {
		log.Printf("AgentService: chat completion error: %v", err)
		return "", err
	}

	choice := response.Choices[0]
	log.Printf("AgentService: received completion (finishReason=%s, toolCalls=%d)", choice.FinishReason, len(choice.Message.ToolCalls))

	if choice.FinishReason == "tool_calls" {
		for _, toolCall := range choice.Message.ToolCalls {
			log.Printf("AgentService: handling tool call id=%s type=%s", toolCall.ID, toolCall.Type)
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
		log.Printf("AgentService: executing function tool %s", functionToolCall.Function.Name)

		var args map[string]interface{}
		json.Unmarshal([]byte(functionToolCall.Function.Arguments), &args)

		if serviceName, ok := tools.GetToolService(functionToolCall.Function.Name); ok {
			log.Printf("AgentService: resolved service %s for tool %s", serviceName, functionToolCall.Function.Name)
			switch serviceName {
			case tools.ServiceRotoReader:
				data, err := s.rotoreader.GetFeeds(ctx, args)
				if err != nil {
					log.Printf("AgentService: rotoreader error for %s: %v", functionToolCall.Function.Name, err)
					return fmt.Sprintf("error: %v", err)
				}
				return data
			case tools.ServiceOddsTracker:
				data, err := s.oddstracker.GetChanges(ctx, args)
				if err != nil {
					log.Printf("AgentService: oddstracker error for %s: %v", functionToolCall.Function.Name, err)
					return fmt.Sprintf("error: %v", err)
				}
				return data
			default:
				log.Printf("AgentService: unsupported service %s for tool %s", serviceName, functionToolCall.Function.Name)
				return fmt.Sprintf("error: unsupported service %s", serviceName)
			}
		}

		switch functionToolCall.Function.Name {
		case "get_roto_data":
			data, err := s.rotoreader.GetFeeds(ctx, args)
			if err != nil {
				return fmt.Sprintf("error: %v", err)
			}
			return data
		case "get_odds_data":
			data, err := s.oddstracker.GetChanges(ctx, args)
			if err != nil {
				return fmt.Sprintf("error: %v", err)
			}
			return data
		default:
			log.Printf("AgentService: unknown function tool %s", functionToolCall.Function.Name)
			return "unknown function"
		}
	default:
		log.Printf("AgentService: unsupported tool type %s", toolCall.Type)
		return "unsupported tool type"
	}
}
