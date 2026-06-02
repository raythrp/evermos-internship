package testutils

import (
	"database/sql"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupMockDB returns a GORM DB backed by go-sqlmock, the mock controller,
// and the underlying *sql.DB so callers can close it in t.Cleanup.
func SetupMockDB() (*gorm.DB, sqlmock.Sqlmock, *sql.DB) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		panic("failed to create sqlmock: " + err.Error())
	}

	// GORM MySQL driver runs SELECT VERSION() during Initialize to detect capabilities.
	mock.ExpectQuery("SELECT VERSION()").
		WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("8.0.33"))

	gormDB, err := gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to open gorm with mock: " + err.Error())
	}
	return gormDB, mock, sqlDB
}
