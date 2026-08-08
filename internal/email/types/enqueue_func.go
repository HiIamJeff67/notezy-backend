package types

import exceptions "github.com/HiIamJeff67/notezy-backend/contracts/types/exceptions"

type EnqueueFunc func(
	emailObject EmailObject,
	taskType EmailTaskType,
	maxRetries int,
	priority int,
) *exceptions.Exception
