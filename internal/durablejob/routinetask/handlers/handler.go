package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

	routinetasktypes "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/types/routine-tasks"
	enums "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"
)

type PurposeHandler struct {
	HandlerFunc PurposeHandlerFunc
}

type PurposeHandlerFunc func(
	context.Context,
	routinetasktypes.RoutineTaskAssignment,
) (*routinetasktypes.PreparedRoutineTask, *exceptions.Exception)

func NewPurposeHandler(validator *validator.Validate) PurposeHandler {
	return PurposeHandler{
		HandlerFunc: func(
			ctx context.Context,
			assignment routinetasktypes.RoutineTaskAssignment,
		) (*routinetasktypes.PreparedRoutineTask, *exceptions.Exception) {
			return prepareAssignment(ctx, validator, assignment)
		},
	}
}

func prepareAssignment(
	_ context.Context,
	validator *validator.Validate,
	assignment routinetasktypes.RoutineTaskAssignment,
) (*routinetasktypes.PreparedRoutineTask, *exceptions.Exception) {
	if assignment.RoutineTaskId == uuid.Nil || assignment.RoutineTaskRecordId == uuid.Nil ||
		assignment.RoutineId == uuid.Nil || assignment.ActorUserId == uuid.Nil ||
		assignment.Purpose == "" || len(assignment.Payload) == 0 {
		return nil, invalidPayloadException(fmt.Errorf("routine task assignment is incomplete"))
	}

	var payload any
	switch assignment.Purpose {
	case enums.RoutineTaskPurpose_CreateRootShelf:
		payload = &routinetasktypes.CreateRootShelfRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_UpdateRootShelf:
		payload = &routinetasktypes.UpdateRootShelfRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_ResetRootShelf:
		payload = &routinetasktypes.ResetRootShelfRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_CreateSubShelf:
		payload = &routinetasktypes.CreateSubShelfRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_UpdateSubShelf:
		payload = &routinetasktypes.UpdateSubShelfRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_ResetSubShelf:
		payload = &routinetasktypes.ResetSubShelfRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_CreateBlockPack:
		payload = &routinetasktypes.CreateBlockPackRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_UpdateBlockPack:
		payload = &routinetasktypes.UpdateBlockPackRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_ResetBlockPack:
		payload = &routinetasktypes.ResetBlockPackRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_AppendBlock:
		payload = &routinetasktypes.AppendBlockRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_UpdateBlock:
		payload = &routinetasktypes.UpdateBlockRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_ResetBlock:
		payload = &routinetasktypes.ResetBlockRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_CreateRoutine:
		payload = &routinetasktypes.CreateRoutineRoutineTaskPayload{}
	case enums.RoutineTaskPurpose_UpdateRoutine:
		payload = &routinetasktypes.UpdateRoutineRoutineTaskPayload{}
	default:
		return nil, invalidPayloadException(fmt.Errorf("unsupported routine task purpose: %s", assignment.Purpose))
	}

	if err := json.Unmarshal(assignment.Payload, payload); err != nil {
		return nil, invalidPayloadException(err)
	}
	if validator != nil {
		if err := validator.Struct(payload); err != nil {
			return nil, invalidPayloadException(err)
		}
	}

	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, invalidPayloadException(err)
	}
	var payloadValue any
	if err := json.Unmarshal(rawPayload, &payloadValue); err != nil {
		return nil, invalidPayloadException(err)
	}
	payloadValue = matchPayloadValue(payloadValue, assignment.PatternValues, true)
	preparedPayload, err := json.Marshal(payloadValue)
	if err != nil {
		return nil, invalidPayloadException(err)
	}

	return &routinetasktypes.PreparedRoutineTask{
		RoutineTaskId:       assignment.RoutineTaskId,
		RoutineTaskRecordId: assignment.RoutineTaskRecordId,
		RoutineId:           assignment.RoutineId,
		ActorUserId:         assignment.ActorUserId,
		Attempt:             assignment.Attempt,
		Purpose:             assignment.Purpose,
		Payload:             preparedPayload,
		PreparedAt:          time.Now().UTC(),
	}, nil
}

func invalidPayloadException(err error) *exceptions.Exception {
	return exceptions.New(
		"InvalidRoutineTaskPayload",
		"RoutineTask",
		"Prepare",
		"Routine task payload is invalid",
		http.StatusBadRequest,
	).WithOrigin(err)
}

func matchPayloadValue(value any, values map[string]string, allowStrings bool) any {
	if len(values) == 0 {
		return value
	}

	switch typed := value.(type) {
	case string:
		if !allowStrings {
			return typed
		}
		matched := typed
		for key, resolvedValue := range values {
			matched = strings.ReplaceAll(matched, "{{"+key+"}}", resolvedValue)
		}
		return matched
	case []any:
		matched := make([]any, len(typed))
		for index, item := range typed {
			matched[index] = matchPayloadValue(item, values, allowStrings)
		}
		return matched
	case map[string]any:
		isTemplateBlock := false
		if props, ok := typed["props"].(map[string]any); ok {
			if template, ok := props["template"].(bool); ok {
				isTemplateBlock = template
				if template {
					delete(props, "template")
				}
			}
		}

		matched := make(map[string]any, len(typed))
		for key, item := range typed {
			if key == "pattern" {
				continue
			}
			childAllowsStrings := allowStrings
			if key == "arborizedEditableBlock" || key == "children" {
				childAllowsStrings = isTemplateBlock
			}
			matched[key] = matchPayloadValue(item, values, childAllowsStrings)
		}
		return matched
	default:
		return value
	}
}
