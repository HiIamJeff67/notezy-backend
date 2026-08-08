package apicontract

const (
	GetMyRoutineTaskByIdOperation                       = "routine-task.get-by-id"
	GetAllMyRoutineTasksByRoutineIdsOperation           = "routine-task.get-all-by-routine-ids"
	GetAllMyRoutineTasksOperation                       = "routine-task.get-all"
	CreateRoutineTaskByRoutineIdOperation               = "routine-task.create-by-routine-id"
	UpdateMyRoutineTaskByIdOperation                    = "routine-task.update"
	PauseMyRoutineTaskByIdOperation                     = "routine-task.pause"
	ResumeMyRoutineTaskByIdOperation                    = "routine-task.resume"
	HardDeleteMyRoutineTaskByIdOperation                = "routine-task.hard-delete"
	HardDeleteMyRoutineTasksByIdsOperation              = "routine-task.hard-delete-many"
	VisualizeMyRoutineTaskStatusCountOperation          = "routine-task.visualize-status-count"
	VisualizeMyRoutineTaskPurposeCountOperation         = "routine-task.visualize-purpose-count"
	VisualizeMyRoutineTaskScheduledAtCountOperation     = "routine-task.visualize-scheduled-at-count"
	VisualizeMyRoutineTaskActualStartedAtCountOperation = "routine-task.visualize-actual-started-at-count"
	VisualizeMyRoutineTaskActualEndedAtCountOperation   = "routine-task.visualize-actual-ended-at-count"
	SearchRoutineTasksOperation                         = "graphql.search-routine-tasks"
)
