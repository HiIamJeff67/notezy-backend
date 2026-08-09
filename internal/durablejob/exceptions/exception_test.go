package exceptions

import "testing"

func TestNewDurableJobException(t *testing.T) {
	exception := NewDurableJobException("RoutineTask")
	if exception.Domain != "RoutineTask" {
		t.Fatalf("domain = %q, want RoutineTask", exception.Domain)
	}
}
