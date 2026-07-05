package matchers

import (
	"encoding/json"
	"fmt"
	"strings"

	dtos "github.com/HiIamJeff67/notezy-backend/app/dtos"
	exceptions "github.com/HiIamJeff67/notezy-backend/app/exceptions"
)

type TemplateBlockMatcherInterface interface {
	MatchString(value string, values map[string]string) string
	MatchArborizedEditableBlock(block dtos.ArborizedEditableBlock, values map[string]string) (dtos.ArborizedEditableBlock, *exceptions.Exception)
}

type TemplateBlockMatcher struct{}

func NewTemplateBlockMatcher() TemplateBlockMatcherInterface {
	return TemplateBlockMatcher{}
}

func (m TemplateBlockMatcher) MatchString(value string, values map[string]string) string {
	if len(values) == 0 || !strings.Contains(value, "{{") {
		return value
	}
	matched := value
	for key, resolvedValue := range values {
		matched = strings.ReplaceAll(matched, "{{"+key+"}}", resolvedValue)
	}
	return matched
}

func (m TemplateBlockMatcher) MatchArborizedEditableBlock(
	block dtos.ArborizedEditableBlock,
	values map[string]string,
) (dtos.ArborizedEditableBlock, *exceptions.Exception) {
	matchedChildren := make([]dtos.ArborizedEditableBlock, len(block.Children))
	for index, child := range block.Children {
		matchedChild, exception := m.MatchArborizedEditableBlock(child, values)
		if exception != nil {
			return dtos.ArborizedEditableBlock{}, exception
		}
		matchedChildren[index] = matchedChild
	}
	block.Children = matchedChildren

	rawBlock, err := json.Marshal(block)
	if err != nil {
		return dtos.ArborizedEditableBlock{}, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	var blockMap map[string]any
	if err := json.Unmarshal(rawBlock, &blockMap); err != nil {
		return dtos.ArborizedEditableBlock{}, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	shouldMatch := false
	if props, ok := blockMap["props"].(map[string]any); ok {
		if template, ok := props["template"].(bool); ok && template {
			shouldMatch = true
		}
		delete(props, "template")
	}

	if shouldMatch {
		if props, exists := blockMap["props"]; exists {
			blockMap["props"] = m.matchJSONValue(props, values)
		}
		if content, exists := blockMap["content"]; exists {
			blockMap["content"] = m.matchJSONValue(content, values)
		}
	}

	rawMatchedBlock, err := json.Marshal(blockMap)
	if err != nil {
		return dtos.ArborizedEditableBlock{}, exceptions.RoutineTask.InvalidDto().WithOrigin(err)
	}

	var matchedBlock dtos.ArborizedEditableBlock
	if err := json.Unmarshal(rawMatchedBlock, &matchedBlock); err != nil {
		return dtos.ArborizedEditableBlock{}, exceptions.RoutineTask.InvalidDto().
			WithOrigin(fmt.Errorf("invalid matched template block: %w", err))
	}

	return matchedBlock, nil
}

func (m TemplateBlockMatcher) matchJSONValue(value any, values map[string]string) any {
	switch typed := value.(type) {
	case string:
		return m.MatchString(typed, values)
	case []any:
		matched := make([]any, len(typed))
		for index, item := range typed {
			matched[index] = m.matchJSONValue(item, values)
		}
		return matched
	case map[string]any:
		matched := make(map[string]any, len(typed))
		for key, item := range typed {
			matched[key] = m.matchJSONValue(item, values)
		}
		return matched
	default:
		return value
	}
}
