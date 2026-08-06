package eventbuilders

import (
	"testing"
	"time"

	"github.com/google/uuid"

	coreeventscontract "github.com/HiIamJeff67/notezy-backend/contracts/core/v1/events"
	routinetasktypes "github.com/HiIamJeff67/notezy-backend/contracts/durablejob/v1/types/routine-tasks"
)

func TestRoutineTaskCompletionEventBuilderBuildsCoreOwnedLifecycleEvent(t *testing.T) {
	taskId := uuid.New()
	recordId := uuid.New()
	workerId := uuid.New()
	completedAt := time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC)
	event := NewRoutineTaskCompletionEventBuilder().Build(
		routinetasktypes.CompletedRoutineTask{
			RoutineTaskId:       taskId,
			RoutineTaskRecordId: recordId,
			CompletedAt:         completedAt,
			PreparedTask: &routinetasktypes.PreparedRoutineTask{
				Attempt: 2,
			},
		},
		workerId,
		completedAt,
	)

	if event.EventType != coreeventscontract.EventType_RoutineTaskCompleted ||
		event.AggregateType != coreeventscontract.AggregateType_RoutineTask ||
		event.AggregateId != taskId || event.KafkaKey != taskId.String() {
		t.Fatalf("unexpected event envelope: %#v", event)
	}
	if event.Data.WorkerId != workerId || event.Data.Attempt != 2 {
		t.Fatalf("unexpected event data: %#v", event.Data)
	}
}
