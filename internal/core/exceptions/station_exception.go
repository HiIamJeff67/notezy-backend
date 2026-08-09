package apiexceptions

type StationException struct {
	CoreException
}

func NewStationException() StationException {
	return StationException{
		CoreException: NewCoreException("Station"),
	}
}
