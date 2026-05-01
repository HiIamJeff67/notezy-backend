package types

type ModeType string

const (
	ModeType_Development ModeType = "development"
	ModeType_Production  ModeType = "production"
	ModeType_Test        ModeType = "test"
)
