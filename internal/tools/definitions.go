package tools

import (
	"github.com/openai/openai-go/v3"
)

func GetTools() []openai.ChatCompletionToolUnionParam {
	return []openai.ChatCompletionToolUnionParam{
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_roto_data",
			Description: openai.String("Get the latest sports news feed from rotoreader"),
			Parameters: openai.FunctionParameters{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		}),
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_odds_data",
			Description: openai.String("Get recent betting odds changes from oddstracker"),
			Parameters: openai.FunctionParameters{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		}),
	}
}
