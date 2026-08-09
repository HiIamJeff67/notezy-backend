package apiexceptions

type RoutineTagException struct {
	CoreException
}

func NewRoutineTagException() RoutineTagException {
	return RoutineTagException{
		CoreException: NewCoreException("RoutineTag"),
	}
}
