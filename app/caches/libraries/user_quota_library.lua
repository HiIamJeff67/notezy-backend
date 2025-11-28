#!lua name=user_quota_library

-- Helper function to safely decode JSON
local function safe_decode(str)
    local ok, result = pcall(cjson.decode, str)
    if ok then return result else return nil end
end

-- Helper function to safely encode JSON
local function safe_encode(tbl)
    local ok, result = pcall(cjson.encode, tbl)
    if ok then return result else return nil end
end

-- check_and_update_user_quota_by_formatted_key:
-- Atomically checks the limit and updates the user quota of the user with the given formatted key
-- keys[1]: The quota hash key
-- argv[1]: User Quota Field (e.g., "blockCount")
-- argv[2]: Change Amount (N) - Positive for increment, Negative for decrement
-- argv[3]: Max Limit (Only used when N > 0)
-- argv[4]: Expiration Time (Seconds)
local CHECK_AND_UPDATE_USER_QUOTA_BY_FORMATTED_KEY_NUM_OF_KEYS = 1
local CHECK_AND_UPDATE_USER_QUOTA_BY_FORMATTED_KEY_NUM_OF_ARGV = 4
local function check_and_update_user_quota_by_formatted_key(keys, argv)
    if #keys ~= CHECK_AND_UPDATE_USER_QUOTA_BY_FORMATTED_KEY_NUM_OF_KEYS
        or #argv ~= CHECK_AND_UPDATE_USER_QUOTA_BY_FORMATTED_KEY_NUM_OF_ARGV then
        return { new_value = -1, error = "Invalid arguments" }
    end

    local key = keys[1]
    local field = argv[1]
    local change_amount = tonumber(argv[2]) or 0
    local max_limit = tonumber(argv[3]) or 0
    local ttl = tonumber(argv[4]) or 0

    local cache_string = redis.call('GET', key)
    if not cache_string then
        return { new_value = -1, error = "Cache not found" }
    end

    local cache = safe_decode(cache_string)
    if not cache then
        return { new_value = -1, error = "Failed to decode JSON" }
    end

    local current = tonumber(cache[field]) or 0
    local new_value = current + change_amount

    if change_amount > 0 then
        if new_value > max_limit then
            return {
                new_value = -1,
                error = string.format("Quota exceeded for %s. Current: %d, Request: +%d, Limit: %d", field, current,
                    change_amount, max_limit)
            }
        end
    else
        if new_value < 0 then
            return {
                new_value = -1,
                error = string.format("Quota cannot be negative for %s. Current: %d, Request: %d", field, current,
                    change_amount)
            }
        end
    end

    cache[field] = new_value
    local new_json = safe_encode(cache)
    if new_json then
        redis.call('SET', key, new_json)
        redis.call('EXPIRE', key, ttl)
        return { new_value = new_value, error = nil }
    else
        return { new_value = -1, error = "Failed to encode JSON" }
    end
end

-- best_effor_batch_check_and_update_user_quotas_by_formatted_keys:
-- Atomically batch checks the limit and updates the user quotas among mutiple users with the given formatted keys
-- (the passing formatted keys may be different to each others)
--
-- keys: array of quota hash keys
-- argv: array of json object containing required data for accounting
-- Format of argv: [field_1, change_amount_1, max_limit_1, ttl_1, field_2, change_amount_2, max_limit_2, ttl_2, ...]
-- Note: Since the update operation is based on mutiple different formatted keys,
--       so we should use the strategy of "Best Effort" which means if one update operation failed, we just ignore it
local BATCH_CHECK_AND_UPDATE_USER_QUOTAS_BY_FORMATTED_KEYS_ARGV_PER_KEY = 4
local function best_effort_batch_check_and_update_user_quotas_by_formatted_keys(keys, argv)
    if #keys == 0 or #argv == 0 then return { updated_count = 0, error = nil } end
    if #argv % BATCH_CHECK_AND_UPDATE_USER_QUOTAS_BY_FORMATTED_KEYS_ARGV_PER_KEY ~= 0 then
        return { updated_count = 0, error = "Argv size mismatch" }
    end

    local updated_count = 0
    for i = 1, #keys do
        local base = (i - 1) * BATCH_CHECK_AND_UPDATE_USER_QUOTAS_BY_FORMATTED_KEYS_ARGV_PER_KEY
        local key = keys[i]
        local field = argv[base + 1]
        local change = tonumber(argv[base + 2]) or 0
        local limit = tonumber(argv[base + 3]) or 0
        local ttl = tonumber(argv[base + 4]) or 0

        local str = redis.call('GET', key)
        if str then
            local cache = safe_decode(str)
            if cache then
                local current = tonumber(cache[field]) or 0
                local new_value = current + change
                local valid = false

                if change > 0 then
                    if new_value <= limit then valid = true end
                else
                    if new_value >= 0 then valid = true end
                end

                if valid then
                    cache[field] = new_value
                    local new_json = safe_encode(cache)
                    if new_json then
                        redis.call('SET', key, new_json)
                        redis.call('EXPIRE', key, ttl)
                        updated_count = updated_count + 1
                    end
                end
            end
        end
    end
    return { updated_count = updated_count, error = nil }
end

-- all_or_nothing_batch_check_and_update_user_quotas_by_formatted_keys:
-- Atomically batch checks the limit and updates the user quotas among mutiple users with the given formatted keys
-- (the passing formatted keys may be different to each others)
--
-- keys: array of quota hash keys
-- argv: array of json object containing required data for accounting
-- Format of argv: [field_1, change_amount_1, max_limit_1, ttl_1, field_2, change_amount_2, max_limit_2, ttl_2, ...]
-- Note: This function use "All or Nothing" strategy to update the user quotas, which means the entire update operation will abort if there's anywhere went wrong
local function all_or_nothing_batch_check_and_update_user_quotas_by_formatted_keys(keys, argv)
    if #keys == 0 or #argv == 0 then return { updated_count = 0, error = nil } end

    local pending_updates = {}

    for i = 1, #keys do
        local base = (i - 1) * BATCH_CHECK_AND_UPDATE_USER_QUOTAS_BY_FORMATTED_KEYS_ARGV_PER_KEY
        local key = keys[i]
        local field = argv[base + 1]
        local change = tonumber(argv[base + 2]) or 0
        local limit = tonumber(argv[base + 3]) or 0
        local ttl = tonumber(argv[base + 4]) or 0

        local str = redis.call('GET', key)
        if not str then return { updated_count = 0, error = "Cache not found for key: " .. key } end

        local cache = safe_decode(str)
        if not cache then return { updated_count = 0, error = "JSON decode failed for key: " .. key } end

        local current = tonumber(cache[field]) or 0
        local new_value = current + change

        if change > 0 then
            if new_value > limit then
                return { updated_count = 0, error = string.format("Quota exceeded for %s in key %s", field, key) }
            end
        else
            if new_value < 0 then
                return { updated_count = 0, error = string.format("Negative quota for %s in key %s", field, key) }
            end
        end

        cache[field] = new_value
        table.insert(pending_updates, { key = key, data = cache, ttl = ttl })
    end

    for _, op in ipairs(pending_updates) do
        local new_json = safe_encode(op.data)
        if new_json then
            redis.call('SET', op.key, new_json)
            redis.call('EXPIRE', op.key, op.ttl)
        end
    end

    return { updated_count = #keys, error = nil }
end

-- best_effort_batch_check_and_update_user_quotas_by_formatted_key
-- Atomically batch checks the limit and updates the user quotas of a single user with the given formatted key
--
-- keys[1]: The quota hash key
-- argv: array of json object containing required data for accounting
-- Format of argv: [field_1, change_amount_1, max_limit_1, ttl_1, field_2, change_amount_2, max_limit_2, ttl_2, ...]
-- Note : This function use the "Best Effort" strategy
local BEST_EFFORT_BATCH_CHECK_AND_UPDATE_USER_QUOTAS_BY_FORMATTED_KEY_NUM_OF_KEYS = 1
local BEST_EFFORT_BATCH_CHECK_AND_UPDATE_USER_QUOTAS_BY_FORMATTED_KEY_BASE_NUM_OF_ARGV = 4
local function best_effort_batch_check_and_update_user_quotas_by_formatted_key(keys, argv)
    local key = keys[1]
    local str = redis.call('GET', key)
    if not str then return { updated_count = 0, error = "Cache not found" } end

    local cache = safe_decode(str)
    if not cache then return { updated_count = 0, error = "JSON decode failed" } end

    local num_ops = #argv / BEST_EFFORT_BATCH_CHECK_AND_UPDATE_USER_QUOTAS_BY_FORMATTED_KEY_BASE_NUM_OF_ARGV
    local updated_count = 0
    local max_ttl = 0
    local is_modified = false

    for i = 0, num_ops - 1 do
        local base = i * BEST_EFFORT_BATCH_CHECK_AND_UPDATE_USER_QUOTAS_BY_FORMATTED_KEY_BASE_NUM_OF_ARGV
        local field = argv[base + 1]
        local change = tonumber(argv[base + 2]) or 0
        local limit = tonumber(argv[base + 3]) or 0
        local ttl = tonumber(argv[base + 4]) or 0

        local current = tonumber(cache[field]) or 0
        local new_value = current + change
        local valid = false

        if change > 0 then
            if new_value <= limit then valid = true end
        else
            if new_value >= 0 then valid = true end
        end

        if valid then
            cache[field] = new_value
            updated_count = updated_count + 1
            is_modified = true
            if ttl > max_ttl then max_ttl = ttl end
        end
    end

    if is_modified then
        local new_json = safe_encode(cache)
        if new_json then
            redis.call('SET', key, new_json)
            redis.call('EXPIRE', key, max_ttl)
        end
    end

    return { updated_count = updated_count, error = nil }
end

-- all_or_nothing_batch_check_and_update_user_quotas_by_formatted_key
-- Atomically batch checks the limit and updates the user quotas of a single user with the given formatted key
--
-- keys[1]: The quota hash key
-- argv: array of json object containing required data for accounting
-- Format of argv: [field_1, change_amount_1, max_limit_1, ttl_1, field_2, change_amount_2, max_limit_2, ttl_2, ...]
-- Note : This function use the "All or Nothing" strategy
local function all_or_nothing_batch_check_and_update_user_quotas_by_formatted_key(keys, argv)
    local key = keys[1]
    local str = redis.call('GET', key)
    if not str then return { updated_count = 0, error = "Cache not found" } end

    local cache = safe_decode(str)
    if not cache then return { updated_count = 0, error = "JSON decode failed" } end

    local num_ops = #argv / BEST_EFFORT_BATCH_CHECK_AND_UPDATE_USER_QUOTAS_BY_FORMATTED_KEY_BASE_NUM_OF_ARGV
    local max_ttl = 0

    local temp_cache = {}

    for i = 0, num_ops - 1 do
        local base = i * BEST_EFFORT_BATCH_CHECK_AND_UPDATE_USER_QUOTAS_BY_FORMATTED_KEY_BASE_NUM_OF_ARGV
        local field = argv[base + 1]
        local change = tonumber(argv[base + 2]) or 0
        local limit = tonumber(argv[base + 3]) or 0
        local ttl = tonumber(argv[base + 4]) or 0

        local current = temp_cache[field] or (tonumber(cache[field]) or 0)
        local new_value = current + change

        if change > 0 then
            if new_value > limit then
                return { updated_count = 0, error = string.format("Quota exceeded for %s", field) }
            end
        else
            if new_value < 0 then
                return { updated_count = 0, error = string.format("Negative quota for %s", field) }
            end
        end

        temp_cache[field] = new_value
        if ttl > max_ttl then max_ttl = ttl end
    end

    for field, val in pairs(temp_cache) do
        cache[field] = val
    end

    local new_json = safe_encode(cache)
    if new_json then
        redis.call('SET', key, new_json)
        redis.call('EXPIRE', key, max_ttl)
    end

    return { updated_count = num_ops, error = nil }
end

redis.register_function('check_and_update_user_quota_by_formatted_key', check_and_update_user_quota_by_formatted_key)
redis.register_function('best_effort_batch_check_and_update_user_quotas_by_formatted_keys',
    best_effort_batch_check_and_update_user_quotas_by_formatted_keys)
redis.register_function('all_or_nothing_batch_check_and_update_user_quotas_by_formatted_keys',
    all_or_nothing_batch_check_and_update_user_quotas_by_formatted_keys)
redis.register_function('best_effort_batch_check_and_update_user_quotas_by_formatted_key',
    best_effort_batch_check_and_update_user_quotas_by_formatted_key)
redis.register_function('all_or_nothing_batch_check_and_update_user_quotas_by_formatted_key',
    all_or_nothing_batch_check_and_update_user_quotas_by_formatted_key)
