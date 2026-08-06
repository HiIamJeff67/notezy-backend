package ratelimitrecord

import types "github.com/HiIamJeff67/notezy-backend/shared/types"

type Config struct {
	ServerRange types.Range[int, int]
}
