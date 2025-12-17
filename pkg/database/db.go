package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDB(dsn string) (*Database, error) {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Database{
		DB: database,
	}, nil
}

func (d *Database) AutoMigrate(models ...interface{}) error {
	if err := d.DB.AutoMigrate(models...); err != nil {
		return err
	}
	return nil
}

// ExecSQLFile 读取并执行某个 .sql 文件
func (d *Database) ExecSQLFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read sql file: %w", err)
	}

	if err := d.DB.Exec(string(content)).Error; err != nil {
		return fmt.Errorf("exec sql file: %w", err)
	}
	return nil
}
