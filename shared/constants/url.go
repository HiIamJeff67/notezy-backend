package constants

const (
	DevelopmentNamespace = "development"
	ProductionNamespace  = ""
	TestNamespace        = "test"
)

const (
	APIBaseURL         = "/api"                                                             // the basic api route url
	DevelopmentBaseURL = APIBaseURL + "/" + DevelopmentNamespace + "/" + DevelopmentVersion // the current development version of api
	ProductionBaseURL  = APIBaseURL + "/" + ProductionNamespace + "/" + ProductionVersion
	TestBaseURL        = APIBaseURL + "/" + TestNamespace + "/" + TestVersion
)

const (
	BaseURL = DevelopmentBaseURL // use in the entire project apis
)
