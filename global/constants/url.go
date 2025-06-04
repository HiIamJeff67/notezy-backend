package constants

const (
	APIBaseURL         = "/api"                                // the basic api route url
	DevelopmentBaseURL = APIBaseURL + "/" + DevelopmentVersion // the current development version of api
	ProductionBaseURL  = APIBaseURL + "/" + ProductionVersion
)

const (
	BaseURL = DevelopmentBaseURL // use in the entire project apis
)
