package caches

import (
	types "notezy-backend/global/types"
)

var (
	RecentPagesRange = types.Range{Start: 8, Size: 8} // server number: 8 - 15 (included)
)
