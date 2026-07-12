package types

type RealtimeTicketAudience string

const (
	RealtimeTicketAudience_Connection RealtimeTicketAudience = "notezy-realtime-connection"
	RealtimeTicketAudience_BlockPack  RealtimeTicketAudience = "notezy-realtime-block-pack"
)

func (a RealtimeTicketAudience) String() string {
	return string(a)
}
