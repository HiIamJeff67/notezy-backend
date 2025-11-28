package types

type UserQuotaField string

const (
	UserQuotaField_RootShelfCount      = "rootShelfCount"
	UserQuotaField_BlockPackCount      = "blockPackCount"
	UserQuotaField_BlockCount          = "blockCount"
	UserQuotaField_MaterialCount       = "materialCount"
	UserQuotaField_WorkflowCount       = "workflowCount"
	UserQuotaField_AdditionalItemCount = "additionalItemCount"
)

func (uqf UserQuotaField) String() string {
	return string(uqf)
}
