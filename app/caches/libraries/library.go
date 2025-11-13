package redisfunctionlibraries

import (
	_ "embed"
)

const (
	RateLimitRecordRedisFunctionsLibrary string = "rate_limit_record_functions_library"
	// BatchSynchronizeRateLimitRecordByFormattedKeys
	// is used to batch synchronize the rate limit records using the given formatted keys
	//
	// Arguments format :
	//   - keys: array of formatted keys
	//   - argv: array of json objects containing synchronizeDto
	// Format of argv: [numOfChangingTokens1, isAccumulated1, numOfChangingTokens2, isAccumulated2, ...]
	BatchSynchronizeRateLimitRecordByFormattedKeysFunction string = "batch_synchronize_rate_limit_record_by_formatted_keys"
	BatchDeleteRateLimitRecordByFormattedKeysFunction      string = "batch_delete_rate_limit_record_by_formatted_keys"
)

var (
	//go:embed rate_limit_record_functions_library.lua
	RateLimitRecordRedisFunctionsLibraryContent string
)
