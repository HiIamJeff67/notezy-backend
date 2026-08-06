package types

import exceptions "github.com/HiIamJeff67/notezy-backend/shared/exceptions"

type EnqueueFunc func(
	emailObject EmailObject,
	taskType EmailTaskType,
	maxRetries int,
	priority int,
) *exceptions.Exception
