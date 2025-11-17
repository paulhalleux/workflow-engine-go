package persistance

import (
	"fmt"

	"github.com/cenkalti/backoff"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateDatabase(dsn string) (*gorm.DB, error) {
	conn := func() (*gorm.DB, error) {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		sqlDB, err := db.DB()
		if err != nil {
			return nil, fmt.Errorf("failed to get database instance: %w", err)
		}
		if err := sqlDB.Ping(); err != nil {
			return nil, fmt.Errorf("database ping failed: %w", err)
		}
		return db, nil
	}

	var db *gorm.DB
	operation := func() error {
		var err error
		db, err = conn()
		return err
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 2 * 60 * 1e9 // 2 minutes

	if err := backoff.Retry(operation, expBackoff); err != nil {
		return nil, fmt.Errorf("could not connect to database after retries: %w", err)
	}

	return db, nil
}
