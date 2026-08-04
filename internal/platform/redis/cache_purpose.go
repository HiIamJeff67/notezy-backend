package redis

type CachePurpose string

const (
	CachePurpose_UserData  CachePurpose = "UserData"
	CachePurpose_RateLimit CachePurpose = "RateLimit"
	CachePurpose_Realtime  CachePurpose = "Realtime"
)

func (purpose CachePurpose) String() string {
	return string(purpose)
}
