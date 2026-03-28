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
	APIGroupBase       = "api"                                                                // the basic api route name space
	DevelopmentBaseURL = APIGroupBase + "/" + DevelopmentNamespace + "/" + DevelopmentVersion // the current development version of api
	ProductionBaseURL  = APIGroupBase + "/" + ProductionNamespace + "/" + ProductionVersion
	TestBaseURL        = APIGroupBase + "/" + TestNamespace + "/" + TestVersion
)

const (
	CurrentBaseURL = DevelopmentBaseURL // use in the entire project apis
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
