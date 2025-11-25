package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"sportsagent/internal/clients"
	"sportsagent/internal/tools"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type AgentService struct {
	llm         llms.Model
	rotoreader  *clients.RotoReaderClient
	oddstracker *clients.OddsTrackerClient
	tools       []llms.Tool
}

func NewAgentService() *AgentService {
	llm, err := openai.New(openai.WithModel("gpt-4o"))
	if err != nil {
		log.Fatalf("Failed to initialize OpenAI client: %v", err)
	}

	return &AgentService{
		llm:         llm,
		rotoreader:  clients.NewRotoReaderClient(),
		oddstracker: clients.NewOddsTrackerClient(),
		tools:       tools.GetTools(),
	}
}

func (s *AgentService) ProcessQuery(ctx context.Context, query string) (string, error) {
	log.Printf("AgentService: processing query (len=%d)", len(query))

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeHuman, query),
	}

	response, err := s.llm.GenerateContent(ctx, messages, llms.WithTools(s.tools))
	if err != nil {
		log.Printf("AgentService: generate content error: %v", err)
		return "", err
	}

	choice := response.Choices[0]
	log.Printf("AgentService: received completion (finishReason=%s, toolCalls=%d)", choice.StopReason, len(choice.ToolCalls))

	// If there are tool calls, execute them and continue the conversation
	if len(choice.ToolCalls) > 0 {
		// Append assistant's response with tool calls to history
		assistantResponse := llms.TextParts(llms.ChatMessageTypeAI, choice.Content)
		for _, tc := range choice.ToolCalls {
			assistantResponse.Parts = append(assistantResponse.Parts, tc)
		}
		messages = append(messages, assistantResponse)

		// Execute each tool call
		for _, toolCall := range choice.ToolCalls {
			log.Printf("AgentService: handling tool call id=%s type=%s name=%s", toolCall.ID, toolCall.Type, toolCall.FunctionCall.Name)
			
			result := s.executeToolCall(ctx, toolCall)

			// Append tool response to history
			messages = append(messages, llms.MessageContent{
				Role: llms.ChatMessageTypeTool,
				Parts: []llms.ContentPart{
					llms.ToolCallResponse{
						ToolCallID: toolCall.ID,
						Name:       toolCall.FunctionCall.Name,
						Content:    result,
					},
				},
			})
		}

		// Get final response from LLM
		finalResponse, err := s.llm.GenerateContent(ctx, messages)
		if err != nil {
			return "", err
		}

		return finalResponse.Choices[0].Content, nil
	}

	return choice.Content, nil
}

func (s *AgentService) executeToolCall(ctx context.Context, toolCall llms.ToolCall) string {
	if toolCall.Type != "function" {
		log.Printf("AgentService: unsupported tool type %s", toolCall.Type)
		return "unsupported tool type"
	}

	functionName := toolCall.FunctionCall.Name
	log.Printf("AgentService: executing function tool %s", functionName)

	var args map[string]interface{}
	if err := json.Unmarshal([]byte(toolCall.FunctionCall.Arguments), &args); err != nil {
		log.Printf("AgentService: failed to unmarshal arguments for %s: %v", functionName, err)
		return fmt.Sprintf("error: invalid arguments: %v", err)
	}

	if metadata, ok := tools.GetToolMetadata(functionName); ok {
		log.Printf("AgentService: resolved service %s for tool %s (method=%s path=%s)", metadata.Service, functionName, metadata.Method, metadata.Path)
		switch metadata.Service {
		case tools.ServiceRotoReader:
			data, err := s.rotoreader.ExecuteOperation(ctx, functionName, args)
			if err != nil {
				log.Printf("AgentService: rotoreader error for %s: %v", functionName, err)
				return fmt.Sprintf("error: %v", err)
			}
			return data
		case tools.ServiceOddsTracker:
			data, err := s.oddstracker.ExecuteOperation(ctx, functionName, args)
			if err != nil {
				log.Printf("AgentService: oddstracker error for %s: %v", functionName, err)
				return fmt.Sprintf("error: %v", err)
			}
			return data
		default:
			log.Printf("AgentService: unsupported service %s for tool %s", metadata.Service, functionName)
			return fmt.Sprintf("error: unsupported service %s", metadata.Service)
		}
	}

	log.Printf("AgentService: unknown function tool %s", functionName)
	return "unknown function"
}
