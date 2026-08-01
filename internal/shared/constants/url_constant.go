package constants

const (
	Protocol = "http"
	Host     = "localhost"
	Port     = "7777" // should be the same as the environment variables of DOCKER_GIN_PORT and GIN_PORT
)

const (
	DevelopmentNamespace = "development"
	ProductionNamespace  = ""
	TestNamespace        = "test"
)

const (
	APIGroupBase          = "api"                                                                // the basic api route name space
	APIDevelopmentBaseURL = APIGroupBase + "/" + DevelopmentNamespace + "/" + DevelopmentVersion // the current development version of api
	APIProductionBaseURL  = APIGroupBase + "/" + ProductionNamespace + "/" + ProductionVersion
	APITestBaseURL        = APIGroupBase + "/" + TestNamespace + "/" + TestVersion
)

const (
	RealtimeGroupBase          = "realtime"
	RealtimeDevelopmentBaseURL = RealtimeGroupBase + "/" + DevelopmentNamespace + "/" + DevelopmentVersion
)

const (
	CurrentAPIBaseURL = APIDevelopmentBaseURL // use in the entire project apis
)

var URLWhiteList = []string{
	"http",
	"https",
	"mailto",
	"tel",
	"ws",
}

var URLBlackList = []string{
	"javascript",
	"vbscript",
	"file",
	"data",
}
