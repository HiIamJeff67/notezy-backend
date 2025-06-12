package util

import types "notezy-backend/app/shared/types"

func GetMinInMap[K comparable, T types.Number](searchMap map[K]T) (res T) {
	for _, value := range searchMap {
		res = min(res, value)
	}
	return res
}

func GetMaxInMap[K comparable, T types.Number](searchMap map[K]T) (res T) {
	for _, value := range searchMap {
		res = max(res, value)
	}
	return res
}
