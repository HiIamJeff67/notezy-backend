package redislibraries

import (
	_ "embed"
)

const (
	RateLimitRecordLibrary string = "rate_limit_record_library"

	// BatchSynchronizeRateLimitRecordByFormattedKeysFunction:
	// Redis functions to batch synchronize the rate limit record by the given formatted keys
	//
	// keys: array of formatted keys
	// argv: array of json objects containing synchronizeDto
	// Format of argv: [num_of_changing_tokens_1, is_accumulated_1, num_of_changing_tokens_2, is_accumulated_2, ...]
	//                 each `NUM_OF_ARGV_PER_KEY` mapping to a key in `keys`
	BatchSynchronizeRateLimitRecordByFormattedKeysFunction string = "batch_synchronize_rate_limit_record_by_formatted_keys"

	// BatchDeleteRateLimitRecordByFormattedKeysFunction:
	// Redis functions to batch delete the rate limit record by the given formatted keys
	//
	// keys: array of formatted keys
	// argv: a placeholder for argv, but we don't use it here
	BatchDeleteRateLimitRecordByFormattedKeysFunction string = "batch_delete_rate_limit_record_by_formatted_keys"
)

const (
	UserQuotaLibrary string = "user_quota_library"
	// CheckAndUpdateUserQuotaByFormattedKeyFunction:
	// Atomically checks the limit and updates the user quota of the user with the given formatted key
	// keys[1]: The quota hash key
	// argv[1]: User Quota Field (e.g., "blockCount")
	// argv[2]: Change Amount (N) - Positive for increment, Negative for decrement
	// argv[3]: Max Limit (Only used when N > 0)
	// argv[4]: Expiration Time (Seconds)
	CheckAndUpdateUserQuotaByFormattedKeyFunction string = "check_and_update_user_quota_by_formatted_key"

	// BestEffortBatchCheckAndUpdateUserQuotasByFormattedKeysFunction:
	// Atomically batch checks the limit and updates the user quotas among mutiple users with the given formatted keys
	// (the passing formatted keys may be different to each others)
	//
	// keys: array of quota hash keys
	// argv: array of json object containing required data for accounting
	// Format of argv: [field_1, change_amount_1, max_limit_1, ttl_1, field_2, change_amount_2, max_limit_2, ttl_2, ...]
	// Note: Since the update operation is based on mutiple different formatted keys,
	//       so we should use the strategy of "Best Effort" which means if one update operation failed, we just ignore it
	BestEffortBatchCheckAndUpdateUserQuotasByFormattedKeysFunction string = "best_effort_batch_check_and_update_user_quotas_by_formatted_keys"

	// AllOrNothingBatchCheckAndUpdateUserQuotasByFormattedKeysFunction:
	// Atomically batch checks the limit and updates the user quotas among mutiple users with the given formatted keys
	// (the passing formatted keys may be different to each others)
	//
	// keys: array of quota hash keys
	// argv: array of json object containing required data for accounting
	// Format of argv: [field_1, change_amount_1, max_limit_1, ttl_1, field_2, change_amount_2, max_limit_2, ttl_2, ...]
	// Note: This function use "All or Nothing" strategy to update the user quotas, which means the entire update operation will abort if there's anywhere went wrong
	AllOrNothingBatchCheckAndUpdateUserQuotasByFormattedKeysFunction string = "all_or_nothing_batch_check_and_update_user_quotas_by_formatted_keys"

	// BestEffortBatchCheckAndUdateUserQuotasByFormattedKeyFunction:
	// Atomically batch checks the limit and updates the user quotas of a single user with the given formatted key
	//
	// keys[1]: The quota hash key
	// argv: array of json object containing required data for accounting
	// Format of argv: [field_1, change_amount_1, max_limit_1, ttl_1, field_2, change_amount_2, max_limit_2, ttl_2, ...]
	// Note : This function use the "Best Effort" strategy
	BestEffortBatchCheckAndUpdateUserQuotasByFormattedKeyFunction string = "best_effort_batch_check_and_update_user_quotas_by_formatted_key"

	// AllOrNothingBatchCheckAndUpdateUserQuotasByFormattedKeyFunction:
	// Atomically batch checks the limit and updates the user quotas of a single user with the given formatted key
	//
	// keys[1]: The quota hash key
	// argv: array of json object containing required data for accounting
	// Format of argv: [field_1, change_amount_1, max_limit_1, ttl_1, field_2, change_amount_2, max_limit_2, ttl_2, ...]
	// Note : This function use the "All or Nothing" strategy
	AllOrNothingBatchCheckAndUpdateUserQuotasByFormattedKeyFunction string = "all_or_nothing_batch_check_and_update_user_quotas_by_formatted_key"
)

var (
	//go:embed rate_limit_record_library.lua
	RateLimitRecordLibraryContent string

	//go:embed user_quota_library.lua
	UserQuotaLibraryContent string
)
