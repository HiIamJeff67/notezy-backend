package util

import "go-gorm-api/global"

func GetMinInMap[K comparable, T global.Number](searchMap map[K]T) (res T) {
	for _, value := range searchMap { res = min(res, value); }
	return res
}

func GetMaxInMap[K comparable, T global.Number](searchMap map[K]T) (res T) {
	for _, value := range searchMap { res = max(res, value); }
	return res
}