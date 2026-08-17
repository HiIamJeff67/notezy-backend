package views

import (
	userview "github.com/HiIamJeff67/notegic-backend/internal/core/data/database/schemas/views/user"
)

var MigratingViewSQLs = []string{
	userview.NotificationUserViewSQL,
}
