package tools

import (
	"encoding/json"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/openai/openai-go/v3"
)

// ConvertOpenAPIToTools converts OpenAPI specifications to OpenAI function tool definitions
func ConvertOpenAPIToTools(specs []*openapi3.T) []openai.ChatCompletionToolUnionParam {
	tools := []openai.ChatCompletionToolUnionParam{}

	for _, spec := range specs {
		for _, pathItem := range spec.Paths.Map() {
			for _, operation := range pathItem.Operations() {
				if operation == nil || operation.OperationID == "" {
					continue
				}

				// Get description (summary preferred, fallback to description)
				desc := operation.Summary
				if desc == "" {
					desc = operation.Description
				}

				// Convert parameters - OpenAPI schema is already JSON Schema compatible
				params := buildParameters(operation)

				tools = append(tools, openai.ChatCompletionFunctionTool(openai.FunctionDefinitionParam{
					Name:        operation.OperationID,
					Description: openai.String(desc),
					Parameters:  params,
				}))
			}
		}
	}

	return tools
}

func buildParameters(operation *openapi3.Operation) openai.FunctionParameters {
	// Start with base structure
	params := openai.FunctionParameters{
		"type":       "object",
		"properties": map[string]any{},
	}

	// If no parameters and no body, return empty object schema
	if len(operation.Parameters) == 0 && operation.RequestBody == nil {
		return params
	}

	properties := map[string]any{}
	required := []string{}

	// Add parameters (query, path, header)
	for _, paramRef := range operation.Parameters {
		if paramRef.Value == nil || paramRef.Value.Schema == nil {
			continue
		}

		// Convert OpenAPI schema to map - they're already JSON Schema compatible
		schemaBytes, _ := paramRef.Value.Schema.MarshalJSON()
		var schemaDef map[string]interface{}
		json.Unmarshal(schemaBytes, &schemaDef)

		if paramRef.Value.Description != "" {
			schemaDef["description"] = paramRef.Value.Description
		}

		properties[paramRef.Value.Name] = schemaDef

		if paramRef.Value.Required {
			required = append(required, paramRef.Value.Name)
		}
	}

	// Add request body schema if present (application/json only)
	if operation.RequestBody != nil && operation.RequestBody.Value != nil {
		if content := operation.RequestBody.Value.Content.Get("application/json"); content != nil && content.Schema != nil {
			schemaBytes, _ := content.Schema.MarshalJSON()
			var bodySchema map[string]interface{}
			json.Unmarshal(schemaBytes, &bodySchema)

			// Merge body properties into parameters
			if props, ok := bodySchema["properties"].(map[string]interface{}); ok {
				for k, v := range props {
					properties[k] = v
				}
			}
			if reqs, ok := bodySchema["required"].([]interface{}); ok {
				for _, r := range reqs {
					if rStr, ok := r.(string); ok {
						required = append(required, rStr)
					}
				}
			}
		}
	}

	params["properties"] = properties
	if len(required) > 0 {
		params["required"] = required
	}

	return params
}
