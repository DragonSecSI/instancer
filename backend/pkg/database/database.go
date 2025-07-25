package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/DragonSecSI/instancer/backend/pkg/config"
	"github.com/DragonSecSI/instancer/backend/pkg/errors"
)

func NewDatabase(config *config.ConfigDatabase) (*gorm.DB, error) {
	gormConfig := gorm.Config{}

	switch config.Dialect {
	case "sqlite":
		db, err := gorm.Open(sqlite.Open(config.ConnectionString), &gormConfig)
		if err != nil {
			return nil, &errors.DatabaseConnectionError{
				Dialect: config.Dialect,
				DSN:     config.ConnectionString,
				Err:     err,
			}
		}
		return db, nil

	case "mysql":
		db, err := gorm.Open(mysql.Open(config.ConnectionString), &gormConfig)
		if err != nil {
			return nil, &errors.DatabaseConnectionError{
				Dialect: config.Dialect,
				DSN:     config.ConnectionString,
				Err:     err,
			}
		}
		return db, nil

	case "postgres":
		db, err := gorm.Open(postgres.Open(config.ConnectionString), &gormConfig)
		if err != nil {
			return nil, &errors.DatabaseConnectionError{
				Dialect: config.Dialect,
				DSN:     config.ConnectionString,
				Err:     err,
			}
		}
		return db, nil

	default:
		return nil, &errors.DatabaseConnectionError{
			Dialect: config.Dialect,
			DSN:     config.ConnectionString,
			Err:     nil,
		}
	}
}
