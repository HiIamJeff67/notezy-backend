package apiexceptions

type RoutineTaskException struct {
	CoreException
}

func NewRoutineTaskException() RoutineTaskException {
	return RoutineTaskException{
		CoreException: NewCoreException("RoutineTask"),
	}
}
