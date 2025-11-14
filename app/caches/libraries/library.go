package redislibraries

import (
	_ "embed"
)

const (
	RateLimitRecordLibrary string = "rate_limit_record_library"
	// BatchSynchronizeRateLimitRecordByFormattedKeys
	// is used to batch synchronize the rate limit records using the given formatted keys
	//
	// Arguments format :
	//   - keys: array of formatted keys
	//   - argv: array of json objects containing synchronizeDto
	// Format of argv: [numOfChangingTokens1, isAccumulated1, numOfChangingTokens2, isAccumulated2, ...]
	BatchSynchronizeRateLimitRecordByFormattedKeysFunction string = "batch_synchronize_rate_limit_record_by_formatted_keys"

	// batch_delete_rate_limit_record_by_formatted_keys:
	// Redis functions to batch delete the rate limit record by the given formatted keys
	// Arguments format :
	//   - keys: array of formatted keys
	//   - _ : placeholder for argv, but we don't use it here
	BatchDeleteRateLimitRecordByFormattedKeysFunction string = "batch_delete_rate_limit_record_by_formatted_keys"
)

var (
	//go:embed rate_limit_record_library.lua
	RateLimitRecordLibraryContent string
)
