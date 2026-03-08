package constants

const (
	DevelopmentVersion = "v1"
	ProductionVersion  = "v1"
	TestVersion        = "v1"
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
