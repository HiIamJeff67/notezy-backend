package models

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	logs "notezy-backend/app/logs"
	schemas "notezy-backend/app/models/schemas"
	enums "notezy-backend/app/models/schemas/enums"
	util "notezy-backend/app/util"
	types "notezy-backend/shared/types"
)

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string // the port inside the container, so please leave this as 5432 for PostgreSQL
}

var (
	// the main database instance of the application (we use a different one for e2e testing, etc.)
	NotezyDB *gorm.DB

	// maintain the static information about the database instance and its config
	DatabaseInstanceToConfig = map[*gorm.DB]DatabaseConfig{}
	DatabaseNameToInstance   = map[string]*gorm.DB{}
)

var (
	PostgresDatabaseConfig = DatabaseConfig{
		Host:     util.GetEnv("DB_HOST", "notezy-db"),
		User:     util.GetEnv("DB_USER", "master"),
		Password: util.GetEnv("DB_PASSWORD", ""),
		DBName:   util.GetEnv("DB_NAME", "notezy-db"),
		Port:     util.GetEnv("DOCKER_DB_PORT", "5432"),
	}
)

func ConnectToDatabase(config DatabaseConfig) *gorm.DB {
	var dbArgs string = fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.DBName,
		config.Password,
	)

	dbConn, err := gorm.Open(postgres.Open(dbArgs), &gorm.Config{})
	if err != nil {
		logs.FError("Error connecting to the %s database\n", config.DBName)
		panic("Connecting to database error : " + err.Error())
	}

	if _, ok := DatabaseInstanceToConfig[dbConn]; !ok {
		logs.FInfo("Storing database of %s into the DatabaseInstanceToConfig...", config.DBName)
		DatabaseInstanceToConfig[dbConn] = config
	}
	if _, ok := DatabaseNameToInstance[config.DBName]; !ok {
		logs.FInfo("Storing database of %s into the DatabaseNameToInstance...", config.DBName)
		DatabaseNameToInstance[config.DBName] = dbConn
	}

	logs.FInfo("%s database connected\n", config.DBName)

	return dbConn
}

func DisconnectToDatabase(db *gorm.DB) bool {
	sqlDB, err := db.DB()
	config, ok := DatabaseInstanceToConfig[db]
	if err != nil || !ok {
		logs.FError("Failed to get the connection of the given database")
		return false
	}

	if err := sqlDB.Close(); err != nil {
		logs.FError("Failed to close the connection of %s database", config.DBName)
		return false
	}

	logs.FInfo("Extracting database of %s into the DatabaseInstanceToConfig...", config.DBName)
	delete(DatabaseInstanceToConfig, db)
	logs.FInfo("Extracting database of %s into the DatabaseNameToInstance...", config.DBName)
	delete(DatabaseNameToInstance, config.DBName)

	logs.FInfo("%s database connection closed", config.DBName)

	return true
}

func TruncateTablesInDatabase(tableName types.ValidTableName, db *gorm.DB) bool {
	result := db.Exec("TRUNCATE TABLE \"%s\" RESTART IDENTITY CASCADE;")
	if err := result.Error; err != nil {
		logs.FError("Failed to truncate %s database %s table", DatabaseInstanceToConfig[db].DBName, tableName)
		return false
	}

	logs.FInfo("%s database %s table truncated", DatabaseInstanceToConfig[db].DBName, tableName)
	return true
}

func MigrateToDatabase(db *gorm.DB) bool {
	logs.Info("Migrating enums found in models/migration.go ...")
	for name, values := range enums.MigratingEnums {
		name = strings.ToLower(name)

		// get current enum value
		var exists bool
		checkEnumSQL := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_type WHERE typname = '%s');", name)
		if err := db.Raw(checkEnumSQL).Scan(&exists).Error; err != nil {
			logs.FError("Failed to check enum %s existence: %v", name, err)
			return false
		}

		if !exists {
			// if the enum does not exist, create it
			enumSQL := fmt.Sprintf("CREATE TYPE %s AS ENUM ('%s');", name, util.JoinValues(values))
			if err := db.Exec(enumSQL).Error; err != nil {
				logs.FError("Failed to create enum %s: %v", name, err)
				return false
			}
			logs.FInfo("Enum %s created with values: %v", name, values)
		} else {
			// get current enum value
			var dbValues []string
			getValuesSQL := `
                SELECT enumlabel FROM pg_enum
                WHERE enumtypid = (SELECT oid FROM pg_type WHERE typname = ?)
                ORDER BY enumsortorder;`
			if err := db.Raw(getValuesSQL, name).Scan(&dbValues).Error; err != nil {
				logs.FError("Failed to get enum %s values: %v", name, err)
				return false
			}

			// add new values to the current enum
			for _, v := range values {
				found := false
				for _, dbv := range dbValues {
					if v == dbv {
						found = true
						break
					}
				}
				if !found {
					addValueSQL := fmt.Sprintf("ALTER TYPE %s ADD VALUE '%s';", name, v)
					if err := db.Exec(addValueSQL).Error; err != nil {
						logs.FError("Failed to add value '%s' to enum %s: %v", v, name, err)
						return false
					}
					logs.FInfo("Added value '%s' to enum %s", v, name)
				}
			}

			// check if there're values to remove
			var toRemove []string
			for _, dbv := range dbValues {
				found := false
				for _, v := range values {
					if v == dbv {
						found = true
						break
					}
				}
				if !found {
					toRemove = append(toRemove, dbv)
				}
			}
			if len(toRemove) > 0 {
				logs.FWarn("Enum %s found in code: %v", name, toRemove)
				// could choose to delete it and rebuild the enum right here
			}
		}
	}

	logs.Info("Migrating tables found in models/migration.go ...")
	for _, table := range schemas.MigratingTables {
		if err := db.AutoMigrate(table); err != nil {
			logs.FError("Failed to migrate table: %v", err)
			return false
		}
	}

	logs.Info("Migration done")

	return true
}
