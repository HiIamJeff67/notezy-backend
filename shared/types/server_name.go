package types

type BackendServerName string

const (
	BackendServerName_EastAsia    = "EastAsia"
	BackendServerName_EastAmerica = "EastAmerica"
	BackendServerName_WestAmerica = "WestAmerica"
	BackendServerName_WestEurope  = "Europe"
)

var AllBackendServerNames = []BackendServerName{
	BackendServerName_EastAsia,
	BackendServerName_EastAmerica,
	BackendServerName_WestAmerica,
	BackendServerName_WestEurope,
}

var _backendServerNames = map[string]BackendServerName{
	"EastAsia":    BackendServerName_EastAsia,
	"EastAmerica": BackendServerName_EastAmerica,
	"WestAmerica": BackendServerName_WestAmerica,
	"WestEurope":  BackendServerName_WestEurope,
}

func (bsn BackendServerName) String() string {
	return string(bsn)
}

func IsBackendServerName(backendServerNameString string) bool {
	_, ok := _backendServerNames[backendServerNameString]
	return ok
}

func ConvertToBackendServerName(backendServerNameString string) (BackendServerName, bool) {
	backendServerName, ok := _backendServerNames[backendServerNameString]
	return backendServerName, ok
}
