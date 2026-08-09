package schemas

var MigratingTables = []any{
	&Notification{},
	&InboxEvent{},
	&OutboxEvent{},
	&UserDeletion{},
}
