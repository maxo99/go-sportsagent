package tools

import (
	"github.com/openai/openai-go/v3"
)

func GetTools() []openai.ChatCompletionToolUnionParam {
	return []openai.ChatCompletionToolUnionParam{
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_roto_data",
			Description: openai.String("Retrieves sports data from the rotoreader service"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]string{
						"type":        "string",
						"description": "The query to retrieve data",
					},
				},
			},
		}),
		openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
			Name:        "get_odds_data",
			Description: openai.String("Retrieves betting odds from the oddstracker service"),
			Parameters: openai.FunctionParameters{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]string{
						"type":        "string",
						"description": "The query to retrieve odds",
					},
				},
			},
		}),
	}
}
