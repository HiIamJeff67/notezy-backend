package notificationeventscontract

import eventcontract "github.com/HiIamJeff67/notegic-backend/contracts/types/events"

const NotificationTopic eventcontract.Topic = "notegic.notification.v1"

const (
	AggregateType_Notification    eventcontract.AggregateType = "Notification"
	EventType_NotificationCreated eventcontract.EventType     = "NotificationCreated"
	EventType_NotificationUpdated eventcontract.EventType     = "NotificationUpdated"
)
