package handlers

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"

	routinetasktypes "github.com/HiIamJeff67/notezy-backend/contracts/durable-job/v1/types/routine-tasks"
	enums "github.com/HiIamJeff67/notezy-backend/contracts/types/enums"

	validation "github.com/HiIamJeff67/notezy-backend/internal/durablejob/validations"
)

func TestPurposeHandlerPreparesAssignmentWithoutDatabaseAccess(t *testing.T) {
	payload, err := json.Marshal(routinetasktypes.CreateRootShelfRoutineTaskPayload{
		Name: "Daily {{date}}",
	})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	assignment := routinetasktypes.RoutineTaskAssignment{
		RoutineTaskId:       uuid.New(),
		RoutineTaskRecordId: uuid.New(),
		RoutineId:           uuid.New(),
		ActorUserId:         uuid.New(),
		Purpose:             enums.RoutineTaskPurpose_CreateRootShelf,
		Payload:             payload,
		Attempt:             1,
		ScheduledAt:         time.Now().UTC(),
		StartedAt:           time.Now().UTC(),
		PatternValues:       map[string]string{"date": "2026-08-05"},
	}

	prepared, exception := NewPurposeHandler(validation.New()).HandlerFunc(t.Context(), assignment)
	if exception != nil {
		t.Fatalf("prepare assignment: %v", exception)
	}
	if prepared == nil || prepared.RoutineTaskId != assignment.RoutineTaskId {
		t.Fatalf("prepared task = %#v", prepared)
	}

	var preparedPayload routinetasktypes.CreateRootShelfRoutineTaskPayload
	if err := json.Unmarshal(prepared.Payload, &preparedPayload); err != nil {
		t.Fatalf("decode prepared payload: %v", err)
	}
	if preparedPayload.Name != "Daily 2026-08-05" {
		t.Fatalf("prepared name = %q", preparedPayload.Name)
	}
}
