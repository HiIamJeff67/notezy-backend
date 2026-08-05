package gatewaycontract

const (
	DevelopmentNamespace = "development"
	ProductionNamespace  = ""
	TestNamespace        = "test"

	APIGroupBase          = "api"
	APIDevelopmentBaseURL = APIGroupBase + "/" + DevelopmentNamespace + "/" + DevelopmentVersion
	APIProductionBaseURL  = APIGroupBase + "/" + ProductionNamespace + "/" + ProductionVersion
	APITestBaseURL        = APIGroupBase + "/" + TestNamespace + "/" + TestVersion
)
