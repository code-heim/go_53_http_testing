package models_test

import (
	"go_http_testing/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupTestDB() error {
	// Setup the MySQL test database connection string
	connStr := "net_http_blogs:tmp_pwd@tcp(127.0.0.1:3306)/net_http_blogs_test?charset=utf8&parseTime=true"
	db, err := gorm.Open(mysql.Open(connStr), &gorm.Config{})
	if err != nil {
		return err
	}

	models.DB = db
	// Run migrations
	models.DBMigrate()

	return nil
}

func TeardownTestDB() {
	var tableNames []string

	// Query to get all table names in the current database
	models.DB.Raw("SHOW TABLES").Scan(&tableNames)

	// Iterate over each table name and drop it
	for _, tableName := range tableNames {
		models.DB.Exec("DROP TABLE IF EXISTS " + tableName)
	}
}
