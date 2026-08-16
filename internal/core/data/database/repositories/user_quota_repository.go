package repositories

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	exceptions "github.com/HiIamJeff67/notegic-backend/contracts/types/exceptions"

	inputs "github.com/HiIamJeff67/notegic-backend/internal/core/data/database/inputs"
	options "github.com/HiIamJeff67/notegic-backend/internal/core/data/database/options"
	schemas "github.com/HiIamJeff67/notegic-backend/internal/core/data/database/schemas"
)

type UserQuotaRepositoryInterface interface {
	GetRoutineTaskCostUnitUsed(ctx context.Context, userId uuid.UUID, opts ...options.RepositoryOptions) (int64, *exceptions.Exception)
	InitializeMissing(ctx context.Context, now time.Time, opts ...options.RepositoryOptions) *exceptions.Exception
	InitializeMissingForUserIds(ctx context.Context, userIds []uuid.UUID, now time.Time, opts ...options.RepositoryOptions) *exceptions.Exception
	ResetDue(ctx context.Context, now time.Time, opts ...options.RepositoryOptions) (int64, *exceptions.Exception)
	ConsumeRoutineTaskCostUnits(ctx context.Context, consumptionInputs []inputs.ConsumeRoutineTaskCostUnitInput, opts ...options.RepositoryOptions) ([]uuid.UUID, *exceptions.Exception)
}

type UserQuotaRepository struct{}

func NewUserQuotaRepository() UserQuotaRepositoryInterface {
	return &UserQuotaRepository{}
}

func (r *UserQuotaRepository) GetRoutineTaskCostUnitUsed(
	ctx context.Context,
	userId uuid.UUID,
	opts ...options.RepositoryOptions,
) (int64, *exceptions.Exception) {
	if userId == uuid.Nil {
		return 0, exceptions.New(
			"InvalidDto",
			"UserQuota",
			"GetRoutineTaskCostUnitUsed",
			"User quota request is invalid",
			http.StatusBadRequest,
		)
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)

	var routineTaskCostUnitUsed int64
	result := parsedOptions.DB.
		WithContext(ctx).
		Model(&schemas.UserQuota{}).
		Raw(`
		SELECT COALESCE(routine_task_cost_unit_used, 0)
		FROM "UserQuotaTable"
		WHERE user_id = ?
		`, userId).
		Scan(&routineTaskCostUnitUsed)
	if result.Error != nil {
		return 0, exceptions.New(
			"FailedToGet",
			"UserQuota",
			"GetRoutineTaskCostUnitUsed",
			"Failed to retrieve routine task quota usage",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	return routineTaskCostUnitUsed, nil
}

func (r *UserQuotaRepository) InitializeMissing(
	ctx context.Context,
	now time.Time,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	return r.initializeMissing(ctx, nil, now, opts...)
}

func (r *UserQuotaRepository) InitializeMissingForUserIds(
	ctx context.Context,
	userIds []uuid.UUID,
	now time.Time,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	if len(userIds) == 0 {
		return nil
	}

	return r.initializeMissing(ctx, userIds, now, opts...)
}

func (r *UserQuotaRepository) initializeMissing(
	ctx context.Context,
	userIds []uuid.UUID,
	now time.Time,
	opts ...options.RepositoryOptions,
) *exceptions.Exception {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	if !parsedOptions.IsTransactionStarted {
		return exceptions.New(
			"TransactionRequired",
			"UserQuota",
			"InitializeMissing",
			"User quota initialization requires a transaction",
			http.StatusInternalServerError,
		)
	}

	userIdFilter := ""
	queryArgs := []any{now, now}
	if len(userIds) > 0 {
		userIdFilter = "user_table.id IN ? AND "
		queryArgs = append(queryArgs, userIds)
	}

	query := fmt.Sprintf(`
		INSERT INTO "UserQuotaTable" (
			id,
			user_id,
			routine_task_cost_unit_used,
			cycle_started_at,
			next_reset_at,
			updated_at,
			created_at
		)
		SELECT
			gen_random_uuid(),
			user_table.id,
			0,
			COALESCE(
				billing.next_billing_date - INTERVAL '30 days',
				billing.start_date,
				user_table.created_at
			),
			COALESCE(
				billing.next_billing_date,
				billing.start_date + INTERVAL '30 days',
				user_table.created_at + INTERVAL '30 days'
			),
			?,
			?
		FROM "UserTable" AS user_table
		LEFT JOIN LATERAL (
			SELECT start_date, next_billing_date
			FROM "UsersToBillingPlansTable"
			WHERE user_id = user_table.id
				AND status = 'ACTIVE'
			ORDER BY start_date DESC
			LIMIT 1
		) AS billing ON TRUE
		WHERE %sNOT EXISTS (
			SELECT 1
			FROM "UserQuotaTable"
			WHERE user_id = user_table.id
		)
		ON CONFLICT (user_id) DO NOTHING
	`, userIdFilter)

	result := parsedOptions.DB.
		WithContext(ctx).
		Model(&schemas.UserQuota{}).
		Exec(query, queryArgs...)
	if result.Error != nil {
		return exceptions.New(
			"FailedToInitialize",
			"UserQuota",
			"InitializeMissing",
			"Failed to initialize missing user quotas",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	return nil
}

func (r *UserQuotaRepository) ResetDue(
	ctx context.Context,
	now time.Time,
	opts ...options.RepositoryOptions,
) (int64, *exceptions.Exception) {
	parsedOptions := options.ParseRepositoryOptions(opts...)
	if !parsedOptions.IsTransactionStarted {
		return 0, exceptions.New(
			"TransactionRequired",
			"UserQuota",
			"ResetDue",
			"User quota reset requires a transaction",
			http.StatusInternalServerError,
		)
	}

	result := parsedOptions.DB.
		WithContext(ctx).
		Model(&schemas.UserQuota{}).
		Exec(`
		UPDATE "UserQuotaTable"
		SET
			routine_task_cost_unit_used = 0,
			cycle_started_at = ?,
			next_reset_at = ? + INTERVAL '30 days',
			updated_at = ?
		WHERE next_reset_at <= ?
		`, now, now, now, now)
	if result.Error != nil {
		return 0, exceptions.New(
			"FailedToReset",
			"UserQuota",
			"ResetDue",
			"Failed to reset due user quotas",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	return result.RowsAffected, nil
}

func (r *UserQuotaRepository) ConsumeRoutineTaskCostUnits(
	ctx context.Context,
	consumptionInputs []inputs.ConsumeRoutineTaskCostUnitInput,
	opts ...options.RepositoryOptions,
) ([]uuid.UUID, *exceptions.Exception) {
	if len(consumptionInputs) == 0 {
		return []uuid.UUID{}, nil
	}

	parsedOptions := options.ParseRepositoryOptions(opts...)
	if !parsedOptions.IsTransactionStarted {
		return nil, exceptions.New(
			"TransactionRequired",
			"UserQuota",
			"ConsumeRoutineTaskCostUnits",
			"Routine task quota consumption requires a transaction",
			http.StatusInternalServerError,
		)
	}

	valuePlaceholders := make([]string, 0, len(consumptionInputs))
	valueArgs := make([]any, 0, len(consumptionInputs)*5)
	for _, consumptionInput := range consumptionInputs {
		if consumptionInput.RoutineTaskId == uuid.Nil || consumptionInput.UserId == uuid.Nil ||
			consumptionInput.CostUnit < 0 || consumptionInput.ScheduledAt.IsZero() {
			return nil, exceptions.New(
				"InvalidDto",
				"UserQuota",
				"ConsumeRoutineTaskCostUnits",
				"Routine task quota consumption input is invalid",
				http.StatusBadRequest,
			)
		}

		valuePlaceholders = append(valuePlaceholders, "(?::uuid, ?::uuid, ?::bigint, ?::integer, ?::timestamptz)")
		valueArgs = append(
			valueArgs,
			consumptionInput.RoutineTaskId,
			consumptionInput.UserId,
			consumptionInput.CostUnit,
			consumptionInput.Priority,
			consumptionInput.ScheduledAt,
		)
	}

	var routineTaskIds []uuid.UUID
	result := parsedOptions.DB.
		WithContext(ctx).
		Model(&schemas.UserQuota{}).
		Raw(fmt.Sprintf(`
		WITH consumption(routine_task_id, user_id, cost_unit, priority, scheduled_at) AS (
			VALUES %s
		),
		eligible_consumption AS (
			SELECT
				ranked_consumption.routine_task_id,
				ranked_consumption.user_id,
				ranked_consumption.cost_unit
			FROM (
				SELECT
					consumption.*,
					sum(cost_unit) OVER (
						PARTITION BY user_id
						ORDER BY priority DESC, scheduled_at ASC, routine_task_id ASC
					) AS accumulated_cost_unit
				FROM consumption
			) AS ranked_consumption
			JOIN "UserQuotaTable" AS user_quota ON user_quota.user_id = ranked_consumption.user_id
			JOIN "UserTable" AS user_table ON user_table.id = ranked_consumption.user_id
			JOIN "PlanLimitationTable" AS plan_limitation ON plan_limitation.key = user_table.plan
			WHERE user_quota.routine_task_cost_unit_used + accumulated_cost_unit <= plan_limitation.max_routine_task_cost_unit_count
		),
		consumption_totals AS (
			SELECT user_id, sum(cost_unit) AS cost_unit
			FROM eligible_consumption
			GROUP BY user_id
		),
		consumed AS (
			UPDATE "UserQuotaTable" AS user_quota
			SET
				routine_task_cost_unit_used = user_quota.routine_task_cost_unit_used + consumption_totals.cost_unit,
				updated_at = NOW()
			FROM consumption_totals
			JOIN "UserTable" AS user_table ON user_table.id = consumption_totals.user_id
			JOIN "PlanLimitationTable" AS plan_limitation ON plan_limitation.key = user_table.plan
			WHERE user_quota.user_id = consumption_totals.user_id
				AND user_quota.routine_task_cost_unit_used + consumption_totals.cost_unit <= plan_limitation.max_routine_task_cost_unit_count
			RETURNING user_quota.user_id
		)
		SELECT eligible_consumption.routine_task_id
		FROM eligible_consumption
		JOIN consumed ON consumed.user_id = eligible_consumption.user_id
		`, strings.Join(valuePlaceholders, ",")), valueArgs...).
		Scan(&routineTaskIds)
	if result.Error != nil {
		return nil, exceptions.New(
			"FailedToConsume",
			"UserQuota",
			"ConsumeRoutineTaskCostUnits",
			"Failed to consume routine task quota",
			http.StatusInternalServerError,
			true,
		).WithOrigin(result.Error)
	}

	return routineTaskIds, nil
}
