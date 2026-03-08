package constants

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
