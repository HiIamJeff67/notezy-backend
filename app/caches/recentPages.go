package caches

import "go-gorm-api/global"

var (
	RecentPagesRange = global.Range{ Start: 10, Size: 10 }
)