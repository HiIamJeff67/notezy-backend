package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	logs "github.com/HiIamJeff67/notezy-backend/shared/platform/observability/logs"
)

func ViewAllDatabaseEnums(db *gorm.DB) error {
	if logs.NotezyLogger == nil {
		return errors.New("observability logger is not initialized")
	}

	var enumInfos []struct {
		Name   string `gorm:"column:enum_name;"`
		Values string `gorm:"column:enum_values;"`
	}
	result := db.Raw(`
		SELECT
			n.nspname || '.' || t.typname AS enum_name,
			string_agg(e.enumlabel, ', ' ORDER BY e.enumsortorder) AS enum_values
		FROM pg_type t
		JOIN pg_enum e ON t.oid = e.enumtypid
		JOIN pg_namespace n ON n.oid = t.typnamespace
		GROUP BY n.nspname, t.typname
		ORDER BY n.nspname, t.typname;
	`).Scan(&enumInfos)
	if result.Error != nil {
		logs.NotezyLogger.Error(context.Background(), result.Error, "Failed to display database enums")
		return result.Error
	}

	logs.NotezyLogger.Info(context.Background(), "=============== Database Enum List ===============")
	if len(enumInfos) == 0 {
		logs.NotezyLogger.Info(context.Background(), "No enums found")
	} else {
		for index, enumInfo := range enumInfos {
			logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("%d. Type: %-30s | Values: %s", index+1, enumInfo.Name, enumInfo.Values))
		}
	}
	logs.NotezyLogger.Info(context.Background(), "=============== Database Enum List ===============")
	return nil
}

func TruncateTablesInDatabase(tableName TableName, db *gorm.DB) error {
	if logs.NotezyLogger == nil {
		return errors.New("observability logger is not initialized")
	}

	if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %q RESTART IDENTITY CASCADE;", tableName)).Error; err != nil {
		logs.NotezyLogger.Error(context.Background(), err, fmt.Sprintf("Failed to truncate database table %s", tableName))
		return err
	}

	logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Database table %s truncated", tableName))
	return nil
}

func MigrateEnumsToDatabase(db *gorm.DB, migratingEnums map[string][]string) error {
	if logs.NotezyLogger == nil {
		return errors.New("observability logger is not initialized")
	}

	logs.NotezyLogger.Info(context.Background(), "Migrating database enums...")
	for name, values := range migratingEnums {
		var exists bool
		if err := db.Raw("SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = ?);", name).Scan(&exists).Error; err != nil {
			logs.NotezyLogger.Error(context.Background(), err, fmt.Sprintf("Failed to check enum %s existence", name))
			return err
		}

		if !exists {
			enumSQL := fmt.Sprintf("CREATE TYPE %q AS ENUM ('%s');", name, strings.Join(values, "', '"))
			if err := db.Exec(enumSQL).Error; err != nil {
				logs.NotezyLogger.Error(context.Background(), err, fmt.Sprintf("Failed to create enum %s", name))
				return err
			}
			logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Enum %s created with values: %v", name, values))
			continue
		}

		var dbValues []string
		if err := db.Raw(`
			SELECT enumlabel
			FROM pg_enum
			WHERE enumtypid = (SELECT oid FROM pg_type WHERE typname = ?)
			ORDER BY enumsortorder;
		`, name).Scan(&dbValues).Error; err != nil {
			logs.NotezyLogger.Error(context.Background(), err, fmt.Sprintf("Failed to get enum %s values", name))
			return err
		}

		for _, value := range values {
			if containsString(dbValues, value) {
				continue
			}
			if err := db.Exec(fmt.Sprintf("ALTER TYPE %q ADD VALUE '%s';", name, value)).Error; err != nil {
				logs.NotezyLogger.Error(context.Background(), err, fmt.Sprintf("Failed to add value %q to enum %s", value, name))
				return err
			}
		}

		for _, dbValue := range dbValues {
			if !containsString(values, dbValue) {
				logs.NotezyLogger.Warn(context.Background(), fmt.Sprintf("Enum %s contains value %q that is not present in code", name, dbValue))
			}
		}
	}

	logs.NotezyLogger.Info(context.Background(), "Migration of enums is done")
	return nil
}

func MigrateTablesToDatabase(db *gorm.DB, migratingTables []any) error {
	if logs.NotezyLogger == nil {
		return errors.New("observability logger is not initialized")
	}

	logs.NotezyLogger.Info(context.Background(), "Migrating database tables...")
	for _, table := range migratingTables {
		if err := db.AutoMigrate(table); err != nil {
			logs.NotezyLogger.Error(context.Background(), err, "Failed to migrate table")
			return err
		}
	}

	logs.NotezyLogger.Info(context.Background(), "Migration of tables is done")
	return nil
}

func MigrateTriggersToDatabase(db *gorm.DB, migratingSQLs []string) error {
	if logs.NotezyLogger == nil {
		return errors.New("observability logger is not initialized")
	}

	return migrateSQL(db, migratingSQLs, "triggers", true)
}

func MigrateConstraintsToDatabase(db *gorm.DB, migratingSQLs []string) error {
	if logs.NotezyLogger == nil {
		return errors.New("observability logger is not initialized")
	}

	return migrateSQL(db, migratingSQLs, "constraints", false)
}

func MigrateViewsToDatabase(db *gorm.DB, migratingSQLs []string) error {
	if logs.NotezyLogger == nil {
		return errors.New("observability logger is not initialized")
	}

	return migrateSQL(db, migratingSQLs, "views", false)
}

func SeedDefaultDataToDatabase(db *gorm.DB, seedingSQLs []string) error {
	if logs.NotezyLogger == nil {
		return errors.New("observability logger is not initialized")
	}

	return migrateSQL(db, seedingSQLs, "default data", false)
}

func migrateSQL(db *gorm.DB, sqls []string, description string, skipAlreadyExists bool) error {
	logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Migrating database %s...", description))

	for _, sql := range sqls {
		for _, statement := range strings.Split(sql, SQLSeparator) {
			statement = strings.TrimSpace(statement)
			if statement == "" {
				continue
			}
			if err := db.Exec(statement).Error; err != nil {
				if skipAlreadyExists && strings.Contains(err.Error(), "SQLSTATE 42710") {
					logs.NotezyLogger.Warn(context.Background(), fmt.Sprintf("Database %s object already exists; skipping: %v", description, err))
					continue
				}
				logs.NotezyLogger.Error(context.Background(), err, fmt.Sprintf("Failed to execute %s SQL statement", description))
				return err
			}
		}
	}

	logs.NotezyLogger.Info(context.Background(), fmt.Sprintf("Migration of %s is done", description))
	return nil
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
